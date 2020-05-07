package message

import (
	"errors"
	"github.com/nats-io/nats.go"
	"os"
)

type natsMessage struct {
	conn *nats.Conn
}

func NewNatsMessageByEnv() (*natsMessage, error) {
	url := os.Getenv("NATS")
	if url == "" {
		return nil, errors.New("please set your NATS environment variable")
	}
	conn, err := nats.Connect("nats://" + url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return &natsMessage{conn: conn}, nil
}

func (nm *natsMessage) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return nm.conn.Subscribe(subj, cb)
}

func (nm *natsMessage) Publish(subj string, data []byte) error {
	return nm.conn.Publish(subj, data)
}