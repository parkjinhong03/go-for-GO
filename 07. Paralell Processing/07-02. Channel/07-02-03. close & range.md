# **close & range**
## **close**
- 채널에 더 이상 전송할 값이 없으면 **채널을 닫을 수 있다.**
- 하지만, 만약에 채널이 **빈 상태**가 아니라면 채널이 **비워진 후**에야 비로서 **close** 된다.
    ~~~go
    close(ch)
    ~~~
- 채널을 **닫은 후**에 메세지를 **전송**하면 **에러**가 발생한다.
    ~~~go
    package main

    import "fmt"

    func main() {
        ch := make(chan int, 2)
        ch<-1
        ch<-2
        go func(){ ch<-3 }()
        fmt.Println(<-ch)
        fmt.Println(<-ch)
        fmt.Println(<-ch)
        close(ch)
        ch<-4
    }
    ~~~
    ~~~
    실행 결과

    1
    2
    3
    panic: send on closed channel

    goroutine 1 [running]:
    main.main()
            /Users/parkjinhong/workspace/go/src/package_practice/pkg/goroutine_practice/main.go:28 +0x4c2
    exit status 2
    ~~~
- 또한 **채널의 수신자**는 채널에서 값을 읽을 때 채널이 **닫힌 상태**인지 아닌지 **두 번째 매개변수**로 확인할 수 있다.
    ~~~go
    package main

    import "fmt"

    func main() {
        ch := make(chan int, 2)
        ch <- 1

        if v, ok := <-ch; ok {
            fmt.Println(v)
        }

        close(ch)

        if v, ok := <-ch; ok {
            fmt.Println(v)
        } else {
            fmt.Println("channel was closed")
        }
    }
    ~~~
    ~~~
    실행 결과

    1
    channel was closed
    ~~~

<br>

---
## **range**
- **for i := range ch**은 채널 c가 닫힐 때 까지 **반복**하며 채널로부터 **수신**을 시도한다.
    ~~~go
    package main

    import "fmt"

    func main() {
        ch := make(chan int, 3)

        ch <- 1
        ch <- 2
        ch <- 3
        go func() {
            ch<-4
            close(ch)
        }()


        for i := range ch{
            fmt.Println(i)
        }

        if v, ok := <-ch; ok {
            fmt.Println(v)
        } else {
            fmt.Println("channel was closed")
        }
    }
    ~~~
    ~~~
    실행 결과

    1
    2
    3
    4
    channel was closed
    ~~~

- 또한, **for range** 구문 전에 ch를 **닫지 않으면 에러**가 발생한다.
    ~~~go
    package main

    import "fmt"

    func main() {
        ch := make(chan int, 3)

        ch <- 1
        ch <- 2
        ch <- 3


        for i := range ch{
            fmt.Println(i)
        }
    }
    ~~~
    ~~~
    실행 결과

    fatal error: all goroutines are asleep - deadlock!

    goroutine 1 [chan receive]:
    main.main()
            /Users/parkjinhong/workspace/go/src/package_practice/pkg/goroutine_practice/main.go:30 +0x398
    exit status 2
    ~~~

<br>

---
## *정리*
- 채널을 닫아주는 것은 **필수가 아니다.**
- 하지만 **for range**와 같이 수신자가 **채널에 더 이상 들어올 값이 없다는 것**을 알아야 할 때만 **채널을 닫아**주면 된다.