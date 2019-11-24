# **병행 처리의 응용**
## **이번 화에 다룰 내용**
- **고루틴**과 **채널**을 활용하여 **병행 처리 코드**를 다양하게 작성할 수 있다.

- 이번 절에서는 **병행 처리**의 **네 가지 예제**를 소개한다
    - **타임아웃** -> 시간이 오래 걸리는 작업에 타임아웃 처리하기
    - **공유 메모리** -> 채널을 사용하여 여러 고루틴의 공유 메모리 접근 제어하기
    - **파이프 라인** -> 여러 고루틴을 파이프라인 형태로 연결하기
    - **맵리듀스** -> 고루틴을 사용하여 맵리듀스 패턴 구현하기.

- 물론 이번 절에서 소개하는 것 외에도 **다양한 패턴**을 만들 수 있다.

<br>

---
## **타임아웃**
- 시간이 **오래** 걸리는 작업에 **select** 문을 이용하면 **타임아웃** 기능을 쉽게 구현할 수 있다.
- 다음은 **process()** 작업에 **타임아웃**(10 밀리초)을 적용한 코드이다.
    ~~~go
    done := process()
    timeout := time.After(10 * time.Millisecond)

    select {
    case d := <-done:
        fmt.Println(d)
    case <-timeout:
        fmt.Println("Time out!")
    }
    ~~~
- 다음 코드는 **select** 문으로 **done**과 **timeout** 채널중 먼저 수신되는 **케이스**를 수행한다.
- 따라서 10 밀리초가 지나면 **timeout** 채널이 수신되어 **process()** 함수에 **타임아웃** 기능이 적용된다.

<br>

- 이 코드에서 **process()** 함수는 **타임아웃** 처리 이후에도 **계속 동작**한다.

- 타임아웃 처리 이후 **process()** 함수를 **제어**하기 위해 다음 **세 가지 방법**을 생각해 볼 수 있다.
    1. **아무 처리도 하지 않음**
    2. **done 채널을 닫음**
    3. **process() 함수에 타임아웃 메세지를 전송**

<br>

**다음은 process() 함수에 타임아웃 메세지를 전달하는 예제이다.**
~~~go
package main

import (
	"fmt"
	"time"
)

func main() {
	quit := make(chan struct{})
	done := process(quit)
	timeout := time.After(10 * time.Millisecond)

	select {
	case d := <-done:
		fmt.Println(d)
	case <-timeout:
		quit <- struct{}{}
		fmt.Println("Time Out...")
	}
}

func process(quit chan struct{}) chan string {
	done := make(chan string)

	go func() {
		ch := func() chan struct{} {
            ch := make(chan struct{})
            // 수행하고 싶은 명령어들을 여기에 적으면 됨
			go func() {
				time.Sleep(100 * time.Millisecond)
				ch <- struct{}{}
			}()
			return ch
		}()

		select {
		case <-ch: // 함수 실행이 먼저 끝나면 해당 수행문 실행
			done<-"complete"
		case <-quit: // 타임아웃이 먼저되면 
		}
	}()

	return done
}
~~~
~~~
실행결과

Time Out...
~~~