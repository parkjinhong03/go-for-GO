// 앞에서 말했던 두 가지의 SOP 정책 해결 방안중 대부분의 사용 비중을 차지하고 있는 CORS(Cross-Origin Resource Sharing) 이다.
// HTTP 클라이언트가 실제 요청 전에 URI에 OPTIONS 요청을 보내어 응답받은 헤더에 따라 CORS를 지원하고 있는 API인지 브라우저에서 구분한다.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type helloWorldResponse struct {
	Message string `json:"message"`
}

func main() {
	port := 8080

	http.Handle("/helloWorld", http.HandlerFunc(helloWorldHandler))

	log.Print(fmt.Sprintf("Server starting on port %v\n", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func helloWorldHandler(rw http.ResponseWriter, r *http.Request) {
	// 만약 요청 메서드가 OPTIONS라면, 즉 CORS를 지원하는 API인지 확인하기 위한 요청이 들어왔다면 아래 구문들을 실행한다.
	if r.Method == "OPTIONS" {
		// 다음 코드로 인해 모든 도메인에서 해당 API로 요청을 할 수 있도록 브라우저가 허락해준다.
		rw.Header().Add("Access-Control-Allow-Origin", "*")
		// 다음 코드로 인해 다른 도메인에서 GET 메서드로는 해당 API로 요청을 할 수 있도록 브러우저가 허락해준다.
		rw.Header().Add("Access-Control-Allow-Methods", "GET")
		// 응답 Body가 존재하지 않기 때문에 Status Code를 204 No Content로 지정하여 OPTIONS 메서드에 대한 요청을 완료한다.
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	response := helloWorldResponse{Message: "Hello World"}
	encoder := json.NewEncoder(rw)
	encoder.Encode(response)
}