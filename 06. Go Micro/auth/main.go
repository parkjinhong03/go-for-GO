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
	br "github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
)

func main() {
	// 플러그인 객체 생성
	rbMQ := broker.ConnRabbitMQ()

	// 의존성 주입을 위한 객체 생성
	conn, err := db.ConnMysql()
	if err != nil {
		log.Fatalf("unable to connect mysql server, err: %v\n", err)
	}
	adc := dao.NewAuthDAOCreator(conn)
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }

	// 서비스 생성
	service := micro.NewService(
		micro.Name("examples.blog.service.auth"),
		micro.Version("latest"),
		micro.Broker(rbMQ),
	)

	// 이벤트 및 rpc 핸들러 객체 생성
	s := subscriber.NewMsgHandler(adc, validate)
	// mq := service.Options().Broker
	h := handler.NewAuth(rbMQ, adc, validate)

	// 초기화 핸들러 함수 생성
	brkHandleFunc := func() (err error) {
		brk := service.Options().Broker
		if err = brk.Connect(); err != nil { log.Fatal(err) }
		_, err = brk.Subscribe(subscriber.CreateAuthEventTopic, s.CreateAuth,
			br.Queue(subscriber.CreateAuthEventTopic), // Queue 정적 이름 설정
			br.DisableAutoAck(), // Ack를 수동으로 실행하게 설정
			rabbitmq.DurableQueue()) // Queue 연결을 종료해도 삭제X 설정
		if err != nil { log.Fatal(err) }
		log.Infof("succeed in connecting to broker!! (name: %s | addr: %s)\n",  brk.String(), brk.Address())
		return
	}

	// 서비스 초기화 등록
	service.Init(
		micro.AfterStart(brkHandleFunc),
	)

	// 핸들러 등록 및 서비스 실행
 	if err := auth.RegisterAuthHandler(service.Server(), h); err != nil {
 		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
