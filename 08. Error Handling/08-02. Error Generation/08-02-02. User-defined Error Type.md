# **사용자 정의 에러 타입**
## **시용자 정의 에러 타입이란?**
- **에러**는 발생한 상황에서 **적절한 조치**를 취할 수 있게 **최대한 구체적**으로 만들어야 한다.
- 에러 메세지 및 **추가 정보**를 담아 에러 타입을 직접 만들면 **활용도**가 높다.
- **Error() string 메서드**만 가지고 있으면 모두 **에러로 사용**할 수 있으므로 **확장성** 또한 높다.

<br>

---
## **사용자 정의 에러 타입 사용 예시**
- **sqrtError**을 정의하여 앞서 구현한 **Sqrt()** 함수를 다시 작성해 보았다.
    ~~~go
    package main

    import (
        "fmt"
        "math"
        "time"
    )

    type SqrtError struct {
        time time.Time
        value float64
        message string
    }

    func (e SqrtError) Error() string {
        return fmt.Sprintf("[%v]ERROR - %s(value: %g)", e.time, e.message, e.value)
    }

    func Sqrt(f float64) (float64, error) {
        if f < 0 {
            return 0, SqrtError{time.Now(), f, "음수는 사용할 수 없습니다."}
        }
        if math.IsInf(f, 1) {
            return 0, SqrtError{time.Now(), f, "무한대 값은 사용할 수 없습니다."}
        }
        if math.IsNaN(f) {
            return 1, SqrtError{time.Now(), f, "잘못된 수 입니다."}
        }
        return math.Sqrt(f), nil
    }
    ~~~
- **SqrtError** 타입은 에러가 발생한 시간과 에러를 발생시킨 값, 에러 메세지를 담고 있다.
- **Error() 메서드**는 이렇게 담고있는 여러 값을 이용하여 **에러 메세지**를 생성하여 반환한다.
    ~~~go
    func main() {
        v, err := Sqrt(9) // 9의 제곱근, 정상 작동
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(v)

        v, err := Sqrt(-9) // -1의 제곱근, 에러 발생
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(v)
    }
    ~~~
    ~~~
    실행 결과

    3
    2019/11/24 23:41:21 [2019-11-24 23:41:21.08171 +0900 KST m=+0.000080033]ERROR - 음수는 사용할 수 없습니다.(value: -9)
    exit status 1
    ~~~

<br>

---
## **사용자 정의 에러 타입 확인**
- 함수나 메서드로 에러가 반환됬을 때 **특정 타입의 에러**인지 확인하는 **두 가지 방법**이 있다.

    1. **타입 어설션으로 확인**
        ~~~go
        if e, ok:= err.(SqrtError); ok {
            fmt.Println("Sqrt Error", e)
        }
        ~~~

    2. **switch 문으로 확인**
        ~~~go
        switch e := err.(type) {
        case SqrtError:
            fmt.Println("sqrt error", e)
        default
            fmt.Println("Default error", e)
        }
        ~~~