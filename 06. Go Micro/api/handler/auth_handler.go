package handler

import (
	"context"
	"gateway/entity"
	authProto "gateway/proto/golang/auth"
	"gateway/tool/jwt"
	_ "github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"net/http"
	"sync"
)

type AuthHandler struct {
	cli authProto.AuthService
	validate *validator.Validate
	registry registry.Registry
	breaker *breaker.Breaker
	mutex sync.Mutex
}

func NewAuthHandler(cli authProto.AuthService, validate *validator.Validate,
	registry registry.Registry, breaker *breaker.Breaker, mutex sync.Mutex) AuthHandler {
	return AuthHandler{
		cli: cli,
		validate: validate,
		registry: registry,
		breaker: breaker,
		mutex: mutex,
	}
}

func (ah AuthHandler) UserIdDuplicateHandler(c *gin.Context) {
	var body entity.UserIdDuplicate
	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := ah.validate.Struct(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	xReqId := c.GetHeader("X-Request-Id")
	if _, err := uuid.Parse(xReqId); err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", xReqId)
	if ss := c.GetHeader("Unique-Authorization"); ss != "" {
		ctx = metadata.Set(ctx, "Unique-Authorization", ss)
	}

	resp := new(authProto.UserIdDuplicatedResponse)
	reqFunc := func() (err error) {
		opts := []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
		if resp, err = ah.cli.UserIdDuplicated(ctx, &authProto.UserIdDuplicatedRequest{
			UserId: body.UserId,
		}, opts...); err != nil {
			return
		}
		if resp.Status == http.StatusInternalServerError {
			err = errors.InternalServerError("go.micro.client", "internal server error")
		}
		return
	}

	var err error
	switch err = ah.breaker.Run(reqFunc); err {
	case nil:
		notified = false
		c.JSON(int(resp.Status), resp)
	case breaker.ErrBreakerOpen:
		c.Status(http.StatusServiceUnavailable)
		if notified == true { break }
		// 처음으로 열린 차단기라면, 알림 서비스 실행
		notified = true
	default:
		err, ok := err.(*errors.Error)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(int(err.Code))
	}
}

func (ah AuthHandler) UserCreateHandler(c *gin.Context) {
	var body entity.UserCreate
	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := ah.validate.Struct(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	xReqId := c.GetHeader("X-Request-Id")
	if _, err := uuid.Parse(xReqId); err != nil {
		c.Status(http.StatusForbidden)
		return
	}

	ss := c.GetHeader("Unique-Authorization")
	if _, err := jwt.ParseDuplicateCertClaimFromJWT(ss); err != nil || ss == "" {
		c.Status(http.StatusForbidden)
		return
	}

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", xReqId)
	ctx = metadata.Set(ctx, "Unique-Authorization", ss)

	resp := new(authProto.BeforeCreateAuthResponse)
	reqFunc := func() (err error) {
		opts := []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
		if resp, err = ah.cli.BeforeCreateAuth(ctx, &authProto.BeforeCreateAuthRequest{
			UserId:       body.UserId,
			UserPw:       body.UserPw,
			Name:         body.Name,
			PhoneNumber:  body.PhoneNumber,
			Email:        body.Email,
			Introduction: body.Introduction,
		}, opts...); err != nil {
			return
		}
		if resp.Status == http.StatusInternalServerError {
			err = errors.InternalServerError("go.micro.client", "internal server error")
		}
		return
	}

	var err error
	switch err = ah.breaker.Run(reqFunc); err {
	case nil:
		notified = false
		c.JSON(int(resp.Status), resp)
	case breaker.ErrBreakerOpen:
		c.Status(http.StatusServiceUnavailable)
		if notified == true { break }
		// 처음으로 열린 차단기라면, 알림 서비스 실행
		notified = true
	default:
		err, ok := err.(*errors.Error)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(int(err.Code))
	}

	c.JSON(int(resp.Status), resp)
	return
}