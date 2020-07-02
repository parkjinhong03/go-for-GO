package handler

import (
	"context"
	"github.com/stretchr/testify/mock"
	"net/http"
	"time"
	"user/dao"
	proto "user/proto/user"
	"user/tool/jwt"
)

func (u *user) EmailDuplicated(ctx context.Context, req *proto.EmailDuplicatedRequest, rsp *proto.EmailDuplicatedResponse) (_ error) {
	if err := u.validate.Struct(req); err != nil {
		rsp.SetStatusAndMsg(http.StatusBadRequest, MessageBadRequest)
		return
	}

	var userId string
	if req.Authorization != "" {
		claim, err := jwt.ParseDuplicateCertClaimFromJWT(req.Authorization)
		if err != nil { rsp.SetStatus(http.StatusForbidden); return }
		userId = claim.UserId
	}

	var ud dao.UserDAOService
	switch ctx.Value("env") {
	case "test":
		mockStore := ctx.Value("mockStore").(*mock.Mock)
		ud = u.udc.GetTestUserDAO(mockStore)
	default:
		ud = u.udc.GetDefaultUserDAO()
	}

	exist, err := ud.CheckIfEmailExist(req.Email)
	if err != nil {
		rsp.SetStatus(http.StatusInternalServerError)
		return
	}

	if exist {
		rsp.SetStatusAndMsg(StatusEmailDuplicated, MessageEmailDuplicated)
		return
	}

	ss, err := jwt.GenerateDuplicateCertJWT(userId, req.Email, time.Hour)
	if err != nil {
		rsp.SetStatus(http.StatusInternalServerError)
	}

	rsp.SetStatusAndMsg(http.StatusOK, MessageEmailNotDuplicated)
	rsp.Authorization = ss
	return
}