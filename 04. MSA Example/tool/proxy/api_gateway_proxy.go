package proxy

import (
	"MSA.example.com/1/tool/message"
	"github.com/go-playground/validator/v10"
)

type ApiGatewayProxy struct {
	natsM message.NatsMessage
	validate *validator.Validate
}

func NewApiGatewayProxy(nastM message.NatsMessage, validate *validator.Validate) *AuthServiceProxy {
	return &AuthServiceProxy{
		natsM:    nastM,
		validate: validate,
	}
}