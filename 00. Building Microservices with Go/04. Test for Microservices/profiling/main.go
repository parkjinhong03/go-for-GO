// 해당 파일은 프로그램의 실행 속도를 살펴보기 위해 가장 좋은 기법인 프로파일링을 하기 위한 파일이다.
// 프로파일링이란 어플리케이션을 실행하는 동안 자동으로 샘플 데이터를 수집하여 특정 기능 및 데이터를 통계 요약하여 계산하는 것 이다.
// Go

package main

import (
	"../data"
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"os"
	"runtime/pprof"
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
	if *cpuprofile == "" {
		fmt.Println("Running with CPU profile")
		// os.Create 함수를 이용하여 매개변수로 넘길 경로에 따라 파일을 만들고, 해당 파일과 연결된 os.File 타입의 객체를 반환한다.
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
}