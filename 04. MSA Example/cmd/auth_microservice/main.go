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
	db, err := dbc.ConnectDB("AuthMicroservice")
	if err != nil {
		log.Fatalf("unable to connect DB server, err: %v\n", err)
	}

	userD := dataservice.GetUserDAO(db)
	natsM, err := message.GetDefaultNatsByEnv()
	if err != nil {
		log.Fatalf("unable to connect NATS server, err: %v\n", err)
	}
	validate := validator.New()
	apiNatsE := natsEncoder.NewJsonEncoder(proxy.NewApiGatewayProxy(natsM, validate))
	userNatsE := natsEncoder.NewJsonEncoder(proxy.NewUserServiceProxy(natsM, validate))
	u := usecase.NewAuthDefaultUseCase(userD, validate, apiNatsE, userNatsE)

	_, err = natsM.Subscribe("auth.signup", u.SignUpMsgHandler)
	if err != nil {
		log.Fatalf("unable to subscribe auth.login from nats message broker, err: %v\n", err)
	}
	_, err = natsM.Subscribe("user.registry.reply", u.RegistryReplyMsgHandler)
	if err != nil {
		log.Fatalf("unable to subscribe user.registry.reply from nats message broker, err: %v\n", err)
	}
	log.Println("Auth message pub/sub server is completely started.")
	runtime.Goexit()
}

