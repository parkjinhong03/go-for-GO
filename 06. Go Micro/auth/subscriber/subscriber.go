package subscriber

import (
	"auth/dao"
	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/broker"
	"github.com/opentracing/opentracing-go"
)

type Auth struct {
	mq       broker.Broker
	adc      *dao.AuthDAOCreator
	validate *validator.Validate
	tr       opentracing.Tracer
}

func NewAuth(mq broker.Broker, adc *dao.AuthDAOCreator, validate *validator.Validate,
	tr opentracing.Tracer) *Auth {

	return &Auth{
		mq:       mq,
		adc:      adc,
		validate: validate,
		tr:       tr,
	}
}

