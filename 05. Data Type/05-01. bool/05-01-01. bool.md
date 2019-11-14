# **불 (boolean)**

## **불(boolean) 타입이란?**
- **불 타입**은 **true나 false** 값을 가진다.
- 값으로 **true나 false**를 직접 사용할 수 있고, **논리 및 비교 연산자**의 결괏값을 사용할 수 있다.

~~~go
var v int = 10

b1 := true
b2 := v == 5    // false
b3 := v == 10   // true
b4 := b1 && b2  // false
b5 := b1 || b2  // true
b6 := !b1       // false
~~~

<br>

---
## **조건부 논리 연산자**
- **조건부 논리 연산자**(&&, ||)는 **단락**(short circuit) 방식으로 동작한다.

- 즉, **앞의 조건**만으로 **결과**를 얻을 수 있다면 **뒤의 조건을 무시**한다.
    - b1 || b2 -> b1 == true -> b2 검사 X -> true 반환
    - b1 && b2 -> b1 == false -> b2 검사 X -> false 반환

<br>

---
## *참고*
- Go에서는 다른 타입 값을 불 타입으로 바꾸는 **암묵적인 변환을 하지 않는다.**
- 즉, **0이나 nil을 false로 변환하지 않는다.**
~~~go
// 다음과 같은 코드는 오류를 발생시킨다.

func checkValue(v interface{}) ( 
    // 빌드 오류 발생 (non-bool v used as if condition)
    if v {
        fmt.Println("value is %v\n", v)
        return
    }
    fmt.Println("value is null")
)
~~~

<br>

- if 문의 **조건식**에는 다음과 같이 반드시 **불 타입 값**을 사용해야 한다.
~~~go
func checkValue(v interface{}) ( 
    if v != nil {
        fmt.Println("value is %v\n", v)
        return
    }
    fmt.Println("value is null")
)
~~~