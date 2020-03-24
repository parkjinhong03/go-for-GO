# **버퍼드 채널**(buffered channel)
## **버퍼드 채널이란?**
- **채널**은 지정한 크기의 **버퍼**를 가질 수 있는데, 그 **버퍼를 가진 채널**을 의미한다.
- 채널을 생성할 때 **버퍼의 크기**를 **make**의 두 번째 매개변수로 전달하면 **버퍼드 채널**을 만들 수 있다.
    ~~~go
    ch := make(chan int, 2)
    ~~~

<br>

---
## **버퍼드 채널의 동작**
- ### 버퍼드 채널은 **비동기 방식**으로 동작한다.
    - 채널이 **꽉 찰 때**까지 채널로 메세지를 **계속 전송**할 수 있음
    - 채널이 **빌 때** 까지 채널로부터 메세지를 **계속 수신**해 올 수 있음

<br>

- **버퍼드 채널 사용 예시**
    ~~~go
    package main

    import "fmt"

    func main() {
        c := make(chan int, 2)
        c <- 1
        c <- 2
        c <- 3
        fmt.Println(<-c)
        fmt.Println(<-c)
        fmt.Println(<-c)
    }
    ~~~

<br>


- 이 코드는 **버퍼가 꽉 참**에도 불구하고 **계속 송신**해서 다음과 같은 **에러**가 발생한다.
    ~~~
    fatal error: all goroutines are asleep - deadlock!
    ~~~

<br>

- 이 코드를 **다음과 같이** 수정해보자.
    ~~~go       
    package main

    import "fmt"

    func main() {
        c := make(chan int, 2)
        c <- 1
        c <- 2
        go func() {
            c <- 3
        }

        fmt.Println(<-c)
        fmt.Println(<-c)
        fmt.Println(<-c)
    }
    ~~~
    ~~~
    실행 결과

    1
    2
    3
    ~~~           

    <br>  

- 채널 c에 **세 번째 값**을 전송하는 부분을 **별도의 고루틴**으로 동작시켰다.
- 고루틴은 **채널 c**에 값을 전송할 수 있을때까지 **대기**하다가, 채널에 들어온 첫 번째 값을 **수신**해가는 **즉시** 바로 채널에 값을 **전송**한다.