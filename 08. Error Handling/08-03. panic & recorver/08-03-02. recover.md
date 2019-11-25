# **recover**
## **recover란?**
- **recover()** 함수는 **패니킹** 작업으로부터 **프로그램의 제어권**을 다시 얻어 계속 **이어나가**게 해준다.
- **recover()** 는 반드시 **defer 안**에서 사용 해야한다.
- defer 구문 안에서 recover()를 호출하면 **패닉 내부 상황**을 **error** 값으로 **복원**할 수 있다.
- **recover()** 로 패닉을 복원한 후에는 **패니킹 상황이 종료**되고 함수 **반환 타입의 제로값**이 **반환**된다.

<br>

**recover 사용 예시 코드**
~~~go
func main() {
    fmt.Println("result:", , divide(1, 0))
}

func divide (a, b int) (int error) {
    defer func(){
        if err := recover(); err != nil {
            fmt.Println(err)
        }
    }()

    return a/b
}
~~~
~~~
실행 결과

runtime error: integer divide by zero
result: 0
~~~

<br>

---
## **recover** VS **try-finally**
- **recover**는 파이썬의 **try-finally** 블록과 유사하다.
- 다음 코드는 매개변수로 받은 함수를 실행한 후 **패닉**이 발생하면 **error 정보**를 가져와 **출력**한다.
    ~~~go
    package main

    import (
        "fmt"
        "log"
    )

    func protect(g func()) {
        defer func() {
            log.Println("done")

            if err := recover(); err != nil {
                log.Printf("run time panic: %v", err)
            }
        }()
        log.Println("start")
        g()
    }

    func main() {
        protect(func() {
            fmt.Println(divide(4, 0))
        })
    }

    func divide(a, b int) int {
        return a/b
    }

    ~~~
- **protect()** 함수에서 **예상치 못한 패닉**에 의해 프로그램이 비정삭적으로 **종료**되는 것을 **막을** 수 있다.

<br>

---
## **recover** VS **try-except-finally**
- 다음은 **panic, defer, recover**를 사용한 예제이다.
    ~~~go
    package main

    import "fmt"

    func badCall() {
        panic("bad end")
    }

    func test() {
        defer func() {
            // except
            if e:=recover(); e!=nil {
                fmt.Printf("Panicking %s\r\n", e)
            }
        }()
        // try
        badCall()
    }

    func main() {
        fmt.Printf("calling test\r\n")
        test()
        // finally
        fmt.Printf("Test completed\r\n")
    }
    ~~~
    ~~~
    실행 결과

    calling test
    anicking bad end
    Test completed
    ~~~
- 이는 파이썬의 **try-except-finally** 블록과 **유사한 효과**를 준다.
- **defer-panic-recover** 패턴은 **패닉**을 복구하여 명시적인 **에러값**으로 **반환**해준다.