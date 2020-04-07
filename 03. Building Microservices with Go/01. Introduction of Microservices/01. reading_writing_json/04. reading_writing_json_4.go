// 이전까지 Go의 구조체를 json으로 인코딩하는 방법을 알아봤다면, 지금부터는 json을 Go의 구조체로 디코딩하는 방법을 알아볼 것 이다.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// JSON을 디코딩하여 저장할 구조체를 정의한다.
// 대소문자를 구분하지 않지만 정확히 일치하는 것이 좋다.
// 하지만 소문자로 시작하는 구조체 필드는 무시하므로 주의해야한다.
type helloWorldRequest1 struct {
	Name string `json:"name"`
}

type helloWorldResponse4 struct {
	Message string `json:"message"`
}

func helloWorldHandler4(w http.ResponseWriter, r *http.Request) {
	// Request.Body는 io.ReadCloser 인터페이스를 구현하고 있으며 []Byte나 문자열을 리턴하지 않는다.
	// 따라서 Request.Body의 값을 얻고 싶으면 ioutil.ReadALL() 함수를 사용하면 바이트 배열로 읽을 수 있다.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad reqeust", http.StatusBadRequest)
		return
	}

	// json 값을 담을 구조체를 선언한다.
	var request helloWorldRequest1
	// json.Unmarshal() 함수를 이용하여 request 객체에 바이트 배열인 body의 값을 집어넣는다.
	// 참고로 request가 아닌 &request를 주는 이유는, 함수에 request 객체의 복사본인 아닌 참조값을 넘겨주기 위해서다.
		err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Bad reqeust", http.StatusBadRequest)
		return
	}

	response := helloWorldResponse4{Message: "Hello " + request.Name}
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

func main() {
	port := 8080

	http.HandleFunc("/helloworld", helloWorldHandler4)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}