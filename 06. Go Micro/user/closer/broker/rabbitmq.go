package broker

import (
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"user/subscriber"
)

func RabbitMQInitializer(s server.Server, us *subscriber.User) func() error {
	return func() error {
		brk := s.Options().Broker
		if err := brk.Connect(); err != nil { log.Fatal(err) }
		options := []broker.SubscribeOption{broker.Queue(subscriber.CreateUserEventTopic), broker.DisableAutoAck(), rabbitmq.DurableQueue()}
		if _, err := brk.Subscribe(subscriber.CreateUserEventTopic, us.CreateUser, options...); err != nil { log.Fatal(err) }
		log.Infof("succeed in connecting to broker!! (name: %s | addr: %s)\n",  brk.String(), brk.Address())
		return nil
	}
}