package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	. "net"
	"os"
	br "user/adapter/broker"
	"user/adapter/db"
	"user/dao"
	"user/handler"
	userProto "user/proto/golang/user"
	"user/subscriber"
	"user/tool/validator"
)

func main() {
	conn, err := db.ConnMysql()
	if err != nil { log.Fatal(err) }
	udc := dao.NewUserDAOCreator(conn)
	validate, err := validator.New()
	if err != nil { log.Fatal(err) }
	rbMQ := br.ConnRabbitMQ()
	cs := consul.NewRegistry(registry.Addrs("http://localhost:8500"))

	sc := &jaegercfg.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1}
	rc := &jaegercfg.ReporterConfig{LogSpans: true, LocalAgentHostPort: "localhost:6831"}
	ujc := jaegercfg.Configuration{ServiceName: "user-service", Sampler: sc, Reporter: rc, Tags: []opentracing.Tag{
		{Key: "environment", Value: getEnvironment()},
		{Key: "host_ip", Value: getLocalAddr().IP},
		{Key: "service", Value: "userService"},
	}}

	utr, c, err := ujc.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil { log.Fatal(err) }
	defer func() { _ = c.Close() }()

	h := handler.NewUser(rbMQ, validate, udc, utr)
	s := subscriber.NewUser(rbMQ, validate, udc)

	service := micro.NewService(
		micro.Name("examples.blog.service.user"),
		micro.Version("latest"),
		micro.Broker(rbMQ),
		micro.Registry(cs),
		micro.Transport(grpc.NewTransport()),
	)

	brkHandler := func() error {
		brk := service.Options().Broker
		if err := brk.Connect(); err != nil { log.Fatal(err) }
		options := []broker.SubscribeOption{broker.Queue(subscriber.CreateUserEventTopic), broker.DisableAutoAck(), rabbitmq.DurableQueue()}
		if _, err := brk.Subscribe(subscriber.CreateUserEventTopic, s.CreateUser, options...); err != nil { log.Fatal(err) }
		log.Infof("succeed in connecting to broker!! (name: %s | addr: %s)\n",  brk.String(), brk.Address())
		return nil
	}

	service.Init(micro.AfterStart(brkHandler))

	if err = userProto.RegisterUserHandler(service.Server(), h); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func getLocalAddr() *UDPAddr {
	conn, err := Dial("udp", "8.8.8.8:80")
	if err != nil { log.Fatal(err) }
	defer func() { _ = conn.Close() } ()
	return conn.LocalAddr().(*UDPAddr)
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