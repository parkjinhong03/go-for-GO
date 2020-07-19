package handler

import (
	"context"
	"fmt"
	"gateway/entity"
	userProto "gateway/proto/golang/user"
	"gateway/tool/conf"
	"gateway/tool/logrusfield"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	var code int
	xid := c.GetHeader("X-Request-Id")
	entry := uh.logger.WithField("segment", "emailDuplicate")
	entry = entry.WithFields(logrusfield.ForHandleRequest(c.Request, c.ClientIP()))

	if v, ok := c.Get("error"); ok {
		code = http.StatusInternalServerError
		err := errors.New(apiGateway, fmt.Sprintf("some error occurs in middleware, err: %v", v.(error)), int32(code))
		entry = entry.WithField("group", "middleware").WithFields(logrusfield.ForReturn(body, code, err))
		entry.Warn()
		return
	}

	v, ok := c.Get("tracer")
	if !ok {
		code = http.StatusInternalServerError
		err := errors.New(apiGateway, "there isn't tracer in *gin.Context", int32(code))
		entry = entry.WithField("group", "middleware").WithFields(logrusfield.ForReturn(body, code, err))
		entry.Warn()
		return
	}

	tr := v.(opentracing.Tracer)
	ps := tr.StartSpan(c.Request.URL.Path)
	ps.SetTag("X-Request-Id", xid).SetTag("segment", "emailDuplicate")
	defer ps.Finish()
	entry = entry.WithField("group", "handler")

	if err := c.BindJSON(&body); err != nil {
		code = http.StatusBadRequest
		c.Status(code)
		err := errors.New(apiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", code).LogFields(log.Error(err))
		return
	}

	if err := uh.validate.Struct(&body); err != nil {
		code = http.StatusBadRequest
		c.Status(code)
		err := errors.New(apiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", code).LogFields(log.Error(err))
		return
	}

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", xid)
	ctx = metadata.Set(ctx, "Unique-Authorization", c.GetHeader("Unique-Authorization"))

	var cs opentracing.Span
	var resp *userProto.EmailDuplicatedResponse
	err := uh.breaker[emailDuplicateIndex].Run(func() (err error) {
		cs = tr.StartSpan(emailDuplicate, opentracing.ChildOf(ps.Context())).SetTag("X-Request-Id", xid)
		req := &userProto.EmailDuplicatedRequest{Email: body.Email}
		opts := []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
		resp, err = uh.cli.EmailDuplicated(ctx, req, opts...)
		cs.LogFields(log.Object("request", req), log.Object("response", resp))
		return
	})

	if err == breaker.ErrBreakerOpen {
		code = http.StatusServiceUnavailable
		c.Status(code)
		err := errors.New(userClient, breaker.ErrBreakerOpen.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Error()
		ps.SetTag("status", code).LogFields(log.Error(err))
		return
	}

	if err != nil {
		code = http.StatusInternalServerError
		if err, ok := err.(*errors.Error); ok { code = int(err.Code) }
		c.Status(code)
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Error()
		ps.SetTag("status", code).LogFields(log.Error(err))
		cs.Finish()
		return
	}

	code = int(resp.Status)
	c.JSON(code, resp)
	entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
	if code == http.StatusInternalServerError {
		entry.Warn()
	} else {
		entry.Info()
	}
	ps.SetTag("status", code)
	cs.Finish()

	return
}
