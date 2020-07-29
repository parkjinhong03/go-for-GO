package checker

import "github.com/micro/go-micro/v2/broker"

type Broker struct {
	broker broker.Broker
}

func NewBroker(br broker.Broker) (brc *Broker, err error) {
	_, err = br.Subscribe("examples.blog.service.auth.HealthCheck", func(e broker.Event) error {
		return nil
	}, broker.Queue("examples.blog.service.auth.HealthCheck"))
	brc = &Broker{
		broker: br,
	}
	return
}

func (b *Broker) Status() (interface{}, error) {
	err := b.broker.Publish("examples.blog.service.auth.HealthCheck", new(broker.Message))
	return nil, err
}