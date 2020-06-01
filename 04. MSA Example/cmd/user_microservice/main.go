package main

import (
	"MSA.example.com/1/dataservice"
	"MSA.example.com/1/tool/dbc"
	natsEncoder "MSA.example.com/1/tool/encoder/nats"
	"MSA.example.com/1/tool/message"
	"MSA.example.com/1/tool/proxy"
	"MSA.example.com/1/usecase"
	"github.com/go-playground/validator/v10"
	"log"
	"runtime"
)

func main() {
	conn, err := dbc.ConnectDB("UserMicroservice")
	if err != nil {
		log.Fatalf("unable to connect UserMicroservice databases, err: %v\n", err)
	}
	userInformD := dataservice.NewUserInformDAO(conn)
	natsM, err := message.GetDefaultNatsByEnv()
	if err != nil {
		log.Fatalf("unable to connect nats server, err: %v\n", err)
	}
	validate := validator.New()
	authP := proxy.NewAuthServiceProxy(natsM, validate)
	authJsonE := natsEncoder.NewJsonEncoder(authP)
	u := usecase.NewUserUseCase(authJsonE, userInformD, validate)

	_, err = natsM.Subscribe("user.registry", u.RegistryMsgHandler)
	if err != nil {
		log.Fatalf("some error occurs while subscribing user.registry channel, err: %v\n", err)
	}
	runtime.Goexit()
}