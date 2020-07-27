package handler

import (
	"context"
	"gateway/adapter/consul"
	"gateway/entity"
	userProto "gateway/proto/golang/user"
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
	"time"
)

type UserHandler struct {
	cli      userProto.UserService
	logger   *logrus.Logger
	validate *validator.Validate
	consul   *api.Client
	tracer   opentracing.Tracer
	//breaker  []*breaker.Breaker
	breakers map[string]*breaker.Breaker
	brConf   conf.BreakerConfig
	nodes    []*registry.Node
	next	 selector.Next
}

func NewUserHandler(cli userProto.UserService, logger *logrus.Logger, validate *validator.Validate,
	consul *api.Client, tracer opentracing.Tracer, bcConf conf.BreakerConfig) UserHandler {

	return UserHandler{
		cli:      cli,
		logger:   logger,
		validate: validate,
		consul:   consul,
		tracer:   tracer,
		//breaker:  []*breaker.Breaker{bk},
		breakers: make(map[string]*breaker.Breaker),
		brConf:   bcConf,
		next:     selector.RoundRobin([]*registry.Service{}),
	}
}

func (uh UserHandler) EmailDuplicateHandler(c *gin.Context) {
	var body entity.EmailDuplicate
	var code int
	xid := c.GetHeader("X-Request-Id")

	ps := uh.tracer.StartSpan(c.Request.URL.Path)
	ps.SetTag("X-Request-Id", xid).SetTag("segment", "emailDuplicate")
	defer ps.Finish()

	entry := uh.logger.WithField("group", "handler").WithField("segment", "emailDuplicate")
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

	if err := uh.validate.Struct(&body); err != nil {
		code = http.StatusBadRequest
		c.Status(code)
		err := errors.New(topic.ApiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Info()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	ctx := context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", xid)
	ctx = metadata.Set(ctx, "Unique-Authorization", c.GetHeader("Unique-Authorization"))

	nds, err := consul.GetServiceNodes(topic.UserService, uh.consul)
	if err != nil || nds == nil {
		code = http.StatusServiceUnavailable
		c.Status(code)
		err := errors.New(topic.ApiGateway, "There are no services registered in consul. or consul error", int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Error()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	if !reflect.DeepEqual(uh.nodes, nds) {
		uh.nodes = nds
		uh.next = selector.RoundRobin([]*registry.Service{{ Nodes: uh.nodes }})
	}

	nd, err := uh.next()
	if err != nil {
		code = http.StatusServiceUnavailable
		c.Status(code)
		err := errors.New(topic.ApiGateway, err.Error(), int32(code))
		entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
		entry.Error()
		ps.SetTag("status", code).LogFields(log.String("message", err.Error()))
		return
	}

	if _, ok := uh.breakers[nd.Id]; !ok {
		uh.breakers[nd.Id] = breaker.New(uh.brConf.ErrorThreshold, uh.brConf.SuccessThreshold, uh.brConf.Timeout)
	}

	var resp *userProto.EmailDuplicatedResponse
	err = uh.breakers[nd.Id].Run(func() (err error) {
		req := body.ToRequestProto()
		opts := append(defaultOpts, client.WithAddress(nd.Address))
		cs := uh.tracer.StartSpan(emailDuplicate, opentracing.ChildOf(ps.Context()))
		cs.SetTag("X-Request-Id", xid).SetTag("Service-Id", nd.Id)
		ctx = metadata.Set(ctx, "Span-Context", cs.Context().(jaeger.SpanContext).String())
		resp, err = uh.cli.EmailDuplicated(ctx, req, opts...)
		md, _ := metadata.FromContext(ctx)
		cs.LogFields(log.Object("req", req), log.Object("resp", resp), log.Object("ctx", md), log.Error(err))
		cs.Finish()
		return
	})

	if err == breaker.ErrBreakerOpen {
		code = http.StatusServiceUnavailable
		c.Status(code)
		err := errors.New(userClient, breaker.ErrBreakerOpen.Error(), int32(code))
		tlsErr := uh.consul.Agent().FailTTL(nd.Metadata["CheckID"], breaker.ErrBreakerOpen.Error())
		if tlsErr != nil { err = tlsErr }
		passFunc := func(){ _ = uh.consul.Agent().PassTTL(nd.Metadata["CheckID"], "circuit breaker is close") }
		time.AfterFunc(uh.brConf.Timeout, passFunc)
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
	err = errors.New(userClient, resp.Message, int32(code))
	entry = entry.WithFields(logrusfield.ForReturn(body, code, err))
	if code == http.StatusInternalServerError {
		entry.Warn()
	} else {
		entry.Info()
	}
	ps.SetTag("status", code).LogFields(log.String("message", err.Error()))

	return
}
