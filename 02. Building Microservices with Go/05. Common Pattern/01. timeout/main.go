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
		// deadline.New 함수를 이용하여 원하는 함수가 넘긴 매개변수 시간 안에 끝나지 않으면 타임 아웃시키는 객체를 생성할 수 있다.
		dl := deadline.New(1 * time.Second)
		// 위에서 만든 객체의 Run 메서드를 이용하여 타임아웃을 실행시킬 수 있다.
		// 매개변수로 넘긴 함수가 해당 시간 안에 끝나지 않으면 즉시 deadline.ErrTimedOut을 반환하고, 시간 안에 끝나면 해당 함수의 반환값을 반환한다.
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
