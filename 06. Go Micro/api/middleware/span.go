package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (m Middleware) AuthServiceSpan() gin.HandlerFunc {
	return func(c *gin.Context) {
		//tr, cl, err := ah.jgConf.NewTracer(
		//	jaegercfg.Logger(jaegerlog.StdLogger),
		//)
		//if err != nil {
		//	c.Status(http.StatusInternalServerError)
		//	ah.setEntryField(entry, c.Request, body, http.StatusInternalServerError, err).Warn()
		//	return
		//}
		//defer func() { _ = cl.Close() } ()
		//
		//opentracing.SetGlobalTracer(tr)
		//ps := tr.StartSpan(c.Request.URL.Path)
		//ps.SetTag("X-Request-Id", c.GetHeader("X-Request-Id")).SetTag("segment", "userIdDuplicate")
		fmt.Println(1)
		c.Next()
	}
}

func (m Middleware) UserServiceSpan() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(2)
		c.Next()
	}
}