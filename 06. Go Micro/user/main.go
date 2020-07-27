package main

import (
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	br "user/adapter/broker"
	"user/adapter/db"
	"user/closer/broker"
	"user/closer/registry"
	"user/dao"
	"user/handler"
	userProto "user/proto/golang/user"
	"user/subscriber"
	"user/tool/addr"
	"user/tool/env"
	"user/tool/validator"
)

func main() {
	ip := addr.GetLocal().IP
	le := env.GetForLogging()

	cs, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.ConnMysql()
	if err != nil { log.Fatal(err) }
	udc := dao.NewUserDAOCreator(conn)
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }
	rbMQ := br.ConnRabbitMQ()

	sc := &jaegercfg.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1}
	rc := &jaegercfg.ReporterConfig{LogSpans: true, LocalAgentHostPort: "localhost:6831"}
	ujc := jaegercfg.Configuration{ServiceName: "user-service", Sampler: sc, Reporter: rc, Tags: []opentracing.Tag{
		{Key: "environment", Value: le},
		{Key: "host_ip", Value: ip.String()},
		{Key: "service", Value: "userService"},
	}}

	utr, c, err := ujc.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil { log.Fatal(err) }
	defer func() { _ = c.Close() }()

	uh := handler.NewUser(rbMQ, validate, udc, utr)
	us := subscriber.NewUser(rbMQ, validate, udc)

	s := micro.NewService(
		micro.Name("examples.blog.service.user"),
		micro.Version("latest"),
		micro.Broker(rbMQ),
		micro.Transport(grpc.NewTransport()),
	)

	s.Init(
		micro.BeforeStart(broker.RabbitMQInitializer(s.Server(), us)),
		micro.AfterStart(registry.ConsulServiceRegistry(s.Server(), cs)),
		micro.BeforeStop(registry.ConsulServiceDeregistry(s.Server(), cs)),
	)

	if err = userProto.RegisterUserHandler(s.Server(), uh); err != nil {
		log.Fatal(err)
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
