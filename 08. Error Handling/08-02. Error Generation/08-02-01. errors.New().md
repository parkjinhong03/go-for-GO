# **에러 생성**
- **수행 결과**로 에러가 반환될 떄도 있지만, 에러를 **의도적으로 발생**시켜야 하는 특별한 **경우도 있다.**
- 따라서 이번 강에서는 **에러를 생성하는 방법**을 소개한다.


## **errors.New()**
- 에러를 생성하는 가장 간단한 방법은 **errors 패키지**의 **New() 함수**를 사용한다는 것 이다.
- errors 패키지의 **errors.New()** 함수가 어떻게 작성되어 있는지 **직접 확인해 보자.**
    ~~~go
    // errors 패키지 소스
    package errors

    func New(text string) error {
        return &errorString{text}
    }

    type errorString struct {
        s string
    }

    func (e *errorString) Error() string {
        return e.s
    }
    ~~~

<br>

- **errors.New()** 함수를 사용하면 **errorString 구조체**가 반환된다.
- 이 구조체에는 **Error() 메서드**가 정의되어 있으므로 **error 타입**으로 사용할 수 있다.
- 이처럼 에러를 **의도적**으로 발생시켜야 할 때는 **errors.New() 함수**로 에러를 생성한다.

<br>

**errors.New() 함수로 에러 생성 및 출력 예제**
~~~go
package main

import (
	"errors"
	"fmt"
)

func main() {
	errNotFound := errors.New("Not found error")

	fmt.Println("errors:", errNotFound)
	fmt.Println("errors:", errNotFound.Error())
}
~~~
~~~
실행 결과

errors: Not found error
errors: Not found error
~~~
- **errNotFound**를 출력한 결과와 **errNotFound.Error()** 를 출력한 결과가 같다.
- 즉, **에러를 출력**하면 **Error()** 메서드의 결과 **문자열이 출력**됨을 알 수 있다.

<br>

**이것을 함수 내부에 적용해 보자**
~~~go
package main

import (
	"errors"
	"fmt"
	"math"
)

func Sqrt(f float64) (float64, error) {
	if f < 0 {
		return 0, errors.New("음수는 사용할 수 없습니다.")
	}
	return math.Sqrt(f), nil
}

func main() {
	if f, err := Sqrt(-1); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Println(f)
	}
}
~~~
~~~
실행 결과

Error: 음수는 사용할 수 없습니다.
~~~

<br>

---
## **fmt.Errorf()**
- **fmt.Errorf()** 는 에러가 발생한 값과 매개변수의 정보를 담아 **에러 메세지**를 만들 수 있다.
- **서식**이 적용된 문자열과 **여러 개의 매개변수**로 문자열을 만드는 방식은 **fmt.Printf()** 와 유사하다.
- 단, 메세지를 화면에 출력하는 게 아니라 **메세지를 담고 있는 error 타입 값**을 만들어낸다.  

<br>

**fmt.Errorf() 함수**
~~~go
func Errorf(format string, a ...interface{}) error {
    return errors.New(Sprintf(format, a...))
}
~~~

<br>

- **fmt.Errorf()** 함수에서는 만들어진 문자열을 **errors.New()** 의 매개변수로 전달해 **error값**을 만들어낸다.
- 앞서 작성한 Sqrt() 함수에 fmt.Errorf()를 적용하여 에러를 생성해 보자.
    ~~~go
    func Sqrt(f float64) (float64, error) {
        if i < 0 {
            return 0, fmt.Errorf("음수(%d)는 사용할 수 없습니다", f)
        }
        return math.Sqrt(f), nil
    }
    ~~~