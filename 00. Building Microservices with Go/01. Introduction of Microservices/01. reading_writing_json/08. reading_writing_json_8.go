// 앞서 언급했던 단점들을 보완하기 위해, 우리는 표준 패키지인 context를 사용할 것 이다.
// 이 context 패키지는 Go 1.7에 표준 패키지에 포함되었다.
// Context 타입은 여러 Go 루틴에서 동시에 안전하게 접근하여 수정 및 사용할 수 있다.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type helloWorldRequest4 struct {
	Name string `json:"name"`
}

type helloWorldResponse7 struct {
	Message string `json:"message"`
}

type validationHandler1 struct {
	next http.Handler
}

func newValidationHandler1(next http.Handler) http.Handler {
	return validationHandler1{next: next}
}

// 단순 문자열 사용시의 키 충돌을 방지하기 위하여 문자열 타입의 객체를 선언한다.
type validationContextKey string

// 이전 코드와 다른점은 context를 추가하고 덧붙이는것 밖에 없다.
// 하지만 이를 통해 Go 루틴이 객체에 안전하게 접근하도록 할 수 있다.
func (h validationHandler1) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request helloWorldRequest4

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}

	// context.WithValue() 함수는 기존의 context를 매개변수로 받아 데이터를 추가하여 새로 생성한 context를 반환한다.
	// 따라서 아래 코드는 r.Context()의 반환값인 context에 매개변수로 받은 키와 값을 추가한 것 이다.
	ctx := context.WithValue(r.Context(), validationContextKey("name"), request.Name)
	// WithContext() 함수를 통해 *http.Request의 인스턴스인 객체에 새로운 ctx라는 content를 추가할 수 있다.
	r = r.WithContext(ctx)

	h.next.ServeHTTP(rw, r)
}

type helloWorldHandler6 struct {}

func newHelloWorldHandler1() http.Handler {
	return helloWorldHandler6{}
}

// validationHandler1 핸들러에서 http.Request의 context에 데이터를 추가한 후 호출되는 핸들러이다.
func (h helloWorldHandler6) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// Context 객체의 메서드안 Value() 메서드를 이용하여 Context에 저장한 값을 빼올 수 있다.
	// 참고로 해당 메서드의 리턴값은 interface{} 타입이므로 type assertion을 통해 원하는 타입으로 바꿔야한다.
	name := r.Context().Value(validationContextKey("name")).(string)
	response := helloWorldResponse7{Message: "Hello " + name}

	encoder := json.NewEncoder(rw)
	encoder.Encode(response)
}

func main() {
	port := 8080

	handler := newValidationHandler1(newHelloWorldHandler1())

	http.Handle("/helloworld", handler)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}