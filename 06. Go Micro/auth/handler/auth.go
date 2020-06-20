package handler

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	proto "auth/proto/auth"
	"auth/tool/jwt"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/stretchr/testify/mock"
	"net/http"
	"time"
)

const (
	StatusUserIdDuplicate = 470
	StatusColumnLengthOver = 570
	StatusBcryptNotHashed = 571
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

	ss, err := jwt.GenerateDuplicateCertJWT(req.UserId, time.Hour)
	if err != nil {
		rsp.SetResponse(http.StatusInternalServerError, err.Error())
		return
	}

	rsp.SetResponse(http.StatusOK, "this user ID can be used")
	rsp.Authorization = ss
	return
}

func (e *auth) CreateAuth(ctx context.Context, req *proto.CreateAuthRequest, rsp *proto.CreateAuthResponse) (_ error) {
	if err := e.validate.Struct(req); err != nil {
		rsp.Status = http.StatusBadRequest
		rsp.Message = err.Error()
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
		rsp.Status = http.StatusInternalServerError
		rsp.Message = err.Error()
		return
	}

	switch exist {
	case true:
		rsp.Status = StatusUserIdDuplicate
		rsp.Message = user.IdDuplicateError.Error()
	case false:
		rsp.Status = http.StatusCreated
		rsp.Message = "succeed in creating new auth"
	}
	return

	// 고루틴으로 실행
	_, err = ad.Insert(&model.Auth{
		UserId: req.UserId,
		UserPw: req.UserPw,
	})

	switch err {
	case nil:
		rsp.Status = http.StatusCreated
		rsp.Message = "create auth success"
		ad.Commit()
	case user.IdDuplicateError:
		rsp.Status = StatusUserIdDuplicate
		rsp.Message = user.IdDuplicateError.Error()
		ad.Rollback()
	case user.DataLengthOverError:
		rsp.Status = StatusColumnLengthOver
		rsp.Message = user.DataLengthOverError.Error()
	case user.BcryptGenerateError:
		rsp.Status = StatusBcryptNotHashed
		rsp.Message = user.BcryptGenerateError.Error()
	default: // Unknown Error
		rsp.Status = http.StatusInternalServerError
		rsp.Message = err.Error()
	}

	log.Info("Received Auth.CreateAuth request")
	return nil
}

