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
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2"
	br "github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	env := getEnvironment()
	ip := getLocalAddr().IP

	// Message Broker 객체 생성
	rbMQ := broker.ConnRabbitMQ()
	//cs := consul.NewRegistry(registry.Addrs("http:localhost:8500"))

	// Service Discovery 객체 생성
	dc := api.DefaultConfig()
	dc.Address = ip.String() + ":8500"
	cs, err := api.NewClient(dc)
	if err != nil {
		log.Fatal(err)
	}

	// 의존성 주입을 위한 객체 생성
	conn, err := db.ConnMysql()
	if err != nil {
		log.Fatalf("unable to connect mysql server, err: %v\n", err)
	}
	adc := dao.NewAuthDAOCreator(conn)
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }

	sc := &jaegercfg.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1}
	rc := &jaegercfg.ReporterConfig{LogSpans: true, LocalAgentHostPort: "localhost:6831"}

	ajc := jaegercfg.Configuration{ServiceName: "auth-service", Sampler: sc, Reporter: rc, Tags: []opentracing.Tag{
		{Key: "environment", Value: env},
		{Key: "host_ip", Value: ip.String()},
		{Key: "service", Value: "authService"},
	}}
	atr, c, err := ajc.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil { log.Fatal(err) }
	defer func() { _ = c.Close() }()

	// 서비스 생성
	s := micro.NewService(
		micro.Name("examples.blog.service.auth"),
		micro.Version("latest"),
		micro.Broker(rbMQ),
		//micro.Registry(cs),
		micro.Transport(grpc.NewTransport()),
	)

	// 이벤트 및 rpc 핸들러 객체 생성
	as := subscriber.NewAuth(rbMQ, adc, validate, atr)
	// mq := service.Options().Broker
	ah := handler.NewAuth(rbMQ, adc, validate, atr)

	// Broker 초기화 핸들러 함수 생성
	brf := func() (err error) {
		brk := s.Options().Broker

		if err = brk.Connect(); err != nil { return }
		_, err = brk.Subscribe(topic.CreateAuthEventTopic, as.CreateAuth,
			br.Queue(topic.CreateAuthEventTopic), // Queue 정적 이름 설정
			br.DisableAutoAck(), // Ack를 수동으로 실행하게 설정
			rabbitmq.DurableQueue()) // Queue 연결을 종료해도 삭제X 설정
		if err != nil { return }

		_, err = brk.Subscribe(topic.ChangeAuthStatusEventTopic, as.ChangeAuthStatus,
			br.Queue(topic.ChangeAuthStatusEventTopic), // Queue 정적 이름 설정
			br.DisableAutoAck(), // Ack를 수동으로 실행하게 설정
			rabbitmq.DurableQueue()) // Queue 연결을 종료해도 삭제X 설정
		if err != nil { return }

		log.Infof("succeed in connecting to broker!! (name: %s | addr: %s)\n",  brk.String(), brk.Address())
		return
	}

	// service discovery(consul)에 서비스 등록 함수 생성
	csf := func() (err error) {
		ps := strings.Split(s.Server().Options().Address, ":")[3]
		port, err := strconv.Atoi(ps)
		if err != nil { log.Fatal(err) }
		sid := s.Server().Options().Name + "-" + s.Server().Options().Id
		cid := "service:" + sid

		asr := &api.AgentServiceRegistration{
			ID:      sid,
			Name:    s.Server().Options().Name,
			Port:    port,
			Address: ip.String(),
		}
		err = cs.Agent().ServiceRegister(asr)
		if err != nil { log.Fatal(err) }

		asc := api.AgentServiceCheck{
			Name:   s.Server().Options().Name,
			Status: "passing",
			TTL:    "8640s",
		}
		acr := &api.AgentCheckRegistration{
			ID:                cid,
			Name:              fmt.Sprintf("service '%s' check", s.Server().Options().Name),
			ServiceID:         sid,
			AgentServiceCheck: asc,
		}
		err = cs.Agent().CheckRegister(acr)
		if err != nil { log.Fatal(err) }

		log.Infof("succeed to registry service and check to consul!! (service id: %s | check id: %s)\n", sid, cid)
		return
	}

	// 서비스 초기화 등록
	s.Init(
		micro.AfterStart(brf),
		micro.AfterStart(csf),
	)

	// 핸들러 등록 및 서비스 실행
 	if err := authProto.RegisterAuthHandler(s.Server(), ah); err != nil {
 		log.Fatal(err)
	}
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

func getLocalAddr() *net.UDPAddr {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil { log.Fatal(err) }
	defer func() { _ = conn.Close() } ()
	return conn.LocalAddr().(*net.UDPAddr)
}

func getEnvironment() (env string) {
	env = os.Getenv("ENV")
	switch env {
	case "DEV":
		env = "development"
	case "PROD":
		env = "production"
	default:
		env = "development"
	}
	return
}