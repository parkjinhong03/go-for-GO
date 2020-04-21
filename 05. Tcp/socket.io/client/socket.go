package client

import (
	"fmt"
	"github.com/zhouhui8915/go-socket.io-client"
	"log"
)

func NewSocketClient(user, pwd string) *socketio_client.Client {
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["user"] = user
	opts.Query["pwd"] = pwd
	uri := "http://172.30.1.24:8080"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Fatalf("Unable to create new client, err: %v\n", err)
	}

	client.On("message", func(name string, message string) {
		fmt.Printf("%s send message '%s'\n", name, message)
	})

	return client
}