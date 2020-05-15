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

type DefaultNatsMessage struct {
	*nats.Conn
}

func GetDefaultNatsByEnv() (*DefaultNatsMessage, error) {
	address := os.Getenv("NATS")
	if address == "" {
		return nil, errors.New("please set your NATS environment variable")
	}
	return GetDefaultNats(address)
}

func GetDefaultNatsByFlag() (*DefaultNatsMessage, error) {
	flag.Parse()
	if *address == "" {
		return nil, errors.New("please set your nats command line flag")
	}
	return GetDefaultNats(*address)
}

func GetDefaultNats(address string) (*DefaultNatsMessage, error) {
	conn, err := nats.Connect("nats://" + address)
	if err != nil {
		return nil, err
	}

	return &DefaultNatsMessage{Conn: conn}, nil
}

func (nm *DefaultNatsMessage) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	return nm.Conn.Subscribe(subj, cb)
}

func (nm *DefaultNatsMessage) Publish(subj string, data []byte) error {
	return nm.Conn.Publish(subj, data)
}