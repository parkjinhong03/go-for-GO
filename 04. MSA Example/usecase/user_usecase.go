package usecase

import (
	"MSA.example.com/1/dataservice"
	natsEncoder "MSA.example.com/1/tool/encoder/nats"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
)

type userUseCase struct {
	authNatsE     natsEncoder.Encoder
	userInformDAO dataservice.UserInformDataService
	validate      *validator.Validate
}

func NewUserUseCase(
	authNatsE natsEncoder.Encoder, userInformDAO dataservice.UserInformDataService, validate *validator.Validate) *userUseCase {
		return &userUseCase{
			authNatsE:     authNatsE,
			userInformDAO: userInformDAO,
			validate:      validate,
		}
}

func (u *userUseCase) RegistryMsgHandler(msg *nats.Msg) {

}