# **panic & revocer**
## **런타임 에러와 패닉**
- 실행 중에 에러가 발생하먄 **Go 런타임**은 **패닉을 발생**시킨다.
- 패닉이 발생하면 **패닉 에러 메세지**가 출력되고 프로그램이 종료된다.
- 물론 **패닉**이 발생해도 프로그램을 종료하지 않고 **계속해서 이어**나갈 수 있다.

<br>

**런타임 에러 예시 코드**
~~~go
package main

import "fmt"

func main() {
	fmt.Println(device(1, 0))
}

func device(a, b int) int {
	return a/b
}
~~~
~~~
실행 결과

panic: runtime error: integer divide by zero

goroutine 1 [running]:
main.device(...)
        /Users/parkjinhong/workspace/go/src/package_practice/pkg/panic/main.go:10
main.main()
        /Users/parkjinhong/workspace/go/src/package_practice/pkg/panic/main.go:6 +0x12
exit status 2
~~~
- 에러 상황이 심각해서 프로그램을 더 이상 실행할 수 없을 때 **panic() 함수**를 사용해 강제로 **패닉을 발생**시키고 **프로그램을 종료**할 수 있다.
- 주로 **error** 값을 **panic()** 함수의 **매개변수**로 전달한다.
~~~go
package main

import "fmt"

func main() {
	fmt.Println("Starting the program")
	panic("A severe error occurred: stopping the program")
	fmt.Println("Ending the program")
}
~~~
~~~
실행 결과
Starting the program
panic: A severe error occurred: stopping the program

goroutine 1 [running]:
main.main()
        /Users/parkjinhong/workspace/go/src/package_practice/pkg/panic/main.go:7 +0x95
exit status 2
~~~

<br>

---
## **패니킹(panicking)이란?**
- **panic()** 이 호출되면 실행을 **종료**하고 모든 **defer** 구문을 실행한 후 호출한 함수로 **패닉에 대한 제어권**을 넘긴다.
- 이러한 작업은 **함수 호출 스택**의 **상위 레벨**로 올라가며 계속 이어진다.
- 그리고 프로그램 스택의 **최상단인 main 함수**에서 프로그램을 **종료**하고 에러 상황을 **출력**한다.