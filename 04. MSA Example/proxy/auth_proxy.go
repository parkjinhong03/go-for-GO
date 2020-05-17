package proxy

import (
	"MSA.example.com/1/tool/message"
	"encoding/json"
)

type AuthServiceProxy struct {
	natsM message.NatsMessage
}

func NewAuthServiceProxy(natsM message.NatsMessage) *AuthServiceProxy {
	return &AuthServiceProxy{
		natsM: natsM,
	}
}

func (ap *AuthServiceProxy) Write(b []byte) (i int, err error) {
	in := input{}
	if err = json.Unmarshal(b, &in); err != nil { return }
	if err = ap.natsM.Publish(in.InputChannel, b); err != nil { return }
	return
}

type input struct {
	InputChannel string
}