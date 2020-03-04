// 해당 파일은 프로그램의 실행 속도를 살펴보기 위해 가장 좋은 기법인 프로파일링을 하기 위한 파일이다.
// 프로파일링이란 어플리케이션을 실행하는 동안 자동으로 샘플 데이터를 수집하여 특정 기능 및 데이터를 통계 요약하여 계산하는 것 이다.
// Go는 다음과 같은 세 가지 유형의 프로파일링을 지원한다.
//	- CPU 		: CPU 실행 시간을 가장 많이 필요로 하는 프로세스를 식별한다.
//	- 힙(Heap)	: 메모리 할당을 가장 많이 하는 구문을 식별한다.
//	- 블로킹		: GoRoutine이 가장 많이 대기하도록 하는 동작을 식별한다.

package main

import (
	"../data"
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

// flag.String 함수를 이용하여 파싱할 커맨드라인 플래그를 등록할 수 있다.
// 아래의 코드는 플래그로 "cpuprofile=example"이 입력되면 "example"을 cpuprofile에 대입하겠다는 뜻이다.
// 하지만 만약 사용자가 cpuprofile 플래그를 입력하지 않으면 cpuprofile 변수에는 두 번째 매개변수로 넘긴 ""이 대입된다.
var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to file")
var memoryprofile = flag.String("memoryprofile", "", "Write memory profile to file")
var store *data.MongoStore


func main() {
	// flag.Parse 함수를 호출하면 flag.String 함수의 반환값을 담은 변수에 해당 함수로 등록한 플래그에 따라 값이 저장된다.
	// 참고로 이러한 프로세스가 가능한 이유는 flag.String 함수의 반환 타입이 *string이기 때문에 가능한 것 이다.
	flag.Parse()

	// cpuprofile 태그에 매칭된 값을 얻고 싶다면 flag.String 함수의 반환값을 담은 변수를 역참조 해야한다.
	if *cpuprofile != "" {
		fmt.Println("Running with CPU profile")
		// os.Create 함수를 이용하여 매개변수로 넘길 경로에 따라 파일을 만들고, 해당 파일과 연결된 os.File 타입의 객체를 반환한다.
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		// pprof.StartCPUProfile 함수를 이용하여 프로파일을 시작하고 초당 100회씩 프로세스를 중지하여 매개변수로 넘김 io.Writer에 데이터를 기록한다.
		pprof.StartCPUProfile(f)
	}

	// 아래에서 선언한 sigs 채널은 os.Signal 타입의 각종 시그널들이 송수신될 통로로, 특정 시그널이 송신될 때 까지 대기하기 위한 용도로 쓰일 것 이다.
	// 참고로 시그널이란 특정 이벤트가 발생되었을 때 프로세스에 전달하는 일종의 신호 같은 개념이다.
	sigs := make(chan os.Signal, 1)
	// 위에서 말한 것 처럼, signal.Notify 함수를 이용하면 특정 시그널 발생 시 위에서 선언한 채널로 해당 시그널을 송신시킬 수 있다.
	// 아래 코드는 syscall.SIGINT, syscall.SIGTERM 시그널이 발생하면 그 시그널을 sigs 채널로 송신시키겠다는 뜻이다.
	// 참고로 각각의 시그널들은 모두 16진수로, syscall.SIGINT(0x2)같은 경우는 ctrl+C로 인터럽트가 발생한 경우이다.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// 위에서 signal.Notify 함수로 등록한 특정 이벤트(시그널)가 발생되어 해당 이벤트가 sigs 채널로 송신될 때 까지 대기한다.
		<- sigs
	}()
}