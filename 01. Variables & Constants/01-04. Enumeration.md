# **열겨형**(Enumeration)
- 열겨형은 **차례로 1씩 증가하는 상수의 묶음**이다.
- Go는 열거형을 따로 만들지 않고 **정수의 그룹**을 **상수로 선언**해서 **열거형**을 표현한다.
~~~go
const (
    Sunday = 0
    Monday = 1
    Tuesday = 2
    Wednesday = 3
    Thursday = 4
    Friday = 5
    Saturday = 6
)
~~~

<br>

- 상수를 열겨형으로 선언할 때 **iota 예약어**를 사용하면 편리하다.
- 상수를 그룹으로 묶어서 선언할 때 const 그룹에서 iota의 값은 **0**이고, 이후로는 **1씩 증가**한다.
- 아래 코드는 위의 코드와 같은 역할을 한다.
~~~go
const (
    Sunday = iota
    Monday
    Tuesday
    Wednesday
    Thursday
    Friday
    Saturday
)
~~~

<br>

- iota는 정수 뿐만 아니라 **부동소수점 타입**에도 사용할 수 있고 계산식과 **혼합**하여 샤용할 수 있다.
~~~go
type ByteSize int64

const (
    _ = iota
    KB ByteSize = 1 << (10 * iota)
    MD
    GB
    TB
    PB
    EB
)

const (
    DEFAULT_RAGE = 5 + 0.3 * iota
    GREEN_RATE
    SILVER_RATE
    GOLD_RATE
)
~~~

<br>

- iota로 **비트 상태**를 표현하는 상수를 정의할 수 있다.
~~~go
package main

import "fmt"

const (
	Running = 1 << iota		// 1 << 0 == 1
	Waiting					// 1 << 1 == 2
	Send					// 1 << 2 == 4
	Receive					// 1 << 3 == 8
)

func main() {
	stat := Running|Send
	display(stat)
}

func display(stat int) {
	if stat & Running == Running {
		fmt.Println("Running")
	}
	if stat & Waiting == Waiting {
		fmt.Println("Waiting")
	}
	if stat & Send == Send {
		fmt.Println("Send")
	}
	if stat & Receive == Receive {
		fmt.Println("Receive")
	}
}
~~~
~~~
실행 결과

Running
Send
~~~