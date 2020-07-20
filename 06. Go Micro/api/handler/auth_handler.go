package handler

import (
	"context"
	"gateway/entity"
	authProto "gateway/proto/golang/auth"
	"gateway/tool/conf"
	"gateway/tool/jwt"
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
	"github.com/uber/jaeger-client-go"
	"net/http"
)

type AuthHandler struct {
	cli      authProto.AuthService
	logger   *logrus.Logger
	validate *validator.Validate
	registry registry.Registry
	tracer	 opentracing.Tracer
	breaker  []*breaker.Breaker
	notified []bool
}

func NewAuthHandler(cli authProto.AuthService, logger *logrus.Logger, validate *validator.Validate,
	registry registry.Registry, tracer opentracing.Tracer, bc conf.BreakerConfig) AuthHandler {

	bk1 := breaker.New(bc.ErrorThreshold, bc.SuccessThreshold, bc.Timeout)
	bk2 := breaker.New(bc.ErrorThreshold, bc.SuccessThreshold, bc.Timeout)

	return AuthHandler{
		cli:      cli,
		logger:   logger,
		validate: validate,
		registry: registry,
		tracer:   tracer,
		breaker:  []*breaker.Breaker{bk1, bk2},
		notified: []bool{false, false},
	}
}

func (ah AuthHandler) UserIdDuplicateHandler(c *gin.Context) {
	var body entity.UserIdDuplicate
	var code int
	xid := c.GetHeader("X-Request-Id")

	ps := ah.tracer.StartSpan(c.Request.URL.Path)
	ps.SetTag("X-Request-Id", xid).SetTag("segment", "userIdDuplicate")
	defer ps.Finish()

	entry := ah.logger.WithField("segment", "userIdDuplicate").WithField("group", "handler")
	entry = entry.WithFields(logrusfield.ForHandleRequest(c.Request, c.ClientIP()))

	if err := c.BindJSON(&body); err != nil {
		code = http.StatusBadRequest
		c.Status(code)
		err = errors.New(apiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", code).LogFields(log.Error(err))
		return
	}

	if err := ah.validate.Struct(&body); err != nil {
		code = http.StatusBadRequest
		c.Status(code)
		err := errors.New(apiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", code).LogFields(log.Error(err))
		return
	}

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", c.GetHeader("X-Request-Id"))
	ctx = metadata.Set(ctx, "Unique-Authorization", c.GetHeader("Unique-Authorization"))

	var resp *authProto.UserIdDuplicatedResponse
	err := ah.breaker[userIdDuplicateIndex].Run(func() (err error) {
		req := body.ToRequestProto()
		opts := []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
		cs := ah.tracer.StartSpan(userIdDuplicate, opentracing.ChildOf(ps.Context())).SetTag("X-Request-Id", xid)
		ctx = metadata.Set(ctx, "Span-Context", cs.Context().(jaeger.SpanContext).String())
		resp, err = ah.cli.UserIdDuplicated(ctx, req, opts...)
		cs.LogFields(log.Object("request", req), log.Object("response", resp), log.Error(err))
		cs.Finish()
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

	return
}

func (ah AuthHandler) UserCreateHandler(c *gin.Context) {
	var body entity.UserCreate
	var code int
	xid := c.GetHeader("X-Request-Id")

	ps := ah.tracer.StartSpan(c.Request.URL.Path)
	ps.SetTag("X-Request-Id", xid).SetTag("segment", "userCreate")
	defer ps.Finish()

	entry := ah.logger.WithField("segment", "userCreate").WithField("group", "handler")
	entry = entry.WithFields(logrusfield.ForHandleRequest(c.Request, c.ClientIP()))

	if err := c.BindJSON(&body); err != nil {
		code = http.StatusBadRequest
		c.Status(code)
		err := errors.New(apiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", code).LogFields(log.Error(err))
		return
	}

	if err := ah.validate.Struct(&body); err != nil {
		code = http.StatusBadRequest
		c.Status(code)
		err := errors.New(apiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", code).LogFields(log.Error(err))
		return
	}

	ss := c.GetHeader("Unique-Authorization")
	if _, err := jwt.ParseDuplicateCertClaimFromJWT(ss); err != nil {
		code = http.StatusForbidden
		c.Status(code)
		err := errors.New(apiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", http.StatusForbidden).LogFields(log.Error(err))
		return
	}

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", xid)
	ctx = metadata.Set(ctx, "Unique-Authorization", ss)

	var resp *authProto.BeforeCreateAuthResponse
	err := ah.breaker[userCreateIndex].Run(func() (err error) {
		opts := []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
		req := body.ToRequestProto()
		cs := ah.tracer.StartSpan(beforeCreateAuth, opentracing.ChildOf(ps.Context())).SetTag("X-Request-Id", xid)
		ctx = metadata.Set(ctx, "Span-Context", cs.Context().(jaeger.SpanContext).String())
		resp, err = ah.cli.BeforeCreateAuth(ctx, req, opts...)
		cs.LogFields(log.Object("request", req), log.Object("response", resp), log.Error(err))
		cs.Finish()
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
		ps.SetTag("status", code)
		return
	}

	code = int(resp.Status)
	c.JSON(code, resp)
	entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
	if resp.Status == http.StatusInternalServerError {
		entry.Warn()
	} else {
		entry.Info()
	}
	ps.SetTag("status", resp.Status)
	return
}
