# **인터페이스**(interface)

## **인터페이스란?**
- **인터페이스**(interface)의 역할은 **객체의 동작**을 표현하는 것이다.
- 각 타입이 실제로 내부에 **어떻게 구현**되어 있는지는 말하지 않고, 단순히 **동작 방식만 표현**한다.
- 따라서 **인터페이스**의 이러한 특징은 **추상 메커니즘**을 제공한다.

    - **매개변수**로 인터페이스를 사용하는 것은 **"값의 타입이 무언인지"** 보다
    - **"값이 무엇을 할 수 있는지"** 에만 집중하게 해준다.

<br>

---
## **Go에서의 인터페이스**
- **Go**의 인터페이스는 **덕 타이핑 방식**을 채택했다.
    - **덕 타이핑 방식** -> 객체의 변수나 메서드의 집합이 **객체의 타입**을 결정하는 것
    - 즉, 인터페이스에 **걸맞은 요소**(메서드)를 가진 타입은 **인터페이스로 사용**할 수 있다.

- 인터페이스의 이러한 특징 덕분에 **정적 타입 언어**인 Go를 **유연**하게 사용할 수 있다.

<br>

---
## **인터페이스 정의**
- 다음과 같이 **메서드 서명**을 묶어 **하나의 인터페이스**로 정의한다.
    ~~~go
    type 인터페이스 명 interface {
        메서드1 (매개변수) 반환타입
        메서드2 (매개변수) 반환타입
    }
    ~~~
- 인터페이스에 정의된 메서드와 **서명이 같은 메서드**를 가진 타입은 **인터페이스로 사용**할 수 있다.

- **다음 코드**를 보면 이해가 될 것 이다.
    - **area() 메서드**를 가진 **shaper 인터페이스**를 정의한다.
    - 그리고 **shaper 인터페이스**를 **매개변수**로 받은 후
    - shaper에 있는 **area() 메서드**의 실행 결과를 출력하는 **describe() 함수**를 정의한다.

    ~~~go
    type shaper interface {
        area() float
    }

    func describe(s shaper) {
        fmt.Println("area:", s.area())
    }
    ~~~

    <br>

    - **rect 구조체**와 **area() 메서드**를 정의한다.
    - **main 함수**에서 **rect 타입**의 값을 **매개변수**로 전달하여 **describe() 함수**를 실행한다.
    ~~~go
    type rect struct { width, height float64 }

    func (r rect) area() float64 {
        return r.width * r.height
    }

    func main() {
        r := rect{3, 4}
        describe(r)
    }
    ~~~
    ~~~
    실행 결과
    
    area: 12
    ~~~

- **rect 구조체**와 **shaper 인터페이스**는 코드에서 아무런 **연결 고리**가 없다.
- 하지만 shaper에서 정의한 **메서드**를 rect 타입이 **제공**하면, 해당 타입을 **인터페이스로 사용**할 수 있다.

<br>

---
## ***Go 코드 컨벤션***
- 인터페이스 이름은 **메서드 이름**에 **er**(또는 r)을 붙여서 짓는다. (ex. Printer, Reader 등등)

- 또한 인터페이스는 **짧은 단위**로 만든다
    - Go의 기본 라이브러리에도 메서드를 **하나만 정의한 인터페이스**가 대부분이다.
    - 많아도 **세 개**를 넘지 않게 한다.

- Go의 **기본 라이브러리**에 정의된 인터페이스를 살펴보자.
    - **io 패키지의 Reader 인터페이스**
    ~~~go
    type Reader interface {
        Read(p []byte) (n int, err error)
    }
    ~~~

    - **io 패키지의 Writer 인터페이스**
    ~~~go
    type Writer interface {
        Write(p []byte) (n int, err error)
    }
    ~~~

    - **io 패키지의 Closer 인터페이스**
    ~~~go
    type Closer interface {
        Close() error
    }
    ~~~

<br>

---
## **익명 인터페이스**
- **인터페이스**도 **타입**을 정의하지 않고 **익명**으로 사용할 수 있다.
    ~~~go
    func display(s interface{ show() }) {
        s.show()
    }
    ~~~

- **display()** 함수에는 **show() 메서드**를 가진 타입만 **매개변수**로 전달 할 수 있다.

<br>

**익명 인터페이스를 이용한 예제 코드**
~~~go
type rect struct{ width, height float64 }

func (r rect) show() {
    fmt.Printf("width: %f, height: %f\n", r.width, r.height)
}

type circle struct { radius float64 }

func (c circle) show() {
    fmt.Printf("radius: %f\n", c.radius)
}

func display(s interface{ show() }) {
    s.show()
}

func main() {
    r := rect{3, 4}
    c := circle{2.5}
    display(r)          // width: 3, height: 4
    display(c)          // radius: 2.5
}
~~~

<br>

---
## **빈 인터페이스**
- **interface{}** 타입은 **메서드**를 정의하지 않은 인터페이스이다.
- **interface{}** 에는 **정의한 메서드**가 없어서 **어떤 값이라도** interface{}로 사용할 수 있다.
- 즉, 함수나 메서드의 **매개변수**를 interface{} 타입으로 정의하면 **어떤 값이든** 전달받을 수 있다.
~~~go
func main() {
    r := rect{3, 4}
    c := circle{2.5}
    display(r)              // width: 3.000000, height: 4.000000
    display(c)              // radius: 2.500000
    display(2.5)            // 2.5
    display("rect struct")  // rect struct
}

func display(s interface{}) {
    fmt.Println(s)
}
~~~