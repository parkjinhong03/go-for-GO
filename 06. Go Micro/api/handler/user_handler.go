package handler

import (
	userProto "gateway/proto/golang/user"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/registry"
)

type UserHandler struct {
	uc userProto.UserService
	validate *validator.Validate
	registry registry.Registry
}

func NewUserHandler(uc userProto.UserService, validate *validator.Validate, registry registry.Registry) UserHandler {
	return UserHandler{
		uc: uc,
		validate: validate,
		registry: registry,
	}
}

func (u UserHandler) EmailDuplicateHandler(c *gin.Context) {

}