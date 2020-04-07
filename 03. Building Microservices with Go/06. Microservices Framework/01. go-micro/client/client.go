package client

import (
	"context"
	"github.com/micro/go-micro/client"
	"go-micro.com/examples/1/proto"
	"log"
)

func RunClient(client client.Client, name string, request *proto.HelloRequest) *proto.HelloResponse {
	greeter := proto.NewGreeterService(name, client)
	resp , err := greeter.Hello(context.TODO(), request)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}