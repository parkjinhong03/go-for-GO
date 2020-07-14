package handler

import (
	userProto "gateway/proto/golang/user"
)

type userHandler struct {
	uc userProto.UserService
}

func NewUserHandler(uc userProto.UserService) *userHandler {
	return &userHandler{
		uc: uc,
	}
}

