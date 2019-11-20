# **타입 변환**
**인터페이스는 어떤 값이라도 동적으로 할당받을 수 있지만 변수에 할당된 값을 실제 타입으로 변환해야 할 때가 종종 있다.** 

<br>

---
## **인터페이스 변수의 실제 타입으로의 변환 방법**
1. ### **타입 어설션**  

2. ### **switch 문**  

3. ### **reflect 패키지**

<br>

---
## **타입 어설션으로 타입 변환**
- **타입 어셜션**(Type Assertion)은 다음과 같이 수행한다.
    ~~~go
    v := interface.(Type)
    ~~~
- **타입 변환 성공** -> 해당 타입으로 변환된 값 반환
- **타입 변환 실패** -> 런타임 오류(panic()) 발생
    ~~~go
    // 다음과 같은 방식으로 안전하게 타입을 변활할 수 있다.
    if v, ok := interface.(Type) ok {
        // ...
    }
    ~~~

<br>

---
## **switch 문으로 타입 변환**
- 인터페이스 타입 변환 시 변환할 타입을 **확실히 알고** 있다면 **타입 어셜션**으로 변환하면 된다.
- 하지만 실제 타입이 무엇인지 **확실하지 않을** 때는 **switch 문**으로 타입을 확인하는 것이 좋다.

<br>

#### 다음 checkType() 함수에서는 switch 문으로 타입을 먼저 확인한 후 실제 값과 타입을 출력한다.
~~~go
func checkType(v interface{}) {
    switch v.(type) {
        case bool:
            fmt.Printf("%t is a bool\n", v)
        case int, int8, int16, int32, int64:
            fmt.Printf("%d is an int\n", v)
        case uint, uint8, uint16, uint32, uint64:
            fmt.Printf("%d is a unsigned int\n", v)
        case float64:
            fmt.Printf("%f is a float64\n", v)
        case complex64, complex128:
            fmt.Printf("%f is a complex\n", v)
        case string:
            fmt.Printf("%s is a string\n", v)
        case nil:
            fmt.Printf("%v is nil\n", v)
        default:
            fmt.Printf("%v is unknown type\n", v)
    }
}
~~~
~~~
3 is an int
1.500000 is a float64
(1.000000+5.000000i) is a complex
true is a bool
s is a string
{} is unknown type
<nil> is nil
~~~