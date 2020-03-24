# **메서드**
- **메서드**는 **사용자 정의 타입 값**에 **호출**할 수 있는 특별한 함수이다.
- **리시버 타입 변수**에 메서드 **호출** -> **변수**가 메서드 내부로 **전달** -> 메서드 내부에서 **사용 가능**

<br>

---
## **리시버 값 전달 방식**
- 함수와 마찬가지로 **메서드**는 **값에 의한 호출**이 기본 방식이다.
- 따라서 메서드 내부에서 리시버 변수의 값을 변경하려면 **변수의 메모리 주소를 전달**해야 한다.
- 참조에 의한 호출로 주소를 전달하려면 **리시버 변수 타입에 * 를 사용**하여 **포인터**로 지정해야 한다.

<br>

~~~go
type quantity int

func (q quantity) greaterThan(i int) bool {
    return int(q) > i
}

func (q *quantity) increment() { *q++ }

func (q *quantity) decrement() { *q-- }

func main() {
    q := quantity(3)
    q.increment()
    fmt.Printf("Is q(%d) greater than %d? %t\n", q, 3, q.greaterThan(3))
    q.decrement()
    fmt.Printf("Is q(%d) greater than %d? %t\n", q, 3, q.greaterThan(3))
}
~~~
~~~
실행 결과

Is q(4) greater than 3? true
Is q(3) greater than 3? false
~~~

<br>

- **05-06-02 장**에서 말했듯이, **필드를 많이 갖는 구조체**라면 복사하기에 **시스템 리소스**가 많이 소요될 것 이다.
- 따라서 리시버 변수의 값을 변경하지 않아도 **리시버를 포인터**로 지정하는 **습관**을 들이는 것이 좋다.

<br>

---
## **리시버 변수 생략**
- 메서드 내부에서 **리시버 변수**를 **사용하지 않을 때**도 있다.
~~~go
type rect struct {
    width float64
    height float64
}

func (rect) new() rect {
    return rect{}
}

func main() {
    r := rect{}.new()
    fmt.Println(r)
}
~~~
- 메서드 내부에서 리시버 변수를 사용하지 않는다면 메서드 정의시 **리시버 변수를 생략**할 수 있다.

<br>

---
## **메서드 함수 표현식**
- **메서드**도 **일급 객체**이기 때문에 **변수에 힐딩**할 수 있고 **함수의 매개변수**로 전달할 수 있다.**
- 메서드의 함수 표현식은 **메서드의 리시버**를 **첫 번째 매개변수**로 전달하는 함수이다.

~~~go
// 다음 예제에서는 rect와 area()와 resize() 메서드를 함수 표현식으로 변환하여 사용한다.

package main

import "fmt"

type myRect struct {
	width, height float64
}

func (r myRect) area() float64 {
	return r.width * r.height
}

func (r *myRect) resize(w, h float64) {
	r.width += w
	r.height += h
}

func main() {
	r := myRect{3, 4}
	fmt.Println("area:", r.area())	// area: 12

	r.resize(10, 10)
	fmt.Println("area:", r.area())	// area: 182

	/* area() 메서드의 함수 표현식
	서명: func(rect) float64 */
	areaFn := myRect.area

	/* resize() 메서드의 함수 표현식
	서명: func(*rect, float64, float64) */
	resizeFn := (*myRect).resize

	fmt.Println("area:", areaFn(r))	// area: 182
	resizeFn(&r, -10, -10)
	fmt.Println("area:", areaFn(r))	// area: 12
}
~~~

<br>

- 이 예제에서는 다음 두 메서드를 **함수 표현식**으로 사용했다.
    - **areaFn** 메서드 -> **func(rect) float64**  
    **첫 번째 매개변수**로 **rect**를 받는다.

    - **resizeFn** 메서드 -> **func(\*rect, float64, float64)**  
    **첫 번째 매개변수**로 **\*rect**를 받고, **두 번째와 세 번째 매개변수**로 **float64**를 받는다.

- 메서드 자체를 **다른 함수의 메개변수**로 전달해야 할 때 메서드를 **함수 표현식**으로 변환하여 전달한다.