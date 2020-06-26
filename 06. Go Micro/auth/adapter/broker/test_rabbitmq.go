package broker

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/stretchr/testify/mock"
)

type testRabbitMq struct {
	mock *mock.Mock
	broker.Broker
}

func NewRabbitMQForTest(mock *mock.Mock) broker.Broker {
	return testRabbitMq{
		mock:   mock,
	}
}
