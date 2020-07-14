package handler

import (
	"auth/dao"
	authProto "auth/proto/golang/auth"
	"auth/tool/jwt"
	"context"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/stretchr/testify/mock"
	"net/http"
	"time"
)

func (e *auth) UserIdDuplicated(ctx context.Context, req *authProto.UserIdDuplicatedRequest, rsp *authProto.UserIdDuplicatedResponse) (_ error) {
	if err := e.validate.Struct(req); err != nil {
		rsp.SetStatusAndMsg(http.StatusBadRequest, MessageBadRequest)
		return
	}
	var md metadata.Metadata
	var ok bool
	if md, ok = metadata.FromContext(ctx); !ok || md == nil {
		log.Info(http.StatusForbidden, MessageUnableGetMetadata)
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	var xId string
	if xId, ok = md.Get("X-Request-Id"); !ok {
		log.Info(http.StatusForbidden, MessageThereIsNoXReqId)
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	if _, err := uuid.Parse(xId); err != nil {
		log.Info(http.StatusForbidden, MessageInvalidXReqId)
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	var email string
	if ss, ok := md.Get("Unique-Authorization"); ok && ss != "" {
		claim, err := jwt.ParseDuplicateCertClaimFromJWT(ss)
		if err != nil {
			log.Info(http.StatusForbidden, MessageUnableParseJwt)
			rsp.SetStatus(http.StatusForbidden)
			return
		}
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
		log.Info(http.StatusInternalServerError, MessageUnableCheckUserId)
		rsp.SetStatus(http.StatusInternalServerError)
		return
	}

	if exist {
		log.Info(StatusUserIdDuplicate, MessageUserIdDuplicate)
		rsp.SetStatusAndMsg(StatusUserIdDuplicate, MessageUserIdDuplicate)
		return
	}

	ss, err := jwt.GenerateDuplicateCertJWT(req.UserId, email, time.Hour)
	if err != nil {
		log.Info(http.StatusInternalServerError, MessageUnableGenerateJwt)
		rsp.SetStatus(http.StatusInternalServerError)
		return
	}

	log.Info(http.StatusOK, MessageUserIdNotDuplicated)
	rsp.SetStatusAndMsg(http.StatusOK, MessageUserIdNotDuplicated)
	rsp.Authorization = ss
	return
}