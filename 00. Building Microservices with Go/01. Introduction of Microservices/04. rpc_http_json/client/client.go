package client

import (
	"../contract"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func PerfumeRequest(port int) contract.HelloWorldResponse {
	// http.Post() 메서드를 사용하여 http 서버에 POST 요청을 보낼 수 있다.
	r, _ := http.Post(
		fmt.Sprintf("http://localhost:%v", port),
		"application/json",
		bytes.NewBuffer([]byte(`{"id": 0, "method": "HelloWorldHandler.HelloWorld", "params": [{"name":"rpc_http_json"}]}`)),
	)
	defer r.Body.Close()

	var response contract.HelloWorldResponse
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&response)

	return response
}