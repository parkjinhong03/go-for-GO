package subscriber

import (
	"fmt"
	"github.com/micro/go-micro/v2/broker"
)

type auth struct{}

func NewAuth() *auth {
	return &auth{}
}

func (a *auth) CreateAuth(e broker.Event) error {
	fmt.Println(string(e.Message().Body))
	return nil
}
