package handler

import (
	"context"
	"encoding/json"
	"gateway/entity"
	userProto "gateway/proto/golang/user"
	"gateway/tool/conf"
	"gateway/tool/serializer"
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
	toLogrus serializer.ToLogrusField
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
		toLogrus: serializer.ToLogrusField{},
	}
}

func (uh UserHandler) EmailDuplicateHandler(c *gin.Context) {
	entry := uh.logger.WithFields(logrus.Fields{
		"group":   "handler",
		"segment": "emailDuplicate",
	})

	var body entity.EmailDuplicate
	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		uh.setEntryField(entry, c.Request, body, http.StatusBadRequest, err).Info()
		return
	}

	if err := uh.validate.Struct(&body); err != nil {
		c.Status(http.StatusBadRequest)
		uh.setEntryField(entry, c.Request, body, http.StatusBadRequest, err).Info()
		return
	}

	xReqId := c.GetHeader("X-Request-Id")
	if _, err := uuid.Parse(xReqId); err != nil {
		c.Status(http.StatusForbidden)
		uh.setEntryField(entry, c.Request, body, http.StatusForbidden, err).Info()
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

		uh.setEntryField(entry, c.Request, body, int(resp.Status), err).Info()
	case breaker.ErrBreakerOpen:
		c.Status(http.StatusServiceUnavailable)

		uh.setEntryField(entry, c.Request, body, http.StatusServiceUnavailable, err).Error()
		if uh.notified[emailDuplicateIndex] == true { break }
		// 처음으로 열린 차단기라면, 알림 서비스 실행
		uh.notified[emailDuplicateIndex] = true
	default:
		var code = http.StatusServiceUnavailable
		err, ok := err.(*errors.Error)
		if ok {
			code = int(err.Code)
		}
		c.Status(code)

		uh.setEntryField(entry, c.Request, body, code, err).Warn()
	}
}

func (uh UserHandler) setEntryField(entry *logrus.Entry, r *http.Request, body interface{}, outcome int, err error) *logrus.Entry {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	b, err := json.Marshal(body)
	if err != nil {
		b = []byte{}
	}

	return entry.WithFields(logrus.Fields{
		"json":    string(b),
		"outcome": outcome,
		"error":   errStr,
	}).WithFields(uh.toLogrus.Request(r))
}