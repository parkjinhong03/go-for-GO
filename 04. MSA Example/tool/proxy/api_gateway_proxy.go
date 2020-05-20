package proxy

import (
	"MSA.example.com/1/tool/message"
	"github.com/go-playground/validator/v10"
)

type apiGatewayProxy struct {
	natsM message.NatsMessage
	validate *validator.Validate
}

func NewApiGatewayProxy(nastM message.NatsMessage, validate *validator.Validate) *apiGatewayProxy {
	return &apiGatewayProxy{
		natsM:    nastM,
		validate: validate,
	}
}