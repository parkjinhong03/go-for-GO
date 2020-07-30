package main

import (
	"auth/adapter/broker"
	"auth/adapter/db"
	brokercloser "auth/closer/broker"
	registrycloser "auth/closer/registry"
	sqlcloser "auth/closer/sql"
	"auth/dao"
	"auth/handler"
	customchecker "auth/plugin/checker"
	authProto "auth/proto/golang/auth"
	"auth/subscriber"
	"auth/tool/addr"
	"auth/tool/env"
	"auth/tool/validator"
	"github.com/InVisionApp/go-health/v2"
	"github.com/InVisionApp/go-health/v2/checkers"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"time"
)

func main() {
	le := env.GetLoggingEnv()
	ip := addr.GetLocalAddr().IP

	// Message Broker 객체 생성
	rbMQ := broker.ConnRabbitMQ()
	if err := rbMQ.Connect(); err != nil { log.Fatal(err) }

	// Service Discovery 객체 생성
	cs, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	// Persistance Layer 객체 생성
	conn, err := db.ConnMysql()
	if err != nil {
		log.Fatalf("unable to connect mysql server, err: %v\n", err)
	}
	adc := dao.NewAuthDAOCreator(conn)

	// 유효성 검사 객체 생성
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }

	// Jaeger 설정 객체 생성
	sc := &jaegercfg.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1}
	rc := &jaegercfg.ReporterConfig{LogSpans: true, LocalAgentHostPort: "localhost:6831"}
	ts := []opentracing.Tag{
		{Key: "environment", Value: le},
		{Key: "host_ip", Value: ip.String()},
		{Key: "service", Value: "authService"},
	}
	ajc := jaegercfg.Configuration{ServiceName: "auth-service", Sampler: sc, Reporter: rc, Tags: ts}

	// Tracer 실행 및 객체 생성
	atr, c, err := ajc.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil { log.Fatal(err) }
	defer func() { _ = c.Close() }()

	// 이벤트 및 rpc 핸들러 객체 생성
	as := subscriber.NewAuth(rbMQ, adc, validate, atr)
	ah := handler.NewAuth(rbMQ, adc, validate, atr)

	// 서비스 생성
	s := micro.NewService(
		micro.Name("examples.blog.service.auth"),
		micro.Version("latest"),
		micro.Broker(rbMQ),
		micro.Transport(grpc.NewTransport()),
	)

	// 서비스 초기화 등록
	s.Init(
		micro.BeforeStart(brokercloser.RabbitMQInitializer(s.Server(), as)),
		micro.AfterStart(registrycloser.ConsulServiceRegistry(s.Server(), cs)),
		micro.AfterStop(registrycloser.ConsulServiceDeregistry(s.Server(), cs)),
	)

	// rpc 핸들러 등록
 	if err := authProto.RegisterAuthHandler(s.Server(), ah); err != nil {
 		log.Fatal(err)
	}

	// health check 객체 생성
	h := health.New()

	// DB checker 객체 생성
	sqlc, err := checkers.NewSQL(&checkers.SQLConfig{ Pinger: conn.DB() })
	if err != nil { log.Fatal(err) }
	sqlh := &health.Config{
		Name:       "SQL-Checker",
		Checker:    sqlc,
		Interval:   time.Second * 5,
		OnComplete: sqlcloser.TTLCheckHandler(s.Server(), cs),
	}

	// Broker Checker 객체 생성
	brc, err := customchecker.NewBroker(rbMQ)
	if err != nil { log.Fatal(err) }
	brh := &health.Config{
		Name:       "Broker-Checker",
		Checker:    brc,
		Interval:   time.Second * 5,
		OnComplete: brokercloser.TTLCheckHandler(s.Server(), cs),
	}

	// DB, MQ health checker 등록
	if err = h.AddChecks([]*health.Config{sqlh, brh}); err != nil {
		log.Fatal(err)
	}

	// health check 실행
	if err = h.Start(); err != nil {
		log.Fatal(err)
	}

	// 서비스 실행
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
