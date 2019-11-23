# **저수준 제어**
## **저수준 제어란?**
- **Go**에는 고루틴과 체널 외에도 **병행 프로그래밍**을 위한 **저수준 제어 기능**이 있다.
- **sync 패키지**에는 **뮤텍스**(mutex)로 **공유 메모리**를 **제어**할 수 있는 기능이 있다.
- **sync/atomic 패키지**에는 **원자성**을 보장할 수 있는 **연산**(add, compare, swap 등)이 있다.

<br>

---
## **sync.Mutex**
- **뮤텍스**는 여러 **고루틴**에서 **공유**하는 데이터를 **보호**해야 할 때 사용한다.

- 뮤텍스 구조체는 **다음 함수**를 제공한다.
    - **func (m \*Mutex) Lock()**    -> 뮤텍스 잠금
    - **func (m \*Mutex) Unlock()**  -> 뮤텍스 잠금 해제

- **임계 영역**의 코드를 실행하기 전에는 **Lock() 메서드**로 잠금을 **걸고** 처리 완료 후에는 **Unlock() 메서드**로 잠금을 **해제**한다.

<br>

---
## **고루틴에서 뮤텍스를 사용하지 않고 공유 데이터에 접근하는 코드**
~~~go
package main

import (
	"fmt"
	"runtime"
)

type counter struct {
	i int64
}

func (c *counter) increment() {
	c.i += 1
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
			c.increment()
			done <- struct{}{}
		}()
	}

	// 모든 고루틴이 완료될 때까지 대기
	for i:=0; i<1000; i++ {
		<-done
	}

	c.display()
}
~~~
~~~
실행 결과

980 ~ 999 사이 랜덤 수
~~~
- 이 프로그램은 **고루틴**으로부터 **incrememnt() 메서드**를 **1000**번 실행시켜 **1000**이 **출력**되어야 한다.
- 하지만 실제로 프로그램을 실행시켜 보면 **1000보다 작은 값**이 출력된다.
- 이는 **여러 고루틴**이 구조체 c의 내부 필드 i의 값을 **동시에 수정**하려 해서 **경쟁상태**가 만들어져서 이다.
- 따라서 **공유 데이터**인 내부 필드 값 i를 변경하는 부분을 **뮤텍스**로 **보호**해줘야 한다.

<br>

### *Note - 예제에 사용된 Go 기본 라이브러리 함수*
- runtime.GOMAXPROCS()  -> 현재 프로그램에서 사용할 CPU의 최대 수
- runtime.NumCPU()      -> CPU 코어 수

<br>

---
## **뮤텍스로 공유 데이터를 보호해주는 코드**
~~~go
package main

import (
	"fmt"
	"runtime"
	"sync"
)

type counter struct {
	i int64
	mu sync.Mutex // 공유 데이터 i를 보호하기 위한 뮤텍스
}

func (c *counter) increment() {
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

1000
~~~

<br>

- **counter 구조체** 정의 시 **내부 필드 값 보호**를 위해 **sync.Mutex**를 **내부 필드** mu로 정의하였다.
- increment() 메서드에서 **c.mu.Lock()** 과 **c.mu.Unlock()** 을 작성하여 **변경**하는 부분을 **뮤텍스**로 **보호**해주었다.
- 따라서 프로그램을 실행시켜보면 **고루틴 수행 수만큼** 값이 정확하게 **증가**한 것을 확인할 수 있다.