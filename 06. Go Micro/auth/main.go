package main

import (
	"auth/adapter/broker"
	"auth/adapter/db"
	"auth/dao"
	"auth/handler"
	auth "auth/proto/auth"
	"auth/subscriber"
	"auth/tool/validator"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

func main() {
	// 플러그인 객체 생성
	rbMQ := broker.ConnRabbitMQ()

	// 서비스 생성
	service := micro.NewService(
		micro.Name("examples.blog.service.auth"),
		micro.Version("latest"),
		micro.Broker(rbMQ),
	)

	// 이벤트 핸들러 객체 생성
	s := subscriber.NewAuth()

	// 메시지 브로커(MQ) 핸들러 생성
	brkHandleFunc := func() (err error) {
		brk := service.Options().Broker
		if err = brk.Connect(); err != nil { log.Fatal(err) }
		if _, err = brk.Subscribe(subscriber.CreateAuthEventTopic, s.CreateAuth); err != nil { log.Fatal(err) }
		log.Infof("succeed in connecting to broker!! (name: %s | addr: %s)\n",  brk.String(), brk.Address())
		return
	}

	// 서비스 초기화
	service.Init(
		micro.AfterStart(brkHandleFunc),
	)

	// 의존성 주입을 위한 객체 생성
	mq := service.Options().Broker
	conn, err := db.ConnMysql()
	if err != nil {
		log.Fatalf("unable to connect mysql server, err: %v\n", err)
	}
	adc := dao.NewAuthDAOCreator(conn)
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }

	// rpc 핸들러 객체 생성
	h := handler.NewAuth(mq, adc, validate)

	// 핸들러 등록 및 서비스 실행
 	if err := auth.RegisterAuthHandler(service.Server(), h); err != nil {
 		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
