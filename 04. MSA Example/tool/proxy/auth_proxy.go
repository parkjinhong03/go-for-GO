package proxy

import (
	"MSA.example.com/1/protocol"
	"MSA.example.com/1/tool/customError"
	"MSA.example.com/1/tool/message"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"time"
)

type AuthServiceProxy struct {
	natsM 	 message.NatsMessage
	validate *validator.Validate
}

func NewAuthServiceProxy(natsM message.NatsMessage, validate *validator.Validate) *AuthServiceProxy {
	return &AuthServiceProxy{
		natsM: natsM,
		validate: validate,
	}
}

func (ap *AuthServiceProxy) Write(b []byte) (int, error) {
	myErr := &customError.ProxyWriteError{}
	
	r := struct {
		Required protocol.RequiredProtocol
	}{}
	if err := json.Unmarshal(b, &r); err != nil {
		myErr.Err = err
		return 0, myErr
	}
	// json 외의 포맷팅 방식을 이용한 Unmarshal 추가 필요
	if err := ap.validate.Struct(&r); err != nil {
		myErr.Err = err
		return 0, myErr
	}

	var err error
	switch r.Required.Usage {
	case "AuthSignUpRequest":
		msg, err := ap.natsM.Request(r.Required.InputChannel, b, 5 * time.Second)
		myErr.Err = err
		myErr.ReturnMsg = msg
		return 0, myErr
	case "UserRegisterResponse":
		err := ap.natsM.Publish(r.Required.InputChannel, b)
		return 0, err
	default:
		err = errors.New("this Usage is undefined so cannot be processed")
	}
	return 0, err
}
