// Encoder와 같이 Unmarshal 함수 대신 Decoder를 사용하면 속도의 성능을 향상할 수 있다.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type helloWorldRequest2 struct {
	Name string `json:"name"`
}

type helloWorldResponse5 struct {
	Message string `json:"message"`
}

func helloWorldHandler5(w http.ResponseWriter, r *http.Request) {
	// 먼저 JSON 값을 디코딩하여 담을 Go 구조체부터 선언한다.
	var request helloWorldRequest2
	// io.Reader 인터페이스를 만족하는 r.Body를 매개변수로 넘겨서 새 decoder를 생성한다.
	// 내부적으로 r.Body의 데이터를 추출하여 Go의 구조체에 디코딩 할 준비를 마친다.
	decoder := json.NewDecoder(r.Body)
	// Decode() 메서드를 실행 하면 매개변수로 넘긴 구조체에 NewDecoder() 호출 시 넘긴 매개변수의 추출된 값이 대입된다.
	// 참고로 주소값을 넘기지 않으면 json: Unmarshal(non-pointer main.helloWorldRequest2) 에러가 발생한다.
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	response := helloWorldResponse5{Message: "Hello " + request.Name}
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

// 위와 같은 방식으로 디코딩을 진행하면, 더 적은 양의 코드로 33%의 성능을 향상시킬 수 있다.
func main() {
	port := 8080

	http.HandleFunc("/helloworld", helloWorldHandler5)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}