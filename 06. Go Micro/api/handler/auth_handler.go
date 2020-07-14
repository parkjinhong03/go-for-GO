package handler

import (
	"context"
	"gateway/entity"
	authProto "gateway/proto/golang/auth"
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
}

func NewAuthHandler(cli authProto.AuthService, validate *validator.Validate, registry registry.Registry) AuthHandler {
	return AuthHandler{
		cli: cli,
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

	c.JSON(int(resp.Status), resp)
}

func (ah AuthHandler) UserCreateHandler(c *gin.Context) {

}