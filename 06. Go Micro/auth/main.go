package main

import (
	"auth/adapter/broker"
	"auth/adapter/db"
	"auth/dao"
	"auth/handler"
	authProto "auth/proto/golang/auth"
	"auth/subscriber"
	"auth/tool/validator"
	topic "auth/topic/golang"
	"github.com/micro/go-micro/v2"
	br "github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/micro/go-plugins/registry/consul/v2"
)

func main() {
	// 플러그인 객체 생성
	rbMQ := broker.ConnRabbitMQ()
	cs := consul.NewRegistry(registry.Addrs("http:localhost:8500"))

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
		micro.Registry(cs),
	)

	// 이벤트 및 rpc 핸들러 객체 생성
	s := subscriber.NewAuth(rbMQ, adc, validate)
	// mq := service.Options().Broker
	h := handler.NewAuth(rbMQ, adc, validate)

	// Broker 초기화 핸들러 함수 생성
	brkHandleFunc := func() (err error) {
		brk := service.Options().Broker

		if err = brk.Connect(); err != nil { return }
		_, err = brk.Subscribe(topic.CreateAuthEventTopic, s.CreateAuth,
			br.Queue(topic.CreateAuthEventTopic), // Queue 정적 이름 설정
			br.DisableAutoAck(), // Ack를 수동으로 실행하게 설정
			rabbitmq.DurableQueue()) // Queue 연결을 종료해도 삭제X 설정
		if err != nil { return }

		_, err = brk.Subscribe(topic.ChangeAuthStatusEventTopic, s.ChangeAuthStatus,
			br.Queue(topic.ChangeAuthStatusEventTopic), // Queue 정적 이름 설정
			br.DisableAutoAck(), // Ack를 수동으로 실행하게 설정
			rabbitmq.DurableQueue()) // Queue 연결을 종료해도 삭제X 설정
		if err != nil { return }

		log.Infof("succeed in connecting to broker!! (name: %s | addr: %s)\n",  brk.String(), brk.Address())
		return
	}

	// 서비스 초기화 등록
	service.Init(
		micro.AfterStart(brkHandleFunc),
	)

	// 핸들러 등록 및 서비스 실행
 	if err := authProto.RegisterAuthHandler(service.Server(), h); err != nil {
 		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
