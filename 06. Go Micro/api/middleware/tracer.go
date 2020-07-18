package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func (m Middleware) AuthTracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		tr, cl, err := m.ajc.NewTracer(
			jaegercfg.Logger(jaegerlog.StdLogger),
		)
		if err != nil {
			c.Set("error", err)
			c.Next()
			return
		}
		defer func() { _ = cl.Close() } ()

		opentracing.SetGlobalTracer(tr)
		c.Set("tracer", tr)
		c.Next()
	}
}

func (m Middleware) UserTracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		tr, cl, err := m.ujc.NewTracer(
			jaegercfg.Logger(jaegerlog.StdLogger),
		)
		if err != nil {
			c.Set("middleware_err", err)
			c.Next()
			return
		}
		defer func() { _ = cl.Close() } ()

		opentracing.SetGlobalTracer(tr)
		c.Set("tracer", tr)
		c.Next()
	}
}