package proxy

import (
	"MSA.example.com/1/tool/message"
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

}