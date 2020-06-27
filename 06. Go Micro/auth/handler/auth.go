package handler

import (
	"auth/dao"
	"auth/dao/user"
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

func (e *auth) CheckIfUserIdExist(ctx context.Context, req *proto.UserIdExistRequest, rsp *proto.UserIdExistResponse) (_ error) {
	if err := e.validate.Struct(req); err != nil {
		rsp.SetResponse(http.StatusBadRequest, err.Error())
		return
	}

	var ad dao.AuthDAOService
	switch ctx.Value("env") {
	case "test":
		mockStore := ctx.Value("mockStore").(*mock.Mock)
		ad = e.adc.GetTestAuthDAO(mockStore)
	default:
		ad = e.adc.GetDefaultAuthDAO()
	}

	exist, err := ad.CheckIfUserIdExists(req.UserId)
	if err != nil {
		rsp.SetResponse(http.StatusInternalServerError, err.Error())
		return
	}

	if exist {
		rsp.SetResponse(StatusUserIdDuplicate, user.IdDuplicateError.Error())
		return
	}

	ss, err := jwt.GenerateDuplicateCertJWT(req.UserId, "", time.Hour)
	if err != nil {
		rsp.SetResponse(http.StatusInternalServerError, err.Error())
		return
	}

	rsp.SetResponse(http.StatusOK, "this user ID can be used")
	rsp.Authorization = ss
	return
}

func (e *auth) BeforeCreateAuth(ctx context.Context, req *proto.BeforeCreateAuthRequest, rsp *proto.BeforeCreateAuthResponse) (_ error) {
	if err := e.validate.Struct(req); err != nil {
		rsp.Status = http.StatusBadRequest
		rsp.Message = err.Error()
		return
	}

	// test 환경 시 context에서 MessageId 추출, 아닐 시 새로 생성
	var mId string
	switch ctx.Value("env") {
	case "test": 	mId = ctx.Value("MessageId").(string)
	default: 		mId = random.GenerateString(32)
	}

	claim, err := jwt.ParseDuplicateCertClaimFromJWT(req.Authorization)
	if err != nil {
		rsp.Status = http.StatusInternalServerError
		return
	}

	if claim.UserId != req.UserId {
		rsp.Status = StatusUserIdDuplicate
		rsp.Message = "this user iD is already in use"
		return
	}

	if claim.Email != req.Email {
		rsp.Status = StatusEmailDuplicate
		rsp.Message = "this email is already in use"
		return
	}

	header := make(map[string]string)
	header["XRequestId"] = req.XRequestID
	header["MessageId"] = mId

	if err = e.mq.Publish(subscriber.CreateAuthEventTopic, &broker.Message{
		Header: header,
		Body:   nil,
	}); err != nil {
		rsp.Status = http.StatusInternalServerError
		return
	}

	rsp.Status = http.StatusCreated
	rsp.Message = "User created reservation has been successfully processed."
	return
}