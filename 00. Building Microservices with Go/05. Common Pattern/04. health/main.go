// 상태 점검(health check)은 마이크로서비스에서 반드시 필요한 사항이다.
// 이는 오작동이나 장애가 발생했을 때, 어플리케이션을 실행해주는 프로세스가 어플리케이션을 다시 시작하거나 종료하도록 해주기 때문에 매우 중요하다.
// 이번 예제에서는 ewma 패키지를 이용하여 현재 응답 시간(이동 평균)을 기록하여 상태 점검을 구현할 것 이다.

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

// 시그니처가 정해진 여러 함수에서 접근해야 하므로 지역변수로 선언해야 한다.
var (
	// ewma.MovingAverage 객체는 이동 평균을 구하여 값을 저장해주는 객체로, 평균 응답 시간을 기록하기 위해 선언해야한다.
	ma ewma.MovingAverage
	// sync.RWMutex 객체는 일정 함수에 2개 이상의 프로세스가 접근하지 못하게 하기 위해 선언해야한다.
	resetMux = sync.RWMutex{}
)

// 아래의 모든 상수들은 단순 구현을 위한 일시적인 값일 뿐, 절대적인 값이 절대 아니다.
const (
	// 평균 응답 시간의 임계치로, 실제 평균 응답 시간이 이 임계치 값을 넘으면 서비스가 충분히 복구될때까지 대기한다.
	threshold = time.Second
	// 위에서 말한 것 처럼 평균 시간이 임계치를 넘으면 아래 선언된 시간만큼 대기시킨다.
	sleepTime = 10 * time.Second
	// 이 상수는 sync.RWMutex.Lock 메서드로 인해 일정 프로세스가 대기 상태가 되었을 때, 대기할 수 있는 최대 시간을 정해논 값이다.
	timeOut = time.Millisecond
)

func main() {
	// ewma.NewMovingAverage 메서드를 이용하여 이동 평균을 구할 수 있는 객체를 생성한다.
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

