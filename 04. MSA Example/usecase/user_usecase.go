package usecase

import (
	"MSA.example.com/1/dataservice"
	"MSA.example.com/1/tool/message"
	"github.com/go-playground/validator/v10"
)

type userUseCase struct {
	authJsonE     message.NatsMessage
	userInformDAO dataservice.UserInformDataService
	validate      *validator.Validate
}

func NewUserUseCase(
	authJsonE message.NatsMessage, userInformDAO dataservice.UserInformDataService, validate *validator.Validate) *userUseCase {
		return &userUseCase{
			authJsonE:     authJsonE,
			userInformDAO: userInformDAO,
			validate:      validate,
		}
}