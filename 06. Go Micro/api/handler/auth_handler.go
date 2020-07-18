package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/entity"
	authProto "gateway/proto/golang/auth"
	"gateway/tool/conf"
	"gateway/tool/jwt"
	"gateway/tool/serializer"
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

type AuthHandler struct {
	cli      authProto.AuthService
	logger   *logrus.Logger
	validate *validator.Validate
	registry registry.Registry
	breaker  []*breaker.Breaker
	mutex    sync.Mutex
	notified []bool
	toLogrus serializer.ToLogrusField
}

func NewAuthHandler(cli authProto.AuthService, logger *logrus.Logger, validate *validator.Validate,
	registry registry.Registry, bc conf.BreakerConfig) AuthHandler {

	bk1 := breaker.New(bc.ErrorThreshold, bc.SuccessThreshold, bc.Timeout)
	bk2 := breaker.New(bc.ErrorThreshold, bc.SuccessThreshold, bc.Timeout)

	return AuthHandler{
		cli:      cli,
		logger:   logger,
		validate: validate,
		registry: registry,
		breaker:  []*breaker.Breaker{bk1, bk2},
		mutex:    sync.Mutex{},
		notified: []bool{false, false},
	}
}

func (ah AuthHandler) UserIdDuplicateHandler(c *gin.Context) {
	var body entity.UserIdDuplicate
	entry := ah.logger.WithField("segment", "userIdDuplicate")

	if v, ok := c.Get("error"); ok {
		c.Status(http.StatusInternalServerError)
		entry = entry.WithField("group", "middleware")
		err := errors.New(apiGateway, fmt.Sprintf("some error occurs in middleware, err: %v", v.(error)), http.StatusInternalServerError)
		ah.setEntryField(entry, c.Request, body, http.StatusInternalServerError, err).Warn()
		return
	}

	v, ok := c.Get("tracer")
	if !ok {
		c.Status(http.StatusInternalServerError)
		entry = entry.WithField("group", "middleware")
		err := errors.New(apiGateway, "there isn't tracer in *gin.Context", http.StatusInternalServerError)
		ah.setEntryField(entry, c.Request, body, http.StatusInternalServerError, err).Warn()
		return
	}

	tr := v.(opentracing.Tracer)
	ps := tr.StartSpan(c.Request.URL.Path)
	ps.SetTag("X-Request-Id", c.GetHeader("X-Request-Id")).SetTag("segment", "userIdDuplicate")
	defer ps.Finish()
	entry = entry.WithField("group", "handler")

	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		err = errors.New(apiGateway, err.Error(), http.StatusBadRequest)
		ah.setEntryField(entry, c.Request, body, http.StatusBadRequest, err).Info()
		ps.SetTag("status", http.StatusBadRequest).LogFields(log.Error(err))
		return
	}

	if err := ah.validate.Struct(&body); err != nil {
		c.Status(http.StatusBadRequest)
		err := errors.New(apiGateway, err.Error(), http.StatusBadRequest)
		ah.setEntryField(entry, c.Request, body, http.StatusBadRequest, err).Info()
		ps.SetTag("status", http.StatusBadRequest).LogFields(log.Error(err))
		return
	}

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", c.GetHeader("X-Request-Id"))
	ctx = metadata.Set(ctx, "Unique-Authorization", c.GetHeader("Unique-Authorization"))

	var cs opentracing.Span
	var resp *authProto.UserIdDuplicatedResponse
	err := ah.breaker[userIdDuplicateIndex].Run(func() (err error) {
		cs = tr.StartSpan(userIdDuplicate, opentracing.ChildOf(ps.Context())).SetTag("X-Request-Id", c.GetHeader("X-Request-Id"))
		opts := []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
		req := &authProto.UserIdDuplicatedRequest{UserId: body.UserId}
		resp, err = ah.cli.UserIdDuplicated(ctx, req, opts...)
		cs.LogFields(log.Object("request", req), log.Object("response", resp))
		return
	})

	if err == breaker.ErrBreakerOpen {
		c.Status(http.StatusServiceUnavailable)
		err := errors.New(userClient, breaker.ErrBreakerOpen.Error(), http.StatusServiceUnavailable)
		ah.setEntryField(entry, c.Request, body, http.StatusServiceUnavailable, err).Error()
		ps.SetTag("status", http.StatusServiceUnavailable).LogFields(log.Error(err))
		return
	}

	if err != nil {
		var code = http.StatusInternalServerError
		if err, ok := err.(*errors.Error); ok { code = int(err.Code) }
		c.Status(code)
		ah.setEntryField(entry, c.Request, body, code, err).Error()
		ps.SetTag("status", code).LogFields(log.Error(err))
		cs.Finish()
		return
	}

	c.JSON(int(resp.Status), resp)
	entry = ah.setEntryField(entry, c.Request, body, int(resp.Status), err)
	if resp.Status == http.StatusInternalServerError {
		entry.Warn()
	} else {
		entry.Info()
	}

	ps.SetTag("status", resp.Status)
	cs.Finish()

	return
}

