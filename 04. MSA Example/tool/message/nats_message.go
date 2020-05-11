package message

import (
	"errors"
	"github.com/nats-io/nats.go"
	"os"
)

type defaultNatsMessage struct {
	conn *nats.Conn
}

func DefaultNatsMessageByEnv() (*defaultNatsMessage, error) {
	url := os.Getenv("NATS")
	if url == "" {
		return nil, errors.New("please set your NATS environment variable")
	}
	conn, err := nats.Connect("nats://" + url)
	if err != nil {
		return nil, err
	}

	return &defaultNatsMessage{conn: conn}, nil
}

func (nm *defaultNatsMessage) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return nm.conn.Subscribe(subj, cb)
}

func (nm *defaultNatsMessage) Publish(subj string, data []byte) error {
	return nm.conn.Publish(subj, data)
}