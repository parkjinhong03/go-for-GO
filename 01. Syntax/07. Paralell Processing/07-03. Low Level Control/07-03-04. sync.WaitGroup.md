# **sync.WaitGroup**
## **sync.WaitGroup란?**
- **sync.WaitGroup**은 **모든 고루틴**이 **종료**될 때 까지 **대기**해야할 때 사용한다.

- 다음은 **sync.WaitGroup**이 제공하는 **메서드**이다.
    - **func (wg \*WaitGroup) Add(delta int)** -> 대기중인 고루틴 개수 추가
    - **func (wg \*WaitGroup) Done()** -> 대기 중인 고루틴의 수행이 종료되는 것을 알려줌
    - **func (wg \*WaitGroup) Wait()** -> 모든 고루틴이 종료될때 까지 대기

<br>

---

## **sync.WaitGroup를 사용한 예제 코드**
~~~go
package main

import (
	"fmt"
	"runtime"
	"sync"
)

type counter struct {
	i int64
	mu sync.Mutex
	once sync.Once
}

const initialValue = -500

func (c *counter) increment() {
	c.once.Do(func() {
		c.i = initialValue
	})

	c.mu.Lock()
	c.i += 1
	c.mu.Unlock()
}

func (c *counter) display() {
	fmt.Println(c.i)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := counter{i:0} 
	wg := sync.WaitGroup{} // WaitGroup 생성

	for i:=0; i<1000; i++ {
		wg.Add(1) // WaitGroup 고루틴 갯수 1 증가
		go func() {
			defer wg.Done() // 고루틴 종료 시 Done() 처리
			c.increment()
		}()
	}

	wg.Wait() // 모든 고루틴이 완료될 때까지 대기

	c.display()
}
~~~
~~~
실행 결과

500
~~~

<br>

- **WaitGroup**을 생성한 뒤 고루틴이 시작될 떄 **Add 메서드**로 대기해야 하는 **고루틴 개수를 추가**한다.
- 그리고 **고루틴**이 **종료**될 때 **Done() 메서드**로 고루틴이 종료되었다고 **알려준다.**
- **Wait() 메서드**를 호출하면 대기중인 **모든 고루틴**이 **종료**될 때까지 **기다린다.**
- 주의할 점은 **Add** 메서드로 추가한 **고루틴 갯수**와 **Done** 메서드를 **호출한 횟수**는 **같아**야 된다. 