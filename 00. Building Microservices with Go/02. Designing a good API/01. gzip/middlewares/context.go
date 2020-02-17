package middlewares

import (
	"../entities"
	"context"
	"encoding/json"
	"net/http"
)

type contextMiddleWare struct {
	next http.Handler
}

func NewContextMiddleware(next http.Handler) *contextMiddleWare {
	return &contextMiddleWare{next: next}
}

type ValidationContextKey string

func (m *contextMiddleWare) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request entities.HelloWorldRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
	}

	ctx := context.WithValue(r.Context(), ValidationContextKey("name"), request.Name)
	r = r.WithContext(ctx)

	m.next.ServeHTTP(rw, r)
}