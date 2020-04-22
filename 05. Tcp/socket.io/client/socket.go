package client

import (
	"github.com/zhouhui8915/go-socket.io-client"
	"log"
)

func NewSocketClient() *socketClient {
	var incoming = make(chan chat)
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	uri := "http://ec2-18-218-105-177.us-east-2.compute.amazonaws.com:8080/socket.io/"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Fatalf("Unable to create new client, err: %v\n", err)
	}

	_ = client.On("MESSAGE", func(name string, message string) {
		incoming <- chat{Name: name, Message: message}
	})

	return &socketClient{
		Client:   client,
		incoming: incoming,
	}
}

func (sc *socketClient) Incoming() <-chan chat {
	return sc.incoming
}