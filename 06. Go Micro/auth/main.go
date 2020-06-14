package main

import (
	"auth/adapter/db"
	"auth/handler"
	auth "auth/proto/auth"
	"auth/subscriber"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)


func main() {
	service := micro.NewService(
		micro.Name("examples.blog.service.auth"),
		micro.Version("latest"),
	)

	conn, err := db.ConnMysql()
	if err != nil {
		log.Fatalf("unable to connect mysql server, err: %v\n", err)
	}
	h := handler.NewAuth(conn)

 	if err := auth.RegisterAuthHandler(service.Server(), h); err != nil {
 		log.Fatal(err)
	}
	if err := micro.RegisterSubscriber("examples.blog.service.auth", service.Server(), new(subscriber.Auth)); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
