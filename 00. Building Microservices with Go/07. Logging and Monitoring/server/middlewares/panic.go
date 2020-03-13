package middlewares

import (
	"fmt"
	"github.com/alexcesaro/statsd"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

// panic이 났을 경우를 대비하여 defer recover 구문으로 서버가 다운되는 것을 막고, 그에 따는 처리를 하는 미들웨이이다.
type panicMiddleware struct {
	statsD *statsd.Client
	logger *logrus.Logger
	next   http.Handler
}

func (pm *panicMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil { pm.handlingPanic(rw, r, err.(error))}
	}()

	pm.next.ServeHTTP(rw, r)
}

func (pm *panicMiddleware) handlingPanic(rw http.ResponseWriter, r *http.Request, err error) {
	pm.logger.WithFields(logrus.Fields{
		"group": "middleware",
		"segment": "panic",
		"method": r.Method,
		"path": r.URL.Path,
		"query": r.URL.RawQuery,
	}).Error(fmt.Sprintf("Error: %v\n%s", err, debug.Stack()))

	rw.WriteHeader(http.StatusInternalServerError)
}

func NewPanicMiddleware(statsD *statsd.Client, logger *logrus.Logger, next http.Handler) *panicMiddleware {
	return &panicMiddleware{
		statsD: statsD,
		logger: logger,
		next:   next,
	}
}