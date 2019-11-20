# **기본 라이브러리 인터페이스의 사용**
## **fmt.String**
- 다음은 Go의 **기본 라이브러리**인 fmt 패키지에서 **Println() 함수**를 정의하는 부분이다.
    ~~~go
    func Println(a ...interface{}) (n int, err error) {
        return FPrintln(os.Stdout, a...)
    }

    func FPrintln(w io.Writter, a ..interface{}) (n int, err error) {
        ...
        // fmt.Stringer 인터페이스 타입일 때 String() 메서드의 결괏값을 출력
    }
    ~~~

<br>

- **os.Stdout**을 첫 번째 **매개변수**로 해서 **FPrintln() 함수**로 전달하고,
- **FPrintln() 함수** 내부에서는 **fmt.String 인터페이스**로 문자열을 출력한다.
    ~~~go
    // fmt 패키지에서 String 인터페이스를 정의하는 부분.
    type Stringer interface {
        String() string
    }
    ~~~

<br>

- 즉, **fmt.String 인터페이스**에 정의된 **String() 메서드**를 가지면 기본 출력 명령인 **fmt.Println() 함수**로 출력될 문자열을 **지정**할 수 있다.

<br>

## *예제*
- 이제 **Item, DiscountItem, Renal, Items**에 **String() 메서드**를 추가하여 **fmt.Println() 함수**로 출력해보자.
- Item의 String() 메서드는 fmt.Sprintf()로 **각 요소의 문자열을 조합**해서 만든다.
~~~go
func (t Item) String() string {
	return fmt.Sprintf("[%s] %0.f", t.name, t.Cost())
}

func (t DiscountItem) String() string {
	return fmt.Sprintf("%s => %.0f(%.0f%s DC)", t.Item.String(), t.Cost(), t.discountRate, "%")
}

func (t Rental) String() string  {
	return fmt.Sprintf("[%s] %0.f", t.name, t.Cost())
}

func (ts Items) String() string {
	var s []string
	for _, t := range ts {
		s = append(s, fmt.Sprint(t))
	}
	return fmt.Sprintf("%d items.total: %.0f\n\t-%s", len(ts), ts.Cost(), strings.Join(s, "\n\t-"))
}
~~~

<br>

- **fmt.Println() 함수**를 실행하면 **String() 메서드**에서 지정한 문자열이 출력된다.
~~~go
func main() {
	shoes := Item{"Sports Shoes", 30000, 2}
	eventShoes := DiscountItem{shoes, 10.00}
	videos := Rental{"Interstellar", 1000, 3, Days}
	items := Items{shoes, eventShoes, videos}

	fmt.Println(shoes)
	fmt.Println(eventShoes)
	fmt.Println(videos)
	fmt.Println(items)
}
~~~
~~~
실행 결과

[Sports Shoes] 60000
[Sports Shoes] 60000 => 54000(10% DC)
[Interstellar] 3000
3 items.total: 117000
        -[Sports Shoes] 60000
        -[Sports Shoes] 60000 => 54000(10% DC)
        -[Interstellar] 3000
~~~

<br>

---
## **io.Writer**
- 다음 io의 **Writer 인터페이스**와 fmt의 **Fprintln() 함수**는 Go 프로그래밍 시 자주 쓰인다.

    - **io.Writer 인터페이스**
        ~~~go
        type Writer interface {
            // 매개변수 p를 받아 byte 길이와 에러 상태 반환 메서드
            Write(p []byte) (n int, err os.Error)
        }
        ~~~

    - **fmt.Fprintln() 함수**
        ~~~go
        // 첫 번째 매개변수 w에 두 번째 매개변수 a를 전달하는 기능을 한다.
        func FPrintln(w io.Writer, a ...interface{}) (n int err error)
        ~~~

        <br>

### **이 둘로 인터페이스가 어떻게 활용되는지 확인해 보자.**
- 매개변수로 받은 **io.Writer**에 **msg**를 **fmt.FPrintln(g)** 로 출력하는 함수를 정의하자.
    ~~~go
    func handle(w io.Writer, msg string) {
        fmt.Fprintln(w, msg)
    }
    ~~~

<br>

- 다음 코드는 **io.Writer**로 **os.Stdout**을 전달하여 msg를 **명령 프롬프트**에 출력한다.
    ~~~go
    func main() {
        msg := []string{"This", "is", "an", "example", "of", "io.Writer"}

        for _, s := range msg {
            time.Sleep(100 * time.Millisecond)
            handle(os.Stdout, s)
        }
    }
    ~~~

<br>

- 이번에는 **handle() 함수**에 다른 **io.Writer**를 전달해 보자.
- **main 함수**에서 HTTP 웹 서버를 동작시키고 **io.Writer**로 **http.ResponseWriter**를 전달한다.
    ~~~go
    func main() {
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            handle(w, r.URL.Path[1:])
        })

        fmt.Println("start listening on port 4000")
        http.ListenAndServe(":4000", nil)
    }
    ~~~

<br>

- **fmt** 패키지와 **http** 패키지는 **연결고리**가 없다.
- 하지만 전달하는 값에 **Write() 메서드**가 정의되어 있으면, **io.Writer**로 사용할 수 있다.

<br>

---
## **정리**
- 이러한 특징들 덕분에 **모듈 간 연계**를 쉽게 할 수 있고, **동적**인 느낌으로 프로그래밍을 할 수 있다.
- 인터페이스에 정의된 메소드는 **코드 패턴의 컨벤션**(관습) 역할을 하고, 컨벤션을 따라 하면 다른 모듈과 **쉽게 호환**할 수 있다.
- 따라서 **컴파일러의 컨벤션 보장**을 받으며 **동적 코딩**을 할 수 있다.