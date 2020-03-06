// 이번에는 전에 나왔던 타임 아웃과 백 오프들을 보완하는 또 다른 패턴인 회로 차단에 대해 설명할 것 이다.
// 회로 차단기는 시스템이 스트레스를 받고 있을 때(계속하여 에러가 발생할 떄) 자동으로 기능을 축소시키는 패턴이다.
// 회로 차단 패턴을 적용하지 않으면 업트스림 서비스에 추가적인 로드를 발생시키고, 이로 인해 전체 서비스가 중단될 수도 있다.
// 따라서 마이크로 서비스에서 매우 중요한 회로 차단을 이번에도 eapache의 breaker 패키지를 이용하여 구현해볼 것 이다.

package main

import (
	"fmt"
	"github.com/eapache/go-resiliency/breaker"
	"time"
)

var a = 1

func main() {
	// breaker.New 함수에 다음과 같은 세 가지 매개 변수를 사용해 회로 차단기를 구성한다.
	// 1. errorThreshold: 해당 횟수만큼 에러가 발생하면 회로를 차단한다.
	// 2. successThreshold: 반 개방 상태에서 해당 횟수만큼 요청 처리에 성공하면 닫힌 상태로 변경한다.
	// 3. timeout: 회로가 차단된 후 반 개방 상태로 변경되기 전까지 대기해야 하는 시간.
	b := breaker.New(3, 1, 5 * time.Second)

	for {
		// b.Run 메서드를 이용하여 매개변수로 넘긴 함수롸 회로 차단기를 연결시킨 후 한번 실행한다.
		// 만약 이 함수가 등록한 임계치만큼 실패한다면(nil 반환 X), 차단기는 켜지게 되고 메서드가 호출되자마자 breaker.ErrBreakerOpen을 반환한다.
		// 차단기가 켜지게 된 후, 등록한 시간이 지나면 회로가 반 개방 상태로 바뀌게 되고, 그 때 함수가 성공한다면(nil 반환) 차단기는 다시 닫히게 된다.
		err := b.Run(func() error {
			time.Sleep(2 * time.Second)
			if a >= 10 {
				return nil
			}
			return fmt.Errorf("time out")

		})

		switch err {
		case nil:
			fmt.Println("Success")
		case breaker.ErrBreakerOpen:
			fmt.Println("Breaker open")
		default:
			fmt.Println(err)
		}
		a++

		time.Sleep(500 * time.Millisecond)
	}
}