func (ah AuthHandler) UserCreateHandler(c *gin.Context) {
	var body entity.UserCreate
	entry := ah.logger.WithField("segment", "userCreate")

	if v, ok := c.Get("error"); ok {
		var code int32 = http.StatusInternalServerError
		c.Status(int(code))
		entry = entry.WithField("group", "middleware")
		err := errors.New(apiGateway, fmt.Sprintf("some error occurs in middlewares, err: %v\n", v.(error)), code)
		ah.setEntryField(entry, c.Request, body, int(code), err).Warn()
		return
	}

	v, ok := c.Get("tracer")
	if !ok {
		var code = http.StatusInternalServerError
		c.Status(code)
		entry = entry.WithField("group", "middleware")
		err := errors.New(apiGateway, "there isn't tracer in *gin.Context", int32(code))
		ah.setEntryField(entry, c.Request, body, code, err)
		return
	}

	xid := c.GetHeader("X-Request-Id")
	tr := v.(opentracing.Tracer)
	ps := tr.StartSpan(c.Request.URL.Path)
	defer ps.Finish()
	ps.SetTag("X-Request-Id", xid).SetTag("segment", "userCreate")
	entry = entry.WithField("group", "handler")

	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		err := errors.New(apiGateway, err.Error(), http.StatusBadRequest)
		ah.setEntryField(entry, c.Request, body, http.StatusBadRequest, err).Info()
		ps.SetTag("status", http.StatusBadRequest).LogFields(log.Error(err))
		return
	}

	if err := ah.validate.Struct(&body); err != nil {
		c.Status(http.StatusBadRequest)
		err := errors.New(apiGateway, err.Error(), http.StatusBadRequest)
		ah.setEntryField(entry, c.Request, body, http.StatusBadRequest, err).Info()
		ps.SetTag("status", http.StatusBadRequest).LogFields(log.Error(err))
		return
	}

	ss := c.GetHeader("Unique-Authorization")
	if _, err := jwt.ParseDuplicateCertClaimFromJWT(ss); err != nil {
		c.Status(http.StatusForbidden)
		err := errors.New(apiGateway, err.Error(), http.StatusForbidden)
		ah.setEntryField(entry, c.Request, body, http.StatusForbidden, err).Info()
		ps.SetTag("status", http.StatusForbidden).LogFields(log.Error(err))
		return
	}

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", xid)
	ctx = metadata.Set(ctx, "Unique-Authorization", ss)

	var cs opentracing.Span
	var resp *authProto.BeforeCreateAuthResponse
	err := ah.breaker[userCreateIndex].Run(func() (err error) {
		cs = tr.StartSpan(beforeCreateAuth, opentracing.ChildOf(ps.Context())).SetTag("X-Request-Id", xid)
		opts := []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
		req := &authProto.BeforeCreateAuthRequest{
			UserId:       body.UserId,
			UserPw:       body.UserPw,
			Name:         body.Name,
			PhoneNumber:  body.PhoneNumber,
			Email:        body.Email,
			Introduction: body.Introduction,
		}
		resp, err = ah.cli.BeforeCreateAuth(ctx, req, opts...)
		cs.LogFields(log.Object("request", req), log.Object("response", resp))
		return
	})

	if err == breaker.ErrBreakerOpen {
		var code = http.StatusServiceUnavailable
		c.Status(code)
		err := errors.New(userClient, breaker.ErrBreakerOpen.Error(), int32(code))
		ah.setEntryField(entry, c.Request, body, code, err).Error()
		ps.SetTag("status", code).LogFields(log.Error(err))
		return
	}

	if err != nil {
		var code = http.StatusInternalServerError
		if err, ok := err.(*errors.Error); ok { code = int(err.Code) }
		c.Status(code)
		ah.setEntryField(entry, c.Request, body, code, err).Error()
		ps.SetTag("status", http.StatusInternalServerError)
		cs.Finish()
		return
	}

	c.JSON(int(resp.Status), resp)
	entry = ah.setEntryField(entry, c.Request, body, int(resp.Status), err)
	if resp.Status == http.StatusInternalServerError {
		entry.Warn()
	} else {
		entry.Info()
	}
	ps.SetTag("status", resp.Status)
	cs.Finish()
	return
}

func (ah AuthHandler) setEntryField(entry *logrus.Entry, r *http.Request, body interface{}, outcome int, err error) *logrus.Entry {
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
	}).WithFields(ah.toLogrus.Request(r))
}