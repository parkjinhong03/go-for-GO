// 마지막으로 기본 패키지에서 제공해주는 핸들러 대신 직접 핸들러를 만들어 등록시키는 것을 구현할 것 이다.
// 참고로 요청의 유효성 검사와 응답을 리턴하는 기능의 핸들러를 따로 분리하기 위해 두개의 핸들러를 만든다.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type helloWorldRequest3 struct {
	Name string `json:"name"`
}

type helloWorldResponse6 struct {
	Message string `json:"message"`
}

// 요청의 유효성을 검사하는 핸들러를 구현하기 위해 구조체를 정의한다.
// 유효성 검사 후 응답을 리턴해야 하기 때문에 다음 핸들러에 대한 참조를 위한 필드도 추가한다.
type validationHandler struct {
	next http.Handler
}

// http.Handler 인터페이스의 객체를 넘겨 유효성 검사 후 다음으로 실행시킬 핸들러를 등록시킨 validationHandler를 반환한다.
func newValidationHandler(next http.Handler) http.Handler {
	return validationHandler{next: next}
}

// validationHandler를 http.Handler로 사용하기 위해 ServeHTTP(http.ResponseWriter, *http.Request) 메서드를 정의한다.
func (h validationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request helloWorldRequest3
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}

	// 유효성 검사 완료 후 next 필드에 저장되어있던 핸들러의 ServeHTTP() 메서드를 실행시킴으로써 응답 프로세스의 소유권을 다음 핸들러로 넘긴다.
	h.next.ServeHTTP(rw, r)
}

// 유효성 검사 후 로직을 실행할 핸들러를 구현하기 위해 구조체를 정의한다.
type helloWorldHandler struct {}

func newHelloWorldHandler() http.Handler {
	return helloWorldHandler{}
}

// 단순히 구조체를 json으로 인코딩하여 응답을 리턴하는 기능의 함수이다.
func (h helloWorldHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	response := helloWorldResponse6{Message: "Hello"}

	encoder := json.NewEncoder(rw)
	encoder.Encode(response)
}

// 기능별로 핸들러를 분리한 이러한 구조의 코드는 마이크로서비스의 취지에 알맞은 구조이며 이 기법 자체는 유용하다.
// 하지만 이렇게 단순한 경우에는 코드를 불필요하게 복잡하게 만들며 실제로 코드의 반복을 줄이지 않는다.
// 또한 http.Handler 인터페이스를 손상시키지 않고 서로 다른 핸들러끼리 서로 데이터를 주고 받을 수 있는 방법이 없다.
// 따라서 다음에는 Context를 이용하여 이러한 단점을 보완하여 조금 더 나은 코드를 작성해 볼 것 이다.
func main() {
	port := 8080

	handler := newValidationHandler(newHelloWorldHandler())

	http.Handle("/helloworld", handler)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}