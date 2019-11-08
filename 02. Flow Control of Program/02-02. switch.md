# **switch**

- **기본 형태**
~~~go
switch 값 {
    case 조건1:
        ...
    case 조건2:
        ...
    default:
        ...
}
~~~
<br>

- if문과 마찬가지로 여는 괄호는 **switch**와 **같은 줄**에 있어야한다.
- 또한 **,(콤마)** 를 이용해서 case뒤에 값을 **여러 개** 쓸 수 있다.
- **Go**에서는 다른 언어와 다르게 기본으로 case 구문 마지막에 **break**가 실행된다.
~~~go
var i int = 1

switch i (
    case 0: // i가 0일 때는 아무것도 수행하지 않음
    case 1: // i가 1일 때만 수행됨
        func()
)
~~~
<br>

- Go는 **fallthrough**를 이용하여 case 구문에서 **다음 case** 구문으로 **넘어가게** 할 수 있다.
- 하지만 다음 case로 넘겨서 한번에 처리하는 경우는 case에 조건을 여러 개 넣는 것이 낫다.
~~~go
var i int = 1

switch i (
    case 0: 
        fallthrough
    case 1: // i가 1일 때와 2일 때 모두 수행함
        func()
)
~~~

<br>

- if문과 마찬가지로 switch 문에서도 **초기화 구문**을 사용할 수 있다.
- **초기화 구문**에 선언된 변수는 **switch 문 내에서만** 사용할 수 있다.
~~~go
switch a, b := x[i], y[j]; {
    case a < b:
        fmt.Println("a는 b보다 작습니다.")
    case a == b:
        fmt.Println("a는 b와 같습니다.")
    case a > b:
        fmt.Println("a는 b보다 큽니다.")
}
~~~