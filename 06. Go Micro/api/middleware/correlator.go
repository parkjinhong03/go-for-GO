package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Correlator() gin.HandlerFunc {
	return func(c *gin.Context) {
		xReqId := uuid.New().String()
		c.Request.Header.Set("X-Request-Id", xReqId)
		c.Header("X-Request-Id", xReqId)
		c.Next()
	}
}
