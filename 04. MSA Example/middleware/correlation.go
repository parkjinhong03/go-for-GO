package middleware

import (
	"github.com/google/uuid"
	"net/http"
)

type correlationMiddleware struct {
	next http.Handler
}

func (cm *correlationMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Request-ID") == "" {
		r.Header.Set("X-Request-ID", uuid.New().String())
	}

	cm.next.ServeHTTP(rw, r)
}

func NewCorrelationMiddleware(next http.Handler) http.Handler {
	return &correlationMiddleware{
		next: next,
	}
}