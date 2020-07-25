package handler

import (
	"context"
	"gateway/adapter/consul"
	"gateway/entity"
	authProto "gateway/proto/golang/auth"
	"gateway/tool/conf"
	"gateway/tool/logrusfield"
	topic "gateway/topic/golang"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"net/http"
	"reflect"
)

type AuthHandler struct {
	cli      authProto.AuthService
	logger   *logrus.Logger
	validate *validator.Validate
	consul   *api.Client
	nodes    []*registry.Node
	tracer   opentracing.Tracer
	//breaker []*breaker.Breaker
	breakers map[string]*breaker.Breaker
	brConf   conf.BreakerConfig
	next     selector.Next
}

var (
	defaultOpts = []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
)

func NewAuthHandler(cli authProto.AuthService, logger *logrus.Logger, validate *validator.Validate,
	consul *api.Client, tracer opentracing.Tracer, bcConf conf.BreakerConfig) *AuthHandler {

	return &AuthHandler{
		cli:      cli,
		logger:   logger,
		validate: validate,
		consul:   consul,
		tracer:   tracer,
		brConf:   bcConf,
		//breaker: []*breaker.Breaker{bk1, bk2},
		breakers: make(map[string]*breaker.Breaker),
		next:     selector.RoundRobin([]*registry.Service{}),
	}
}

func (ah *AuthHandler) UserIdDuplicateHandler(c *gin.Context) {
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
		err = errors.New(topic.ApiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	if err := ah.validate.Struct(&body); err != nil {
		code = http.StatusBadRequest
		c.Status(code)
		err := errors.New(topic.ApiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", c.GetHeader("X-Request-Id"))
	ctx = metadata.Set(ctx, "Unique-Authorization", c.GetHeader("Unique-Authorization"))

	nds, err := consul.GetServiceNodes(ah.consul)
	if err != nil || nds == nil {
		code = http.StatusServiceUnavailable
		c.Status(code)
		err := errors.New(topic.ApiGateway, "There are no services registered in consul.", int32(code))
		entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Error()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	if !reflect.DeepEqual(ah.nodes, nds) {
		ah.nodes = nds
		ah.next = selector.RoundRobin([]*registry.Service{{ Nodes: ah.nodes }})
	}

	nd, err := ah.next()
	if err != nil {
		code = http.StatusServiceUnavailable
		c.Status(code)
		err := errors.New(topic.ApiGateway, err.Error(), int32(code))
		entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Error()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	if _, ok := ah.breakers[nd.Id]; !ok {
		ah.breakers[nd.Id] = breaker.New(ah.brConf.ErrorThreshold, ah.brConf.SuccessThreshold, ah.brConf.Timeout)
	}

	var resp *authProto.UserIdDuplicatedResponse
	err = ah.breakers[nd.Id].Run(func() (err error) {
		req := body.ToRequestProto()
		opts := append(defaultOpts, client.WithAddress(nd.Address))
		cs := ah.tracer.StartSpan(userIdDuplicate, opentracing.ChildOf(ps.Context())).SetTag("X-Request-Id", xid)
		ctx = metadata.Set(ctx, "Span-Context", cs.Context().(jaeger.SpanContext).String())
		resp, err = ah.cli.UserIdDuplicated(ctx, req, opts...)
		md, _ := metadata.FromContext(ctx)
		cs.LogFields(log.Object("req", req), log.Object("resp", resp), log.Object("ctx", md), log.Error(err))
		cs.Finish()
		return
	})

	if err == breaker.ErrBreakerOpen {
		code = http.StatusServiceUnavailable
		c.Status(code)
		err := errors.New(authClient, breaker.ErrBreakerOpen.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Error()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	if err != nil {
		code = http.StatusInternalServerError
		if err, ok := err.(*errors.Error); ok { code = int(err.Code) }
		c.Status(code)
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Error()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	code = int(resp.Status)
	c.JSON(code, resp)
	err = errors.New(authClient, resp.Message, int32(code))
	entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
	if code == http.StatusInternalServerError {
		entry.Warn()
	} else {
		entry.Info()
	}
	ps.SetTag("status", code).LogFields(log.String("message", err.Error()))

	return
}

func (ah *AuthHandler) UserCreateHandler(c *gin.Context) {
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
		err := errors.New(topic.ApiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	if err := ah.validate.Struct(&body); err != nil {
		code = http.StatusBadRequest
		c.Status(code)
		err := errors.New(topic.ApiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	// 여기서 인증 관련 기능 추가 예정 (proto 변경 필요)

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", xid)
	ctx = metadata.Set(ctx, "Unique-Authorization", c.GetHeader("Unique-Authorization"))

	var resp *authProto.BeforeCreateAuthResponse
	err := ah.breakers["userCreateIndex"].Run(func() (err error) {
		opts := []client.CallOption{client.WithDialTimeout(DefaultDialTimeout), client.WithRequestTimeout(DefaultRequestTimeout)}
		req := body.ToRequestProto()
		cs := ah.tracer.StartSpan(beforeCreateAuth, opentracing.ChildOf(ps.Context())).SetTag("X-Request-Id", xid)
		ctx = metadata.Set(ctx, "Span-Context", cs.Context().(jaeger.SpanContext).String())
		resp, err = ah.cli.BeforeCreateAuth(ctx, req, opts...)
		md, _ := metadata.FromContext(ctx)
		cs.LogFields(log.Object("req", req), log.Object("resp", resp), log.Object("ctx", md), log.Error(err))
		cs.Finish()
		return
	})

	if err == breaker.ErrBreakerOpen {
		code = http.StatusServiceUnavailable
		c.Status(code)
		err := errors.New(authClient, breaker.ErrBreakerOpen.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Error()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	if err != nil {
		code = http.StatusInternalServerError
		if err, ok := err.(*errors.Error); ok { code = int(err.Code) }
		c.Status(code)
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Error()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	code = int(resp.Status)
	c.JSON(code, resp)
	err = errors.New(authClient, resp.Message, int32(code))
	entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
	if resp.Status == http.StatusInternalServerError {
		entry.Warn()
	} else {
		entry.Info()
	}
	ps.SetTag("status", resp.Status).LogFields(log.String("message", err.Error()))
	return
}
