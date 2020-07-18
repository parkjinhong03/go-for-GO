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
	"github.com/google/uuid"
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

	//xid := c.GetHeader("X-Request-Id")
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

	cs := tr.StartSpan(userIdDuplicate, opentracing.ChildOf(ps.Context())).SetTag("X-Request-Id", c.GetHeader("X-Request-Id"))
	opts := []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
	req := &authProto.UserIdDuplicatedRequest{
		UserId: body.UserId,
	}
	resp, err := ah.cli.UserIdDuplicated(ctx, req, opts...)
	cs.LogFields(log.Object("request", req), log.Object("response", resp))

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

	////cs := tr.StartSpan("callRpcHandler", opentracing.ChildOf(ps.Context()))
	//switch err = ah.breaker[userIdDuplicateIndex].Run(reqFunc); err {
	//case nil:
	//	c.JSON(int(resp.Status), resp)
	//
	//	ah.notified[userIdDuplicateIndex] = false
	//	ah.setEntryField(entry, c.Request, body, int(resp.Status), err).Info()
	//
	//case breaker.ErrBreakerOpen:
	//	c.Status(http.StatusServiceUnavailable)
	//
	//	if ah.notified[userIdDuplicateIndex] == false {
	//		// 처음으로 열린 차단기라면, 알림 서비스 실행
	//		ah.notified[userIdDuplicateIndex] = true
	//	}
	//
	//	ah.setEntryField(entry, c.Request, body, http.StatusServiceUnavailable, err).Error()
	//default:
	//	var code = http.StatusInternalServerError
	//	err, ok := err.(*errors.Error)
	//	if ok {
	//		code = int(err.Code)
	//	}
	//	c.Status(code)
	//
	//	ah.setEntryField(entry, c.Request, body, code, err).Warn()
	//}

	return
}

func (ah AuthHandler) UserCreateHandler(c *gin.Context) {
	entry := ah.logger.WithFields(logrus.Fields{
		"group":   "handler",
		"segment": "userCreate",
	})

	var body entity.UserCreate
	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		ah.setEntryField(entry, c.Request, body, http.StatusBadRequest, err).Info()
		return
	}

	if err := ah.validate.Struct(&body); err != nil {
		c.Status(http.StatusBadRequest)
		ah.setEntryField(entry, c.Request, body, http.StatusBadRequest, err).Info()
		return
	}

	xReqId := c.GetHeader("X-Request-Id")
	if _, err := uuid.Parse(xReqId); err != nil {
		c.Status(http.StatusForbidden)
		ah.setEntryField(entry, c.Request, body, http.StatusForbidden, err).Info()
		return
	}

	ss := c.GetHeader("Unique-Authorization")
	if _, err := jwt.ParseDuplicateCertClaimFromJWT(ss); err != nil || ss == "" {
		c.Status(http.StatusForbidden)
		ah.setEntryField(entry, c.Request, body, http.StatusForbidden, err).Info()
		return
	}

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", xReqId)
	ctx = metadata.Set(ctx, "Unique-Authorization", ss)

	resp := new(authProto.BeforeCreateAuthResponse)
	reqFunc := func() (err error) {
		opts := []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
		if resp, err = ah.cli.BeforeCreateAuth(ctx, &authProto.BeforeCreateAuthRequest{
			UserId:       body.UserId,
			UserPw:       body.UserPw,
			Name:         body.Name,
			PhoneNumber:  body.PhoneNumber,
			Email:        body.Email,
			Introduction: body.Introduction,
		}, opts...); err != nil {
			return
		}
		if resp.Status == http.StatusInternalServerError {
			err = errors.InternalServerError("go.micro.client", "internal server error")
		}
		return
	}

	var err error
	switch err = ah.breaker[userCreateIndex].Run(reqFunc); err {
	case nil:
		c.JSON(int(resp.Status), resp)

		ah.notified[userCreateIndex] = false
		ah.setEntryField(entry, c.Request, body, int(resp.Status), err).Info()
	case breaker.ErrBreakerOpen:
		c.Status(http.StatusServiceUnavailable)

		ah.setEntryField(entry, c.Request, body, http.StatusServiceUnavailable, err).Error()
		if ah.notified[userCreateIndex] == false {
			// 처음으로 열린 차단기라면, 알림 서비스 실행
			ah.notified[userCreateIndex] = true
		}
		// 블록된 요청에 대한 추가 작업 실행
	default:
		var code = http.StatusInternalServerError
		err, ok := err.(*errors.Error)
		if ok {
			code = int(err.Code)
		}
		c.Status(code)

		ah.setEntryField(entry, c.Request, body, code, err).Warn()
	}

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