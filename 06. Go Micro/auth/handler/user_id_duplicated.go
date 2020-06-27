package handler

import (
	"auth/dao"
	proto "auth/proto/auth"
	"auth/tool/jwt"
	"context"
	"github.com/stretchr/testify/mock"
	"net/http"
	"time"
)

func (e *auth) UserIdDuplicated(ctx context.Context, req *proto.UserIdDuplicatedRequest, rsp *proto.UserIdDuplicatedResponse) (_ error) {
	if err := e.validate.Struct(req); err != nil {
		rsp.SetStatusAndMsg(http.StatusBadRequest, MessageBadRequest)
		return
	}

	var email string
	if req.Authorization != "" {
		claim, err := jwt.ParseDuplicateCertClaimFromJWT(req.Authorization)
		if err != nil { rsp.SetStatus(http.StatusForbidden); return }
		email = claim.Email
	}

	var ad dao.AuthDAOService
	switch ctx.Value("env") {
	case "test":
		mockStore := ctx.Value("mockStore").(*mock.Mock)
		ad = e.adc.GetTestAuthDAO(mockStore)
	default:
		ad = e.adc.GetDefaultAuthDAO()
	}

	// 로깅 추가
	exist, err := ad.CheckIfUserIdExist(req.UserId)
	if err != nil {
		rsp.SetStatus(http.StatusInternalServerError)
		return
	}

	if exist {
		rsp.SetStatusAndMsg(StatusUserIdDuplicate, MessageUserIdDuplicate)
		return
	}

	ss, err := jwt.GenerateDuplicateCertJWT(req.UserId, email, time.Hour)
	if err != nil {
		rsp.SetStatus(http.StatusInternalServerError)
		return
	}

	rsp.SetStatusAndMsg(http.StatusOK, MessageUserIdNotDuplicated)
	rsp.Authorization = ss
	return
}