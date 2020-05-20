package proxy

import (
	"MSA.example.com/1/protocol"
	"MSA.example.com/1/tool/message"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
)

type ApiGatewayProxy struct {
	natsM message.NatsMessage
	validate *validator.Validate
}

func NewApiGatewayProxy(nastM message.NatsMessage, validate *validator.Validate) *ApiGatewayProxy {
	return &ApiGatewayProxy{
		natsM:    nastM,
		validate: validate,
	}
}

func (ap *ApiGatewayProxy) Write(b []byte) (i int, err error) {
	r := struct {
		Required protocol.RequiredProtocol
	}{}
	if err = json.Unmarshal(b, &r); err != nil { return }
	if err = ap.validate.Struct(&r); err != nil { return }


	switch r.Required.Usage {
	case "AuthSignUpResponse":
		err = ap.natsM.Publish(r.Required.InputChannel, b)
	default:
		err = errors.New("this Usage is undefined so cannot be processed")
	}
	return
}