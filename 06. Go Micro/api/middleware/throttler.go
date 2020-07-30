package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type throttler struct {
	limiter  map[string]*uint32
	rejected map[string]bool
	locker   sync.Mutex
}

func Throttler() *throttler {
	return &throttler{
		limiter:  make(map[string]*uint32),
		rejected: make(map[string]bool),
	}
}

func (t *throttler) Throttling(c *gin.Context) {
	cip := c.ClientIP()

	t.locker.Lock()
	defer t.locker.Unlock()

	if _, ok := t.limiter[cip]; !ok {
		var init uint32 = 0
		t.limiter[cip] = &init
	}
	if _, ok := t.rejected[cip]; !ok {
		t.rejected[cip] = false
	}

	if t.rejected[cip] {
		c.Status(http.StatusForbidden)
		return
	}

	if *t.limiter[cip] >= 10 {
		c.Status(http.StatusTooManyRequests)
		t.rejected[cip] = true
		return
	}

	atomic.AddUint32(t.limiter[cip], 1)
	time.AfterFunc(time.Second, func() {
		t.locker.Lock()
		*t.limiter[cip] -= 1
		t.locker.Unlock()
	})
}