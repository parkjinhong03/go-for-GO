// 백 오프도 타임 아웃처럼 마이크로서비스에서 다운스트림 서비스의 장애로부터 시스템을 보호하는 데 도움이 된다.
// 백 오프란 에러가 발생하였을 때, 즉시 다시 시도하지 않고 점점 늘어나는 일정 시간동안 대기한 후 다시 재시도하는 전략이다.
// 이번 예제에서도 eapache의 retrier 패키지를 이용하여 구현해볼 것 이다.

package main

import (
	"fmt"
	"github.com/eapache/go-resiliency/retrier"
	"time"
)

func main() {
	var n = 1

	// retirer.New 함수를 이용하여 백 오프를 실행시킬 수 있는 객체를 생성할 수 있다.
	// 첫 번째 매개 변수는 time.Duration의 슬라이스로, 인덱스의 갯수만큼 재시도하고 각각의 값만큼 시간을 대기한다.
	// 두 번째 매개 변수는 분류자(Classifier)로, 재시도가 허용되는 에러 타입 또는 즉시 실패시키는 에러 타입을 제어할 수 있다.
	r := retrier.New(retrier.ConstantBackoff(3, 1 * time.Second), nil)
	// 위에서 만든 객체의 Run 메서드를 이용하여 백 오프를 실행시킬 수 있다.
	// 매개변수로 넘긴 함수가 한번이라도 성공(nil 반환)하면 nil을 반환하고, 위에서 정한 횟수만큼 재시도해도 실패한다면 해당 함수가 반환한 에러를 반환한다.
	err := r.Run(func() error {
		fmt.Printf("Attept: %d\n", n)
		n++
		return fmt.Errorf("failed")
	})

	if err != nil {
		fmt.Println(err)
	}
}