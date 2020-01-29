// 구조체 필드의 속성을 이용하여 인코딩되는 JSON의 속성을 변경하는 코드

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// 01. reading_writing_json_1.go에 나와있는 구조체와는 다르게 필드의 속성(태그)을 추가했다.
// 이렇게 필드의 속성으로 `json:"message"`를 달면 JSON으로 인코딩될 때 {"message": "..."}으로 된다.
// 애초에 필드의 이름을 message로 하지 않는 이유는 Go의 문법과 연관되어 있다.
// Go에서는 소문자로 된 프로퍼티(멤버 변수)는 외부에서 접근할 수 없으므로 Marshal 함수에서도 무시하기 때문이다.
// 따라서 Go의 낙타표기법과 JSON의 팟홀표기법을 연결하고 싶다면, 구조체애 태그를 설정해주면 된다.
type helloworldResponse2 struct {
	Message string `json:"message"`
}

func helloWorldHandler2(w http.ResponseWriter, r *http.Request) {
	response := helloworldResponse2{Message: "HelloWorld"}
	data, err := json.Marshal(response)
	if err != nil {
		panic("Oops")
	}

	// Response Body에 {"message":""HelloWorld"}가 찍히는 것을 볼 수 있다.
	fmt.Fprint(w, string(data))
}

func main() {
	port := 8080

	http.HandleFunc("/helloworld", helloWorldHandler2)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}