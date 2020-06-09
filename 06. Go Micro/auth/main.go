package main

import (
	"auth/handler"
	auth "auth/proto/auth"
	"auth/subscriber"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("examples.blog.service.auth"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	auth.RegisterAuthHandler(service.Server(), new(handler.Auth))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("examples.blog.service.auth", service.Server(), new(subscriber.Auth))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
