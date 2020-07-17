package main

import (
	"gateway/handler"
	"gateway/middleware"
	authProto "gateway/proto/golang/auth"
	userProto "gateway/proto/golang/user"
	"gateway/tool/conf"
	"gateway/tool/validator"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/sirupsen/logrus"
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

	fields := logrus.Fields{
		"environment": getEnvironment(),
		"host_ip":     getLocalAddr().IP,
		"service":     "apiGateway",
	}

	// 훅(파일 로깅) 생성
	ahk := logrustash.New(af, logrustash.DefaultFormatter(fields))
	uhk := logrustash.New(uf, logrustash.DefaultFormatter(fields))

	// 로깅 객체 생성
	al := logrus.New()
	ul := logrus.New()

	// 훅(파일 로깅) 등록
	al.Hooks.Add(ahk)
	ul.Hooks.Add(uhk)

	// rpc 클라이언트 객체 생성
	opts := []client.Option{client.Registry(cs)}
	ac := authProto.NewAuthService(AuthServiceName, grpc.NewClient(opts...))
	uc := userProto.NewUserService(UserServiceName, grpc.NewClient(opts...))

	// 핸들러 객체 생성
	ah := handler.NewAuthHandler(ac, al, v, cs, bc)
	uh := handler.NewUserHandler(uc, ul, v, cs, bc)

	// 핸들러 라우팅
	router := gin.Default()
	v1 := router.Group("/v1")
	v1.Use(middleware.Correlation())
	{
		v1.GET("/user-ids/duplicate", ah.UserIdDuplicateHandler)
		v1.POST("/users", ah.UserCreateHandler)
	}
	{
		v1.GET("/emails/duplicate", uh.EmailDuplicateHandler)
	}

	// api gateway 실행
	if err := router.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}

func getLocalAddr() *net.UDPAddr {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil { log.Fatal(err) }
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