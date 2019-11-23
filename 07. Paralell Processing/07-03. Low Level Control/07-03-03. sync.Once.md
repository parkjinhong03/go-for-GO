# **sync.Once**
## **sync.Once란?**
- 특정 함수를 **한 번만 수행**해야 할 때 **sync.Once**를 사용한다.
- **sync.Once** 구조체는 다음과 같은 **메서드**를 제공한다.
    ~~~go
    func (o *Once) Do(func())
    ~~~
- **한 번만 수행할 함수**를 **Do() 메서드**의 매개변수로 전달하면 **여러 고루틴**에서 **실행**한다 해도 해당 함수는 **한 번만 수행**된다.

<br>

---
## **sync.Once 사용 예시 코드**
~~~go
package main

import (
	"fmt"
	"runtime"
	"sync"
)

const initialValue = -500

type counter struct {
	i int64
    mu sync.Mutex // 공유 데이터 i를 보호하기 위한 뮤텍스
    once sync.Once // 한 번만 수행할 함수를 지정하기 위한 Once 구조체
}

func (c *counter) increment() {
    // i 값 초기화 작업은 한 번만 수행되도록 once의 Do() 메서드로 실행 
    c.once.Do(func(){
        c.i = initialValue
    })

	c.mu.Lock() // i 값을 변경하는 부분을 뮤텍스로 잠굼 
	c.i += 1   // 공유 데이터 변경
	c.mu.Unlock() // i 값을 변경 완료한 후 뮤텍스 잠금 해제
}

func (c *counter) display() {
	fmt.Println(c.i)
}

func main() {
	// 모든 CPU를 사용하게 함
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := counter{i:0} // 키운터 생성
	done := make(chan struct{})	// 완료 신호 수신용 채널

	// c.increment()를 실행하는 고루틴 1000개 실행
	for i:=0; i<1000; i++ {
		go func() {
			c.increment()   // 카운터의 값을 1 증가시킴
			done <- struct{}{}  // done 채널에 완료 신호 전송
		}()
	}

	// 모든 고루틴이 완료될 때까지 대기
	for i:=0; i<1000; i++ {
		<-done
	}

	c.display() // c의 값 출력
}
~~~
~~~
실행 결과

500
~~~

<br>

- 다음 코드에서는 카운터의 내부 필드 값 i의 **초기화 작업**을 **Once 구조체**의 **Do() 메서드**로 지정했다.
- 따라서 **c.increment()** 는 **1000번 수행**됬기 때문에, 프로그램 **출력 값**은 **500**인 걸 확인할 수 있다.