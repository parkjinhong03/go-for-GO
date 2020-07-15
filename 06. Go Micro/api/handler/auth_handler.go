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
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"net/http"
)

type AuthHandler struct {
	cli authProto.AuthService
	validate *validator.Validate
	registry registry.Registry
	breaker *breaker.Breaker
}

func NewAuthHandler(cli authProto.AuthService, validate *validator.Validate,
	registry registry.Registry, breaker *breaker.Breaker) AuthHandler {
	return AuthHandler{
		cli: cli,
		validate: validate,
		registry: registry,
		breaker: breaker,
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

	resp, _ := ah.cli.UserIdDuplicated(ctx, &authProto.UserIdDuplicatedRequest{
		UserId: body.UserId,
	})
	// timeout 처리

	c.JSON(int(resp.Status), resp)
	return
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

	resp, _ := ah.cli.BeforeCreateAuth(ctx, &authProto.BeforeCreateAuthRequest{
		UserId:       body.UserId,
		UserPw:       body.UserPw,
		Name:         body.Name,
		PhoneNumber:  body.PhoneNumber,
		Email:        body.Email,
		Introduction: body.Introduction,
	})
	// timeout 처리

	c.JSON(int(resp.Status), resp)
	return
}