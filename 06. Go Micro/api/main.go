package main

import (
	"gateway/handler"
	"gateway/middleware"
	authProto "gateway/proto/golang/auth"
	userProto "gateway/proto/golang/user"
	"gateway/tool/validator"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"log"
	"sync"
	"time"
)

const (
	DefaultErrorThreshold = 3
	DefaultSuccessThreshold = 3
	DefaultTimeout = time.Minute
)

func main() {
	cs := consul.NewRegistry(registry.Addrs("http://localhost:8500"))
	v, err := validator.New()
	if err != nil { log.Fatal(err) }
	bk := breaker.New(DefaultErrorThreshold, DefaultSuccessThreshold, DefaultTimeout)
	mtx := sync.Mutex{}

	opts := []client.Option{client.Registry(cs)}
	ac := authProto.NewAuthService("examples.blog.service.auth", grpc.NewClient(opts...))
	uc := userProto.NewUserService("examples.blog.service.user", grpc.NewClient(opts...))

	ah := handler.NewAuthHandler(ac, v, cs, bk, mtx)
	uh := handler.NewUserHandler(uc, v, cs, bk, mtx)

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

	if err := router.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}