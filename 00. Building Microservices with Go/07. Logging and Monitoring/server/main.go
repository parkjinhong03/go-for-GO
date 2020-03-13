// 이번 07장에서는 측정 지표(metric)를 이용하여 유용한 로깅 및 모니터링을 하는 방법에 대해 알아볼 것 이다.
// logrus-logstash-hook 그리고 statsd 패키지를 이용하여 로깅과 모니터링을 구현할 것 이다.
// 참고로 로깅을 할 때, 명명 규칙(naming convention)을 정의하는 것은 매우 중요하다.
// 이번 예제에서는, 아주 대중적이고 가독성 있는 점 표기법을 이용하여 서비스의 이름을 나눴다.
// - 형식: environment.host.service.group.segment.outcome
// - environment: 운영이나 스테이징과 같은 작업 환경
// - host: 어플리케이션 서비스를 실행하는 서버의 호스트 이름
// - service: 서비스의 이름
// - group: 최상위 그룹으로, API인 경우 핸들러일 수 있음
// - segment: 그룹의 하위 레벨 정보, 일반적으로 API인 경우 핸들러의 이름
// - outcome: 동적의 결과를 나타내며 API인 경우 HTTP 상태 코드를 사용할 수 있음
// ex) production.172.30.1.24.helloWorldServer.handlers.helloWorld.ok

package main

import (
	"building-microservices-with-go.com/logging/handlers"
	"building-microservices-with-go.com/logging/middlewares"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alexcesaro/statsd"
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
)

const port = 8091

func main() {
	statsD, err := createStatsDClient(os.Getenv("STATSD"))
	if err != nil {
		log.Fatal(err)
	}

	logger, err := createLogger(os.Getenv("LOGSTASH"))
	if err != nil {
		log.Fatal(err)
	}

	setupHandlers(statsD, logger)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))

}

func setupHandlers(statsD *statsd.Client, logger *logrus.Logger) {
	helloWorldHandler := middlewares.NewValidationMiddleware(
		statsD,
		logger,
		handlers.NewHelloWorldHandler(statsD, logger),
	)

	bangHandler := middlewares.NewPanicMiddleware(
		statsD,
		logger,
		handlers.NewBangHandler(),
	)

	http.Handle("/helloworld", middlewares.NewCorrelationMiddleware(helloWorldHandler))
	http.Handle("/bang", middlewares.NewCorrelationMiddleware(bangHandler))
}

// statsD 서버의 주소를 매개변수로 받아 해당 서버의 클라이언트를 생성하여 반환하는 함수이다.
func createStatsDClient(address string) (*statsd.Client, error) {
	if address == "" {
		return &statsd.Client{}, errors.New("please set environment variable 'STATSD'")
	}

	// statsd.New 함수를 이용하여 서버 주소를 전달해 클라이언트를 생성할 수 있다.
	return statsd.New(statsd.Address(address))
}

// logrus와 logrustash 패키지를 이용하여 logrustash를 포함하고 있는 로거를 생성하여 반환하는 함수이다.
func createLogger(address string) (*logrus.Logger, error) {
	if address == "" {
		return &logrus.Logger{}, errors.New("please set environment variable 'LOGSTASH'")
	}

	// logrus.New 함수를 이용하여 기본값이 지정되어있는 로거 객채를 얻을 수 있다.
	l := logrus.New()
	// os.Hostname 함수를 이용하여 서비스 이름의 host 부분에 명시할 문자열을 얻을 수 있다.
	hostname, _ := os.Hostname()

	for i:=0; i<10; i++ {
		// 매개변수로 받은 주소에 tcp 프로토콜을 이용하여 연결을 시도한다.
		conn, err := net.Dial("tcp", address)

		// 만약 연결에 성공하지 못했다면, 로그를 찍고 1초 후에 다시 시도한다.
		if err != nil {
			log.Println("Unable to connect to logstash, retrying")
			time.Sleep(1 * time.Second)
			continue
		}

		// 만약 연결에 성공했다면, logrustash.New 함수를 이용하여 logstash 객체를 생성한다.
		hook := logrustash.New(
			// 출력은 위에서 생성한 tcp 연결 객체인 conn에 하도록 등록한다.
			conn,
			// 그리고 로그 메시지에 기본적으로 {"hostname": os.Hostname()}을 추가하도록 설정한다.
			logrustash.DefaultFormatter(
				logrus.Fields{"hostname": hostname},
			),
		)

		// logrus로 선언한 로깅 객체의 훅에 logrustash로 선언한 훅을 추가한다.
		l.Hooks.Add(hook)
		return l, nil
	}

	// 만약 10번의 시도 모두 성공하지 못했다면 에러를 반환시킨다.
	return nil, errors.New("unable to connect to logstash")
}
