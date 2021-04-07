# **for**

**Go에서는 모든 반복문을 for 문으로 작성한다.**

- **기본 문법**
~~~go
for 초기화구문;조건식;후속작업 {
    ...
}
~~~

<br>

- 초기화 구문과 후속 작업은 **생략**할 수 있고, 각각은 서로 **세미콜론**으로 구분된다.
- 생략했을 경우에는 조건식이 **true가 될 때 까지** 내부 코드를 **반복**하는데, 다른 언어의 **while**문 처럼 동작한다.
~~~go
for {
    ...
}
~~~

<br>

- **for문의** 전체적인 프로세스 **실행 순서**는 다음과 같다
    - 반복문을 시작할 때, **초기화 구문**을 실행한다.
    - 조건식을 확인한 후, **조건식**이 **true**가 될 때 까지 내부 코드를 **반복**하여 **수행**한다.
    - 각 반복 작업 후에는 **후속 작업**을 실행한다.

~~~go
// 특정 횟수만큼 반복하는 작업은 다음과 같이 작성한다.
for i:=0; i<COUNT; i++ {
    ...
}
~~~

<br>

## **break와 continue**
- for문을 **강제 종료**해야 할 경우에는 **break** 키워드를 사용한다.
- 현재 수행하는 작업을 **건너띄고** 다음 반복 작업을 수행해야 할 때에는 **continue** 키워드를 사용한다.
~~~go
// 100부터 0사이의 홀수를 출력하는 코드

i := 100

for {
    if i < 0 {
        break
    }
    if i%2==0 {
        continue
    }

    fmt.Println(i)
}
~~~

<br>

## **for 문에 레이블 사용**
- for문 앞에 **콜론(:)으로 끝나는 문자**가 있으면 **레이블**로 인식한다.
- **continue, break, 레이블**을 함께 사용하면 반복문을 유연하게 제어할 수 있다.
~~~go
x := 7
table := [][]int{{1, 5, 9}, {2, 6, 5, 13}, {5, 3, 7. 4}}

LOOP:
for row:=0; row<len(table); row++ {
    for col:=0; col<len(table[row]); col++ {
        if table[row][col] == x {
            fmt.Println("found %d(row:%d, col:%d)\n", x, row, col)
            // LOOP로 지정된 for문을 빠져나옴
            break LOOP
        }
    }
}
~~~
~~~
실행 결과
fount 7(row:2, col:2)
~~~

<br>

- 또한, 다음 코드와 같이 **continue**에도 **레이블**을 사용할 수 있다.
~~~go
x := 5
table := [][]int{{1, 5, 9}, {2, 6, 5, 13}, {5, 3, 7, 4}}

for row:=0; row<len(table); row++ {
    LOOP:
    for col:=0; col<len(table[row]); col++ {
        if table[row][col] == x {
            fmt.Printf("found %d(row:%d, col:%d)\n", x, row, col)
            // LOOP로 지정된 for문을 다음 명령어를 수행
            break LOOP
        }
    }
}
~~~
~~~
실행 결과
found 5(row:0, col:1)
found 5(row:1, col:2)
found 5(row:2, col:0)
~~~

<br>

> 이러첨, 반복문 여러 개가 **중첩**될 떄 **레이블**을 사용하면 **가독성**이 높은 코드를 작성할 수 있다.