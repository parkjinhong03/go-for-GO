package main

import (
	"gateway/handler"
	md "gateway/middleware"
	authProto "gateway/proto/golang/auth"
	userProto "gateway/proto/golang/user"
	"gateway/tool/conf"
	"gateway/tool/validator"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	clientgrpc "github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	transportgrpc "github.com/micro/go-micro/v2/transport/grpc"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"log"
	"net"
	"os"
	"time"
)

const (
	DefaultErrorThreshold = 3
	DefaultSuccessThreshold = 3
	DefaultTimeout = time.Minute
	AuthLogFilePath = "/Users/parkjinhong/log/gateway/auth.log"
	UserLogFilePath = "/Users/parkjinhong/log/gateway/user.log"
	AuthServiceName = "examples.blog.service.auth"
	UserServiceName = "examples.blog.service.user"
)

func main() {
	// service discovery 의존성 객체 생성
	cs := consul.NewRegistry(registry.Addrs("http://localhost:8500"))

	// Round Robin Selector 의존성 객체 생성
	rrs := selector.NewSelector(selector.SetStrategy(selector.RoundRobin), selector.Registry(cs))

	rrs.Select()
	// 유효성 검사 의존성 객체 생성
	v, err := validator.New()
	if err != nil { log.Fatal(err) }

	// 회로 차단기의 설정값을 포함하는 객체 생성
	bc := conf.BreakerConfig{
		ErrorThreshold:   DefaultErrorThreshold,
		SuccessThreshold: DefaultSuccessThreshold,
		Timeout:          DefaultTimeout,
	}

	// log 파일 연결 객체 생성
	af, err := os.OpenFile(AuthLogFilePath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil { log.Fatal(err) }
	uf, err := os.OpenFile(UserLogFilePath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil { log.Fatal(err) }

	env := getEnvironment()
	addr := getLocalAddr()

	// 훅(파일 로깅) 생성
	ahk := logrustash.New(af, logrustash.DefaultFormatter(logrus.Fields{
		"environment": env,
		"host_ip":     addr.IP,
		"service":     "apiGateway",
	}))
	uhk := logrustash.New(uf, logrustash.DefaultFormatter(logrus.Fields{
		"environment": env,
		"host_ip":     addr.IP,
		"service":     "apiGateway",
	}))

	// 로깅 객체 생성
	al := logrus.New()
	ul := logrus.New()

	// 훅(파일 로깅) 등록
	al.Hooks.Add(ahk)
	ul.Hooks.Add(uhk)

	// jaeger tracer 생성을 위한 설정 객체 생성
	sc := &jaegercfg.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1}
	rc := &jaegercfg.ReporterConfig{LogSpans: true, LocalAgentHostPort: "localhost:6831"}

	ajc := jaegercfg.Configuration{ServiceName: "auth-service", Sampler: sc, Reporter: rc, Tags: []opentracing.Tag{
		{Key: "environment", Value: env},
		{Key: "host_ip", Value: addr.IP},
		{Key: "service", Value: "authService"},
	}}
	atr, c, err := ajc.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil { log.Fatal(err) }
	opentracing.SetGlobalTracer(atr) // 이게 맞는건가?
	defer func() { _ = c.Close() }()

	ujc := jaegercfg.Configuration{ServiceName: "user-service", Sampler: sc, Reporter: rc, Tags: []opentracing.Tag{
		{Key: "environment", Value: env},
		{Key: "host_ip", Value: addr.IP},
		{Key: "service", Value: "userService"},
	}}
	utr, c, err := ujc.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil { log.Fatal(err) }
	opentracing.SetGlobalTracer(utr) // 이게 맞는건가?
	defer func() { _ = c.Close() }()

	// rpc 클라이언트 객체 생성
	opts := []client.Option{client.Registry(cs), client.Transport(transportgrpc.NewTransport()), client.Selector(rrs)}
	ac := authProto.NewAuthService(AuthServiceName, clientgrpc.NewClient(opts...))
	uc := userProto.NewUserService(UserServiceName, clientgrpc.NewClient(opts...))

	// 핸들러 객체 생성
	ah := handler.NewAuthHandler(ac, al, v, cs, atr, bc)
	uh := handler.NewUserHandler(uc, ul, v, cs, utr, bc)

	// 핸들러 라우팅
	router := gin.Default()
	v1 := router.Group("/v1")
	v1.Use(md.Correlator())

	ar := v1.Group("/")
	{
		ar.GET("/user-ids/duplicate", ah.UserIdDuplicateHandler)
		ar.POST("/users", ah.UserCreateHandler)
	}

	ur := v1.Group("/")
	{
		ur.GET("/emails/duplicate", uh.EmailDuplicateHandler)
	}

	// api gateway 실행
	if err := router.Run(":8000"); err != nil {
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