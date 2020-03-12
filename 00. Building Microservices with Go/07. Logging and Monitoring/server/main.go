// 이번 07장에서는 측정 지표(metric)를 이용하여 유용한 로깅 및 모니터링을 하는 방법에 대해 알아볼 것 이다.
// logrus-logstash-hook 그리고 statsd 패키지를 이용하여 로깅과 모니터링을 구현할 것 이다.

package main

import (
	"errors"
	"log"
	"os"

	"github.com/alexcesaro/statsd"
)

const port = 8080

func main() {
	statsd, err := createStatsDClient(os.Getenv("STATSD"))
	if err != nil {
		log.Fatal(err)
	}
}

// statsD 서버의 주소를 매개변수로 받아 해당 서버의 클라이언트를 생성하여 반환하는 함수이다.
func createStatsDClient(address string) (*statsd.Client, error) {
	if address == "" {
		return &statsd.Client{}, errors.New("please set environment variable 'address'")
	}

	// statsd.New 함수를 이용하여 서버 주소를 전달해 클라이언트를 생성할 수 있다.
	return statsd.New(statsd.Address(address))
}

