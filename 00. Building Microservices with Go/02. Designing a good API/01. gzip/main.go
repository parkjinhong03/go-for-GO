// REST 엔드 포인트는 가능항 경우 gzip과 deflate 인코딩은을 항상 지원해야 한다.
// gzip, deflate란 html, javascript, css 등을 압축해줘서 리소스를 받는 로딩시간을 줄여주게 해주는 응용 소프트웨어이다.
// 그중 gzip을 이용하여 response writer를 작성하는 것을 구현해볼 것 이다.

package main

import (
	"./handlers"
	"./middlewares"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 8080

	http.Handle("/helloWorld", middlewares.NewGzipMiddleWare(handlers.NewHelloWorldHandler()))

	log.Println(fmt.Sprintf("Server starting on port %v\n", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
