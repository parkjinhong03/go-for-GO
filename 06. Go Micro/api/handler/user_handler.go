package handler

import (
	"context"
	"gateway/entity"
	userProto "gateway/proto/golang/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"net/http"
)

type UserHandler struct {
	cli userProto.UserService
	validate *validator.Validate
	registry registry.Registry
}

func NewUserHandler(cli userProto.UserService, validate *validator.Validate, registry registry.Registry) UserHandler {
	return UserHandler{
		cli: cli,
		validate: validate,
		registry: registry,
	}
}

func (u UserHandler) EmailDuplicateHandler(c *gin.Context) {
	var body entity.EmailDuplicate
	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := u.validate.Struct(&body); err != nil {
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

	resp, _ := u.cli.EmailDuplicated(ctx, &userProto.EmailDuplicatedRequest{
		Email: body.Email,
	})

	c.JSON(int(resp.Status), resp)
}