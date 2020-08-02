package main

import (
	"encoding/json"
	"fmt"
	registryentity "gateway/entity/registry"
	"gateway/handler"
	md "gateway/middleware"
	authProto "gateway/proto/golang/auth"
	userProto "gateway/proto/golang/user"
	"gateway/tool/conf"
	"gateway/tool/validator"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2/client"
	clientgrpc "github.com/micro/go-micro/v2/client/grpc"
	transportgrpc "github.com/micro/go-micro/v2/transport/grpc"
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
	Port = 8000
	DefaultErrorThreshold = 3
	DefaultSuccessThreshold = 3
	DefaultTimeout = time.Minute
	AuthLogFileDirectory = "/Users/parkjinhong/log/gateway"
	UserLogFileDirectory = "/Users/parkjinhong/log/gateway"
	AuthLogFilePath = "/Users/parkjinhong/log/gateway/auth.log"
	UserLogFilePath = "/Users/parkjinhong/log/gateway/user.log"
	AuthServiceName = "examples.blog.service.auth"
	UserServiceName = "examples.blog.service.user"
)

func main() {
	env := getEnvironment()
	addr := getLocalAddr()

	if os.Getenv("CONSUL_ADDRESS") == "" {
		log.Fatal("Please set CONSUL_ADDRESS in env")
	}
	cfg := api.DefaultConfig()
	cfg.Address = os.Getenv(os.Getenv("CONSUL_ADDRESS"))
	cs, err := api.NewClient(cfg)
	if err != nil { log.Fatal(err) }

	// 유효성 검사 의존성 객체 생성
	v, err := validator.New()
	if err != nil { log.Fatal(err) }

	// 회로 차단기의 설정값을 포함하는 객체 생성
	bc := conf.BreakerConfig{
		ErrorThreshold:   DefaultErrorThreshold,
		SuccessThreshold: DefaultSuccessThreshold,
		Timeout:          DefaultTimeout,
	}

	// log 디렉토리 초기화
	if _, err := os.Stat(AuthLogFileDirectory); os.IsNotExist(err) {
		if err = os.MkdirAll(AuthLogFileDirectory, os.ModePerm); err != nil { log.Fatal(err) }
	}
	if _, err := os.Stat(UserLogFileDirectory); os.IsNotExist(err) {
		if err = os.MkdirAll(UserLogFileDirectory, os.ModePerm); err != nil { log.Fatal(err) }
	}

	// log 파일 연결 객체 생성
	af, err := os.OpenFile(AuthLogFilePath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil { log.Fatal(err) }
	uf, err := os.OpenFile(UserLogFilePath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil { log.Fatal(err) }

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

	// jaeger address를 얻기 위한 consul KV 파싱
	kp, _, err := cs.KV().Get("jaeger", nil)
	if err != nil { log.Fatal(err) }
	var body registryentity.Jaeger
	if err := json.Unmarshal(kp.Value, &body); err != nil { log.Fatal(err) }

	// jaeger tracer 생성을 위한 설정 객체 생성
	sc := &jaegercfg.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1}
	rc := &jaegercfg.ReporterConfig{LogSpans: true, LocalAgentHostPort: body.Addr}
	ajc := jaegercfg.Configuration{ServiceName: "auth-service", Sampler: sc, Reporter: rc, Tags: []opentracing.Tag{
		{Key: "environment", Value: env},
		{Key: "host_ip", Value: addr.IP},
		{Key: "service", Value: "authService"},
	}}
	ujc := jaegercfg.Configuration{ServiceName: "user-service", Sampler: sc, Reporter: rc, Tags: []opentracing.Tag{
		{Key: "environment", Value: env},
		{Key: "host_ip", Value: addr.IP},
		{Key: "service", Value: "userService"},
	}}

	// tracer 시작 및 객체 생성
	atr, c, err := ajc.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil { log.Fatal(err) }
	opentracing.SetGlobalTracer(atr)
	defer func() { _ = c.Close() }()
	utr, c, err := ujc.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil { log.Fatal(err) }
	opentracing.SetGlobalTracer(utr)
	defer func() { _ = c.Close() }()

	// rpc 클라이언트 객체 생성
	opts := []client.Option{client.Transport(transportgrpc.NewTransport())}
	ac := authProto.NewAuthService(AuthServiceName, clientgrpc.NewClient(opts...))
	uc := userProto.NewUserService(UserServiceName, clientgrpc.NewClient(opts...))

	// 핸들러 객체 생성
	ah := handler.NewAuthHandler(ac, al, v, cs, atr, bc)
	uh := handler.NewUserHandler(uc, ul, v, cs, utr, bc)

	// 핸들러 라우팅
	router := gin.Default()
	v1 := router.Group("/v1")
	v1.Use(md.Correlator())
	v1.Use(md.Throttler().Throttling)

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
	if err := router.Run(fmt.Sprintf(":%d", Port)); err != nil {
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