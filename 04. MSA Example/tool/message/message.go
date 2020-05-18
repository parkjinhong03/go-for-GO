package message

import (
	"github.com/nats-io/nats.go"
	"time"
)

type NatsMessage interface {
	Publish(subj string, data []byte) error
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
	Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
}