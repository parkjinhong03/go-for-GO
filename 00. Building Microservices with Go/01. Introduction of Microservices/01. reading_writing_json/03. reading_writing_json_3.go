// 이번에는 json.Marshal()이 아닌 다른 방식으로 구조체를 JSON으로 인코딩하는 방법이 나와있다.
// 이전에 나온 코드들을 보면 구조체를 바이트 배열로 디코딩한 다음 문자열로 변환시켜 응답 스트림에 쓴다.
// 이러한 방식의 인코딩은 특별히 효율적으로 보이지도 않고 실제로도 효율적이지 않다.
// 따라서 Go에서 지원하는 스트림에 직접 쓸 수 있는 인코더 및 디코더를 알아볼 것 이다.

package _1__reading_writing_json

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type helloworldResponse3 struct {
	Message string `json:"message"`
}

func helloWorldHandler3(w http.ResponseWriter, r *http.Request) {
	response := helloworldResponse3{Message: "HelloWorld"}
	// http.ResponseWriter가 io.Writer 인터페이스 만족하고 있기 때문에 json.NewEncoder() 함수의 매개변수로 사용할 수 있다.
	encoder := json.NewEncoder(w)
	// encoder.Encode() 메서드를 실행하면 내부에서 다음과 같은 동작들이 실행된다.
	// 1. 매개변수로 받은 구조체를 Marshal 한다.
	// 2. 그 결과를 바이트 배열에 저장하지 않고 바로 HTTP 응답에 쓴다.
	// 3. 바이트 배열로 마샬링하는 것 보다 50% 이상의 속도 향상 결과를 볼 수 있다.
	encoder.Encode(&response)
}

// 참고로 이러한 표준 패키지 동작 방식의 이해는 프레임워크에 대해 자세히 이해하기 위해 알아야 할 것 들이다.
func main() {
	port := 8080

	http.HandleFunc("/helloworld", helloWorldHandler3)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}