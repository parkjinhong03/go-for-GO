package middleware

import (
	"github.com/gin-gonic/gin"
)

type Status int

func (s Status) ifIncludeIn(arr []int) bool {
	for _, i := range arr {
		if i == int(s) { return true }
	}
	return false
}

type customWriter struct {
	gin.ResponseWriter
	s Status
}

func (cw *customWriter) WriteHeader(status int) {
	cw.s = Status(status)
	cw.ResponseWriter.WriteHeader(status)
}

func CustomWriter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer = &customWriter{ResponseWriter: c.Writer}
		c.Next()
	}
}
