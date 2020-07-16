package handler

import (
	"context"
	"gateway/entity"
	userProto "gateway/proto/golang/user"
	"gateway/tool/conf"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type UserHandler struct {
	cli      userProto.UserService
	logger 	 *logrus.Logger
	validate *validator.Validate
	registry registry.Registry
	breaker  []*breaker.Breaker
	mutex    sync.Mutex
	notified []bool
}

func NewUserHandler(cli userProto.UserService, logger *logrus.Logger, validate *validator.Validate,
	registry registry.Registry, bc conf.BreakerConfig) UserHandler {

	bk := breaker.New(bc.ErrorThreshold, bc.SuccessThreshold, bc.Timeout)

	return UserHandler{
		cli:      cli,
		logger:   logger,
		validate: validate,
		registry: registry,
		breaker:  []*breaker.Breaker{bk},
		mutex:    sync.Mutex{},
		notified: []bool{false},
	}
}

func (uh UserHandler) EmailDuplicateHandler(c *gin.Context) {
	var body entity.EmailDuplicate
	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := uh.validate.Struct(&body); err != nil {
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

	resp := new(userProto.EmailDuplicatedResponse)
	reqFunc := func() (err error) {
		opts := []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
		if resp, err = uh.cli.EmailDuplicated(ctx, &userProto.EmailDuplicatedRequest{
			Email: body.Email,
		}, opts...); err != nil {
			return
		}
		if resp.Status == http.StatusInternalServerError {
			err = errors.InternalServerError("go.micro.client", "internal server error")
		}
		return
	}

	var err error
	switch err = uh.breaker[emailDuplicateIndex].Run(reqFunc); err {
	case nil:
		uh.notified[emailDuplicateIndex] = false
		c.JSON(int(resp.Status), resp)
	case breaker.ErrBreakerOpen:
		c.Status(http.StatusServiceUnavailable)
		if uh.notified[emailDuplicateIndex] == true { break }
		// 처음으로 열린 차단기라면, 알림 서비스 실행
		uh.notified[emailDuplicateIndex] = true
	default:
		err, ok := err.(*errors.Error)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(int(err.Code))
	}
}