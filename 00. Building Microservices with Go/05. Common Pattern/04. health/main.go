package main

import (
	"fmt"
	"github.com/VividCortex/ewma"
	"github.com/eapache/go-resiliency/deadline"
	"github.com/labstack/gommon/log"
	"net/http"
	"sync"
	"time"
)

var (
	ma ewma.MovingAverage
	resetMux = sync.RWMutex{}
)
const (
	threshold = time.Second
	sleepTime = 10 * time.Second
	timeOut = time.Millisecond
)

func main() {
	ma = ewma.NewMovingAverage()

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/health", healthHandler)

	log.Fatal(http.ListenAndServe(":8008", nil))
}

func mainHandler(rw http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	if !isHealthy() {
		respondServiceUnhealthy(rw)
		return
	}

	time.Sleep(time.Second)
	rw.WriteHeader(http.StatusOK)

	duration := time.Now().Sub(startTime)
	ma.Add(float64(duration))
	_, _ = fmt.Fprintf(rw, "Average request time: %f (ms)\n", ma.Value()/1000000)
}

func healthHandler(rw http.ResponseWriter, r *http.Request) {
	if !isHealthy() {
		rw.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(rw, "OK")
}

func respondServiceUnhealthy(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusServiceUnavailable)


	dl := deadline.New(timeOut)
	err := dl.Run(func(i <-chan struct{}) error {
		resetMux.Lock()
		return nil
	})

	switch err {
	case nil:
		go sleepAndResetAverage()
		return
	default:
		return
	}

}

func sleepAndResetAverage() {
	time.Sleep(sleepTime)
	ma = ewma.NewMovingAverage()
	resetMux.Unlock()
}

func isHealthy() bool {
	return ma.Value() < float64(threshold)
}

