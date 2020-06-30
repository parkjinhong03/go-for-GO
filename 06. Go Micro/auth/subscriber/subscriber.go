package subscriber

import (
	"auth/dao"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/broker"
)

var (
	ErrorBadRequest = errors.New("bad request")
	ErrorDuplicatedMessage = errors.New("massage duplicated")
)

type auth struct {
	mq		 broker.Broker
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

