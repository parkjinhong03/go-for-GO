package middleware

import (
	"github.com/gin-gonic/gin"
)

type customWriter struct {
	gin.ResponseWriter
	status int
}

func (cw *customWriter) WriteHeader(status int) {
	cw.status = status
	cw.ResponseWriter.WriteHeader(status)
}

func (m Middleware) CustomWriter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer = &customWriter{ResponseWriter: c.Writer}
		c.Next()
	}
}