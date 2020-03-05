// 타임 아웃은 마이크로서비스에서 다른 서비스나 데이터 저장소와 통신할 때 매우 유용한 패턴이다.
// 이를 이용해 다운스트림 서비스의 에러를 탐지하여 다시 시도하거나 업스트림 서비스에 실패 메세지를 보내는 등과 같은 실패 처리 패턴을 작성할 수 있다.
// 이번 예제에서는 go-kit 툴킷에서 추천하는 eapache의 deadline 패키지를 사용할 것 이다.

package main

import (
	"flag"
	"fmt"
	"github.com/eapache/go-resiliency/deadline"
	"time"
)

func main() {
	var timeout = flag.String("timeout", "off", "Flag to set whether to apply a timeout.")
	var err error = nil
	flag.Parse()

	if *timeout == "on" {
		dl := deadline.New(1 * time.Second)
		err = dl.Run(func(stopper <-chan struct{}) error {
			slowFunction()
			return nil
		})
	} else {
		slowFunction()
	}

	switch err {
	case deadline.ErrTimedOut:
		fmt.Println("Time Out")
	default:
		fmt.Println(err)
	}
}

func slowFunction() {
	for i:=0; i<10; i++ {
		fmt.Printf("Loop: %d\n", i)
		time.Sleep(1 * time.Second)
	}
}
