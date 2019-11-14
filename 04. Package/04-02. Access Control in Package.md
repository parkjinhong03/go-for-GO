# **접근 제어**

## **Go에서의 접근 제어**
- 앞에 예제에서 본 IsDigit **함수**처럼 **변수, 상수, 함수, 사용자 정의 타입, 메서드, 구조체**의 필드 등의 식별자가 **대문자**로 시작하면 패키지 **외부**에서 **접근** 할 수 있다.

- 하지만, 식별자가 **소문자**로 시작하면 패키지 **외부에서 접근할 수 없다.** 단, 패키지 내부에서는 모든 요소에 접근할 수 있다.

<br>

~~~go
// $GOPATH/src/package_practice/lib/lib.go
package lib

func IsDigit(c int32) bool {
    return '0' <= c && c <= '9'
}

func isSpace(c int32) bool {
	switch c {
	case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		return true
	}
	return false
}
~~~
~~~go
// $GOPATH/src/package_practice/pkg/main.go
package main

import (
    "fmt"
    "package_practice/lib"
)

func main() {
    lib.IsDigit('1')
    lib.IsDigit('a')
    lib.isSpace('\t') // 소문자로 시작하는 함수는 패키지 외부에서 사용 불가능
}
~~~
~~~
실행 결과

# command-line-arguments
./main.go:11:14: cannot refer to unexported name lib.isSpace
./main.go:11:14: undefined: lib.isSpace

~~~