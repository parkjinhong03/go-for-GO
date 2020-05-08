package message

import "github.com/nats-io/nats.go"

type NatsMessage interface {
	Publish(subj string, data []byte) error
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
}