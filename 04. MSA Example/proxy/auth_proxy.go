package proxy

import (
	"MSA.example.com/1/tool/customError"
	"MSA.example.com/1/tool/message"
	"encoding/json"
	"errors"
	"time"
)

type AuthServiceProxy struct {
	natsM message.NatsMessage
}

func NewAuthServiceProxy(natsM message.NatsMessage) *AuthServiceProxy {
	return &AuthServiceProxy{
		natsM: natsM,
	}
}

func (ap *AuthServiceProxy) Write(b []byte) (int, error) {
	myErr := &customError.ProxyWriteError{}
	in := input{}
	if err := json.Unmarshal(b, &in); err != nil {
		myErr.Err = err
		return 0, myErr
	}

	switch in.InputChannel {
	case "auth.signup":
		msg, err := ap.natsM.Request(in.InputChannel, b, 1 * time.Second)
		if err != nil {
			myErr.Err = err
			return 0, myErr
		}
		myErr.ReturnMsg = msg
	default:
		err := errors.New("this channel is currently unavailable")
		myErr.Err = err
	}

	return 0, myErr
}

type input struct {
	InputChannel string
}