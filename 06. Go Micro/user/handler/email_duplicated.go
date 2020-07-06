package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/mock"
	"net/http"
	"time"
	"user/dao"
	userProto "user/proto/golang/user"
	"user/tool/jwt"
)

func (u *user) EmailDuplicated(ctx context.Context, req *userProto.EmailDuplicatedRequest, rsp *userProto.EmailDuplicatedResponse) (_ error) {
	if err := u.validate.Struct(req); err != nil {
		rsp.SetStatusAndMsg(http.StatusBadRequest, MessageBadRequest)
		return
	}

	var md metadata.Metadata
	var ok bool
	if md, ok = metadata.FromContext(ctx); !ok || md == nil {
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	var xReqId string
	if xReqId, ok = md.Get("XRequestID"); !ok || xReqId == "" {
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	if _, err := uuid.Parse(xReqId); err != nil {
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	var userId string
	if ss, ok := md.Get("Authorization"); ok && ss != "" {
		claim, err := jwt.ParseDuplicateCertClaimFromJWT(ss)
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
		return
	}

	rsp.SetStatusAndMsg(http.StatusOK, MessageEmailNotDuplicated)
	rsp.Authorization = ss
	return
}