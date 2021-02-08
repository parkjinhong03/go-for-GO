// 상태 점검(health check)은 마이크로서비스에서 반드시 필요한 사항이다.
// 이는 오작동이나 장애가 발생했을 때, 어플리케이션을 실행해주는 프로세스가 어플리케이션을 다시 시작하거나 종료하도록 해주기 때문에 매우 중요하다.
// 이번 예제에서는 ewma 패키지를 이용하여 현재 응답 시간(이동 평균)을 기록하여 상태 점검을 구현할 것 이다.

package main

import (
	"fmt"
	"github.com/VividCortex/ewma"
	"github.com/labstack/gommon/log"
	"net/http"
	"sync"
	"time"
)

// 시그니처가 정해진 여러 함수에서 접근해야 하므로 지역변수로 선언해야 한다.
var (
	// ewma.MovingAverage 객체는 이동 평균을 구하여 값을 저장해주는 객체로, 평균 응답 시간을 기록하기 위해 선언해야한다.
	ma ewma.MovingAverage
	// sync.RWMutex 객체는 일정 코드에 2개 이상의 프로세스가 접근하지 못하게 하기 위해 선언해야한다.
	resetMux = sync.RWMutex{}
	// 이동 평균 값을 초기화 하는 함수에 오직 한 프로세스밖에 접근하지 못하도록 하기 위해서 선언하는 변수이다.
	resetting = false
)

// 아래의 모든 상수들은 단순 구현을 위한 일시적인 값일 뿐, 절대적인 값이 절대 아니다.
const (
	// 평균 응답 시간의 임계치로, 실제 평균 응답 시간이 이 임계치 값을 넘으면 서비스가 충분히 복구될때까지 대기한다.
	threshold = time.Second
	// 위에서 말한 것 처럼 평균 시간이 임계치를 넘으면 아래 선언된 시간만큼 대기시킨다.
	sleepTime = 5 * time.Second
)

func main() {
	// ewma.NewMovingAverage 메서드를 이용하여 이동 평균을 구할 수 있는 객체를 생성한다.
	ma = ewma.NewMovingAverage()

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/health", healthHandler)

	log.Fatal(http.ListenAndServe(":8008", nil))
}

func mainHandler(rw http.ResponseWriter, r *http.Request) {
	// 핸들러의 총 실행시간을 기록하기 위해 time.Now 함수를 이용하여 현재 시간을 기록해논다.
	startTime := time.Now()

	// 건강 상태를 체크한 후, 상태가 불안정하다면 respondServiceUnhealthy 함수를 실행시키고 반환한다.
	if !isHealthy() {
		respondServiceUnhealthy(rw)
		return
	}

	// 서비스가 과부화 되었다는 상황을 묘사하기 위해 임의로 1초 대기시간을 가지게 한다.
	time.Sleep(time.Second)
	rw.WriteHeader(http.StatusOK)

	// time.Time.Sub 메서드를 이용하여 위에서 저장해 놨던 시간과 현재 시간의 차를 구하여 응답 시간을 구한다.
	duration := time.Now().Sub(startTime)
	// ewma.MovingAverage.Add 메서드에 응답 시간을 매개변수로 넘겨 이동 평균을 다시 구한다.
	ma.Add(float64(duration))
	// ewma.MovingAverage.Value 메서드를 이용하여 얻은 이동 평균을 1000000으로 나누어서 ms 단위로 바꿀 수 있다.
	_, _ = fmt.Fprintf(rw, "Average request time: %f (ms)\n", ma.Value()/1000000)
}

func healthHandler(rw http.ResponseWriter, r *http.Request) {
	// 건강 상태를 체크한 후, 상태가 불안정 하다면 503 ServiceUnavailable 상태코드와 함께 반환한다.
	if !isHealthy() {
		rw.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintln(rw, "OK")
}

// 상태가 불안정할 때 호출하는 함수로, 응답을 처리하고 상태를 복구시키는(지정 시간동안 대기) 기능을 하는 함수이다.
func respondServiceUnhealthy(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusServiceUnavailable)

	// 동시에 여러 요청이 들어왔을 때, 그 요청을 순차적으로 처리하기 위해 sync.RWMutex.Lock 함수를 사용한다.
	resetMux.Lock()
	defer resetMux.Unlock()

	if !resetting {
		go sleepAndResetAverage()
	}
}

// 서비스가 정상적으로 복구될 수 있을 만큼 대기한 후 이동 평균 값을 초기화하는 함수이다.
func sleepAndResetAverage() {
	// 초기화 중이라는 표시를 하기 위해 resetting 값을 true로 변경한다.
	resetting = true

	time.Sleep(sleepTime)
	ma = ewma.NewMovingAverage()

	// 초기화가 끝났다는 표시를 하기 위해 resetting 값을 false로 변경한다.
	resetting = false
}

// 임계치와 이동 평균값을 비교하여 현재 서비스의 상태가 건강한지 확인하는 함수이다.
func isHealthy() bool {
	return ma.Value() < float64(threshold)
}

