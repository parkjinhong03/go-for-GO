package main

import (
	"MSA.example.com/1/dataservice"
	"MSA.example.com/1/dataservice/userdata"
	"MSA.example.com/1/tool/message"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"runtime"
)

func main() {
	dbc, err := dataservice.ConnectDB("AuthMicroservice")
	if err != nil {
		log.Fatalf("unable to connect DB server, err: %v\n", err)
	}

	_ = userdata.GetUserDAO(dbc)
	natsM, err := message.GetDefaultNatsByEnv()
	if err != nil {
		log.Fatalf("unable to connect NATS server, err: %v\n", err)
	}

	_, err = natsM.Subscribe("auth.signup", signUpMsgHandler)
	// handler들을 메서드로 가지고 있는 usecase 구조체 생성 추가 (db, nats 필드 소유)
	if err != nil {
		log.Fatalf("unable to subscribe auth.login from nats message broker, err: %v\n", err)
	}
	runtime.Goexit()
}

func signUpMsgHandler(msg *nats.Msg) {
	fmt.Println(string(msg.Data))
	msg.Reply = "auth.signup.reply"
	err := msg.Respond([]byte("ok"))
	fmt.Println(err)
}
