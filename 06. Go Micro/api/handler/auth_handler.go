package handler

import (
	"context"
	"gateway/entity"
	authProto "gateway/proto/golang/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"log"
	"net/http"
)

type AuthHandler struct {
	ac authProto.AuthService
	validate *validator.Validate
	registry registry.Registry
}

func NewAuthHandler(ac authProto.AuthService, validate *validator.Validate, registry registry.Registry) AuthHandler {
	return AuthHandler{
		ac: ac,
		validate: validate,
		registry: registry,
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
	if xReqId == "" {
		c.Status(http.StatusForbidden)
	}

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", xReqId)
	resp, err := ah.ac.UserIdDuplicated(ctx, &authProto.UserIdDuplicatedRequest{
		UserId: body.UserId,
	})

	if err != nil {
		log.Fatal(err)
	}
	c.Status(int(resp.Status))
}

func (ah AuthHandler) UserCreateHandler(c *gin.Context) {

}