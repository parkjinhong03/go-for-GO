// HTTP 상에서 서로 다른 도메인에 대한 요청을 할 때, 이를 SOP(Same-Origin Policy) 정책으로 보안상 제한을 한다.
// 이 문제를 해결하기 위해서 나온 방법이 두 가지 있는데, 바로 JSONP와 CORS가 있다.
// 그중 JSONP는 CORS가 활성화 되기 전에 사용됐던 방법으로, script 태그는 SOP 정책에 적용되지 않는다는 것을 이용한 것 이다.
// 하지만 보안상의 이슈 때문에, 지금은 W3C에서 채택된 CORS 방식의 HTTP 통신을 사용하여 SOP를 피하는 것이 좋다.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type jsonpResponse struct {
	Message string `json:"message"`
}

func main() {
	port := 8080

	http.Handle("/helloWorld", http.HandlerFunc(jsonpHandler))

	log.Print(fmt.Sprintf("Server starting on port %v\n", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// 참고로 JSONP 형태로 반환될 데이터에 대한 요청을 나타내기 위해 일반적으로 URL 매개변수에 "callback=함수이름"을 추가한다.
func jsonpHandler(rw http.ResponseWriter, r *http.Request) {
	response := jsonpResponse{Message: "Hello World"}
	data, err := json.Marshal(response)
	if err != nil {
		panic("Oops!")
	}

	// URL 매개변수에서 callback 값을 가져온다.
	callback := r.URL.Query().Get("callback")
	if callback != "" {
		// 응답 바디 내용이 더 이상 json이 아니고 javascript이므로 Content-Type을 application/javascript로 설정한다.
		rw.Header().Add("Content-Type", "application/javascript")
		// 매개변수로 받은 함수 이름으로 json 값을 괄호로 싸서 반환하는 것으로 JSONP 구현을 마친다.
		fmt.Fprintf(rw, "%s(%s)", callback, string(data))
	} else {
		fmt.Fprint(rw, string(data))
	}
}