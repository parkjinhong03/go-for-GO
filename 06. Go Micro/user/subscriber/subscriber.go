package subscriber

import (
	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/broker"
	"user/dao"
)

type User struct {
	mq       broker.Broker
	udc      *dao.UserDAOCreator
	validate *validator.Validate
}

func NewUser(mq broker.Broker, validate *validator.Validate, udc *dao.UserDAOCreator) *User {
	return &User{
		mq:       mq,
		udc:      udc,
		validate: validate,
	}
}
