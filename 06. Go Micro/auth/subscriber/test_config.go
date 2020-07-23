package subscriber

import (
	br "auth/adapter/broker"
	"auth/dao"
	"auth/tool/validator"
	"github.com/micro/go-micro/v2/broker"
	"github.com/stretchr/testify/mock"
	"github.com/uber/jaeger-client-go"
	"log"
)

var mockStore mock.Mock
var authId uint
var psMsgId uint
var h *Auth
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
	noneInt uint32 = 67981345
)

func init() {
	adc := dao.NewAuthDAOCreator(nil)
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }
	rbMQ := br.NewRabbitMQForTest(&mockStore)
	// jaeger.Tracer 모의 객체 추가 예정
	h = NewAuth(rbMQ, adc, validate, &jaeger.Tracer{})
}

func setUp() {
	mockStore = mock.Mock{}
	event.clearMessage()
	authId = 1
	psMsgId = 1
}