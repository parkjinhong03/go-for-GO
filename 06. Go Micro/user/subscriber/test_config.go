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
)

var mockStore mock.Mock
var h *user
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
	userId = 0
	psMsgId = 0
}