// 구조체를 이용하여 JSON으로 마샬링하여 출력하는 첫 번째 방법

package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
)

// json으로 마샬링 하여 출력할 구조체를 정의한다
type helloWorldResponse1 struct {
	Message string
}

// Go의 구조체를 JSON으로 마샬링하여 출력하는 간단한 핸들러 정의
func helloWorldHandler1(w http.ResponseWriter, r *http.Request) {
	// 위에서 정의한 구조체의 인스턴스를 만들고 원하는 메세지를 설정한다.
	response := helloWorldResponse1{Message: "HelloWorld"}
	// 리턴하기 전에 json.Marshal() 함수를 사용해 문자열로 인코딩한다.
	// 참고로 json.Marshal() 함수가 반환하는 데이터의 자료형은 []byte 이다.
	// 그 반환된 값을 string 자료형으로 타입 캐스팅을 하면 {"Message":"HelloWorld"}이 된다.
	data, err := json.Marshal(response)
	// 만약 err가 nil이 아닐 경우, 즉 인코딩 과정에서 문제가 생겼을 경우 바로 panic 상태로 들어가 서버가 중단된다.
	if err != nil {
		panic("Oops")
	}

	// fmt.Fprint() 함수를 이용하여 http.ResponseWriter에 string으로 바꾼 json 값을 출력함으로써 응답이 완료된다.
	fmt.Fprint(w, string(data))
}

func main() {
	// 서버를 바인딩 시킬 포트 번호를 8080으로 지정한다.
	port := 8080

	// 위에서 정의한 핸들러 함수를 /helloworld 결로로 라우팅한다.
	http.HandleFunc("/helloworld", helloWorldHandler1)

	// 서버 구동을 시작할 때 로그를 찍기 위하여 log.Printf() 함수 사용
	log.Printf("Server starting on port %v\n", port)

	// http.ListenAndServe() 함수의 인자값에 ":8080"을 넣음으로써 8080의 포트로 서버를 바인딩한다.
	// 서버 구동 시 에러가 발생하였을 때 로그를 찍기 위하여 log.Fatal() 함수 사용
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}