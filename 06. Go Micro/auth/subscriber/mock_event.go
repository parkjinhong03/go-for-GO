package subscriber

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/stretchr/testify/mock"
)

func setUp() {
	mockStore = mock.Mock{}
	event.clearMessage()
	authId = 0
}

type CustomEvent struct {
	mock *mock.Mock
	msg *broker.Message
}

func (e *CustomEvent) Ack() error {
	args := e.mock.Called()
	return args.Error(0)
}

func (e *CustomEvent) Error() error {
	args := e.mock.Called()
	return args.Error(0)
}

func (e *CustomEvent) Topic() string {
	args := e.mock.Called()
	return args.String(0)
}

func (e *CustomEvent) Message() *broker.Message {
	return e.msg
}

func (e *CustomEvent) setMessage(header map[string]string, body []byte) {
	e.msg.Header = header
	e.msg.Body = body
}

func (e *CustomEvent) clearMessage() {
	e.msg.Header = map[string]string{}
	e.msg.Body = []byte{}
}