package handler

import (
	authProto "gateway/proto/golang/auth"
)

type authHandler struct {
	ac authProto.AuthService
}

func NewAuthHandler(ac authProto.AuthService) *authHandler {
	return &authHandler{
		ac: ac,
	}
}
