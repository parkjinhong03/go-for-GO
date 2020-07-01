package broker

import (
	"fmt"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"os"
)

func ConnRabbitMQ() broker.Broker {
	user := os.Getenv("RABBITMQ_DEFAULT_USER")
	pwd := os.Getenv("RABBITMQ_DEFAULT_PASS")
	addr := fmt.Sprintf("amqp://%s:%s@localhost:5672", user, pwd)

	return rabbitmq.NewBroker(broker.Addrs(addr))
}