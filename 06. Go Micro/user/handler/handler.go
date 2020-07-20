package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/broker"
	"github.com/opentracing/opentracing-go"
	"user/dao"
)

type user struct {
	mq       broker.Broker
	validate *validator.Validate
	udc      *dao.UserDAOCreator
	tracer   opentracing.Tracer
}

func NewUser(mq broker.Broker, validate *validator.Validate, udc *dao.UserDAOCreator,
	tracer opentracing.Tracer) *user {
	return &user{
		mq:       mq,
		validate: validate,
		udc:      udc,
		tracer:   tracer,
	}
}
