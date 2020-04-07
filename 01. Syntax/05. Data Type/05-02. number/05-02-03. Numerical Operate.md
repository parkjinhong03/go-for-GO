# **숫자 연산**
## **숫자 연산의 기본적인 규칙**
- **숫자 연산**(산술 연산, 비교 연산)은 **타입이 같은 숫자**끼리만 할 수 있다.
- **타입이 다른 숫자**끼리 **연산**하려면 반드시 **타입을 변환**해주어야 한다.

~~~go
// 타입이 다른 숫자끼리 연산하려고 하면 컴파일 에러가 발생한다.

i := 100000
j := int16(10000)
k := uint8(100)

fmt.Println(i + j)  // 컴파일 에러
fmt.Println(i + k)  // 컴파일 에러
fmt.Println(j > k)  // 컴파일 에러
~~~
~~~
실행 결과

# command-line-arguments
./main.go:20:16: invalid operation: i + j (mismatched types int and int16)
./main.go:21:16: invalid operation: i + k (mismatched types int and uint8)
./main.go:22:16: invalid operation: j > k (mismatched types int16 and uint8)
~~~

<br>

---
## **타입 캐스팅**
- 따라서 반드시 **같은 타입**으로 **변환**한 후에 **연산**을 해 줘야 한다.
~~~go
i := 100000
j := int16(10000)
k := uint8(100)

fmt.Println(i + int(j)) // 110000
fmt.Println(i + int(k)) // 100100
~~~
- 타입 변환 연산은 **오류를 발생시키지 않는다.**
- 하지만 원래 값보다 **작은 범위**를 다루는 타입으로 변환 시, **실제 값**은 기대했던 결과와 **다를 수 있다.**
~~~go
i := 100000

fmt.Println(int16(i)) // -31072
fmt.Println(uint8(i)) // 160
~~~
- 그러므로 타입을 변환하려면 **변환할 수 있는 값인지 먼저 점검**한 후에 변환해야 한다.

<br>

~~~go
// 타입 변환을 위한 안전한 코드

package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(intToUint8(100))
	fmt.Println(intToUint8(1000))
}

func intToUint8(i int) (uint8, error) {
	if 0 <= i && i <= math.MaxUint8 {
		return uint8(i), nil
    }
    // fmt.Errorf() -> 문자열을 기반으로 error를 만들어 반환해줌
	return 0, fmt.Errorf("%d cannot convert to unit8 type", i)
}
~~~
~~~
실행 결과

100 <nil>
0 1000 cannot convert to unit8 type
~~~

<br>

---
## **증감 연산자(++, --)**
- **Go**에서 증감 연산자는 **후치 연산**으로만 사용되고 반환 값은 없다.
- 이는 증감 연산자의 **과도한 사용**으로 **코드의 가독**성이 떨어지는 것을 막기 위해서다.

<br>

**다음 코드**는 모두 Go에서 **사용할 수 없는 코드**이다.
~~~go
++x             // 컴파일 에러
y = x++         // 컴파일 에러
print(x++)      // 컴파일 에러
print(++x)      // 컴파일 에러
if x++ > 0 {    // 컴파일 에러
    ...
}
~~~