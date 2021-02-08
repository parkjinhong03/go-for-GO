package middlewares

import (
	"github.com/google/uuid"
	"net/http"
)

// 사용자의 하나의 요청에 대한 모든 서비스 호출을 동일한 ID로 표시하여 요청을 디버깅할 수 있게 하기 위한 미들웨어이다.
type correlationMiddleware struct {
	next http.Handler
}

func (cm *correlationMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Request-ID") == "" {
		r.Header.Set("X-Request-ID", uuid.New().String())
	}

	cm.next.ServeHTTP(rw, r)
}

func NewCorrelationMiddleware(next http.Handler) *correlationMiddleware {
	return &correlationMiddleware{next: next}
}