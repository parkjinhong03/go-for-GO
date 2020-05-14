package message

import (
	"errors"
	"flag"
	"github.com/nats-io/nats.go"
	"os"
)

var address *string
func init() {
	address = flag.String("nats", "", "Nats server address")
}

type defaultNatsMessage struct {
	*nats.Conn
}

func DefaultNatsMessageByEnv() (*defaultNatsMessage, error) {
	address := os.Getenv("NATS")
	if address == "" {
		return nil, errors.New("please set your NATS environment variable")
	}
	return DefaultNatsMessage(address)
}

func DefaultNatsMessageByFlag() (*defaultNatsMessage, error) {
	flag.Parse()
	if *address == "" {
		return nil, errors.New("please set your nats command line flag")
	}
	return DefaultNatsMessage(*address)
}

func DefaultNatsMessage(address string) (*defaultNatsMessage, error) {
	conn, err := nats.Connect("nats://" + address)
	if err != nil {
		return nil, err
	}

	return &defaultNatsMessage{Conn: conn}, nil
}

func (nm *defaultNatsMessage) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return nm.Conn.Subscribe(subj, cb)
}

func (nm *defaultNatsMessage) Publish(subj string, data []byte) error {
	return nm.Conn.Publish(subj, data)
}