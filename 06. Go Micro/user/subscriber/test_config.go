package subscriber

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/stretchr/testify/mock"
	"log"
	br "user/adapter/broker"
	"user/dao"
	"user/tool/validator"
)

type method string
type returns []interface{}

const (
	none = "none"
	defaultName = "박진홍"
	defaultPN = "01088378347"
	defaultEmail = "jinhong0719@naver.com"
	nilInt = 71098
)

var mockStore mock.Mock
var h *User
var userId uint
var psMsgId uint
var event = &CustomEvent{
	mock: &mockStore,
	msg:  &broker.Message{},
}

func init() {
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }
	rbMQ := br.NewRabbitMQForTest(&mockStore)
	udc := dao.NewUserDAOCreator(nil)
	h = NewUser(rbMQ, validate, udc)
}

func setUpEnv() {
	mockStore = mock.Mock{}
	event.clearMessage()
	userId = 1
	psMsgId = 1
}