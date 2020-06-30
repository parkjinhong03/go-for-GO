package handler

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/broker"
	"user/dao"
	proto "user/proto/user"
)

type user struct {
	mq       broker.Broker
	validate validator.Validate
	udc      dao.UserDAOCreator
}

func NewUser(mq broker.Broker, validate validator.Validate, udc dao.UserDAOCreator) *user {
	return &user{
		mq:       mq,
		validate: validate,
		udc:      udc,
	}
}

func (u *user) EmailDuplicated(ctx context.Context, req *proto.EmailDuplicatedRequest, resp *proto.EmailDuplicatedResponse) (_ error) {
	return
}