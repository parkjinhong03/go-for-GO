package subscriber

import (
	"auth/dao"
	"auth/tool/validator"
	"github.com/micro/go-micro/v2/broker"
	"github.com/stretchr/testify/mock"
	"log"
)

var mockStore mock.Mock
var authId uint
var h *auth
var event = &CustomEvent{
	mock: &mockStore,
	msg:  &broker.Message{},
}

const (
	none = "none"
	defaultUserId = "TestId"
	defaultUserPw = "TestPw"
	defaultName = "박진홍"
	defaultPN = "01088378347"
	defaultEmail = "jinhong0719@naver.com"
)

func init() {
	adc := dao.NewAuthDAOCreator(nil)
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }
	h = NewAuth(adc, validate)
}

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