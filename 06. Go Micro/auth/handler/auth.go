package handler

import (
	"auth/dao"
	"auth/dao/user"
	"auth/model"
	"context"
	"github.com/go-playground/validator/v10"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/stretchr/testify/mock"
	"net/http"

	proto "auth/proto/auth"
)

const (
	StatusUserIdDuplicate = 470
	StatusColumnLengthOver = 570
	StatusBcryptNotHashed = 571
)

type auth struct{
	adc *dao.AuthDAOCreator
	validate *validator.Validate
}

func NewAuth(adc *dao.AuthDAOCreator, validate *validator.Validate) *auth {
	return &auth{
		adc:      adc,
		validate: validate,
	}
}

func (e *auth) CreateAuth(ctx context.Context, req *proto.CreateAuthRequest, rsp *proto.CreateAuthResponse) error {
	if err := e.validate.Struct(req); err != nil {
		rsp.Status = http.StatusBadRequest
		rsp.Message = err.Error()
		return nil
	}

	var ad dao.AuthDAOService
	switch ctx.Value("env") {
	case "test":
		mockStore := ctx.Value("mockStore").(*mock.Mock)
		ad = e.adc.GetTestAuthDAO(mockStore)
	default:
		ad = e.adc.GetDefaultAuthDAO()
	}

	_, err := ad.Insert(&model.Auth{
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
