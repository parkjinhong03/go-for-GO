package broker

import (
	"github.com/micro/go-micro/v2/broker"
	"github.com/stretchr/testify/mock"
)

type testRabbitMQ struct {
	mock *mock.Mock
	broker.Broker
}

func NewRabbitMQForTest(mock *mock.Mock) broker.Broker {
	return &testRabbitMQ{
		mock: mock,
	}
}

func (r *testRabbitMQ) Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {
	args := r.mock.Called(topic, m)
	return args.Error(0)
}