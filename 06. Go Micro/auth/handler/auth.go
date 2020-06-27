package handler

import (
	"auth/dao"
	proto "auth/proto/auth"
	"auth/subscriber"
	"auth/tool/jwt"
	"auth/tool/random"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/broker"
	"github.com/stretchr/testify/mock"
	"net/http"
	"time"
)

type auth struct {
	mq       broker.Broker
	adc      *dao.AuthDAOCreator
	validate *validator.Validate
}

func NewAuth(mq broker.Broker, adc *dao.AuthDAOCreator, validate *validator.Validate) *auth {
	return &auth{
		mq:       mq,
		adc:      adc,
		validate: validate,
	}
}

func (e *auth) UserIdDuplicated(ctx context.Context, req *proto.UserIdDuplicatedRequest, rsp *proto.UserIdDuplicatedResponse) (_ error) {
	if err := e.validate.Struct(req); err != nil {
		rsp.SetStatusAndMsg(http.StatusBadRequest, MessageBadRequest)
		return
	}

	var email string
	if req.Authorization != "" {
		claim, err := jwt.ParseDuplicateCertClaimFromJWT(req.Authorization)
		if err != nil {
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

func (e *auth) BeforeCreateAuth(ctx context.Context, req *proto.BeforeCreateAuthRequest, rsp *proto.BeforeCreateAuthResponse) (_ error) {
	if err := e.validate.Struct(req); err != nil {
		rsp.SetStatusAndMsg(http.StatusBadRequest, MessageBadRequest)
		return
	}

	// test 환경 시 context에서 MessageId 추출, 아닐 시 새로 생성
	var mId string
	switch ctx.Value("env") {
	case "test": 	mId = ctx.Value("MessageId").(string)
	default:		mId = random.GenerateString(32)
	}

	claim, err := jwt.ParseDuplicateCertClaimFromJWT(req.Authorization)
	if err != nil {
		rsp.SetStatus(http.StatusForbidden)
		return
	}

	if claim.UserId != req.UserId {
		rsp.SetStatusAndMsg(StatusUserIdDuplicate, MessageUserIdDuplicate)
		return
	}

	if claim.Email != req.Email {
		rsp.SetStatusAndMsg(StatusEmailDuplicate, MessageEmailDuplicate)
		return
	}

	header := make(map[string]string)
	header["XRequestId"] = req.XRequestID
	header["MessageId"] = mId

	// body Marshaling 추가
	if err = e.mq.Publish(subscriber.CreateAuthEventTopic, &broker.Message{
		Header: header,
		Body:   nil,
	}); err != nil {
		rsp.SetStatus(http.StatusInternalServerError)
		return
	}

	rsp.SetStatusAndMsg(http.StatusCreated, MessageAuthCreated)
	return
}