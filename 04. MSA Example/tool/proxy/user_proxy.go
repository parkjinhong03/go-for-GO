package proxy

import (
	"MSA.example.com/1/protocol"
	"MSA.example.com/1/tool/message"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
)

type UserServiceProxy struct {
	natsM    message.NatsMessage
	validate *validator.Validate
}

func NewUserServiceProxy(natsM message.NatsMessage, validate *validator.Validate) *UserServiceProxy {
	return &UserServiceProxy{
		natsM:    natsM,
		validate: validate,
	}
}

func (up *UserServiceProxy) Write(p []byte) (i int, err error) {
	r := struct {
		Required protocol.RequiredProtocol
	}{}
	if err = json.Unmarshal(p, &r); err != nil { return }
	if err = up.validate.Struct(&r); err != nil { return }

	switch r.Required.Usage {
	case "UserRegistryPublish":
		err = up.natsM.Publish("user.registry", p)
	default:
		err = errors.New(r.Required.Usage + " is undefined usage, so cannot be processed")
	}
	return
}