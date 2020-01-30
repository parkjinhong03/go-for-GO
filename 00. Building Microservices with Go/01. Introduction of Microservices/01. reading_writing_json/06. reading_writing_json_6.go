// 이번에는 net/http 패키지에서 제공하는 여러 핸들러들을 이용하여 정적 파일을 반환하는 엔드 포인트를 만든다.

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 8080

	// http.FileServer 함수는 HTTP 요청에 대해 파일 시스템의 내용을 제공하는 핸들러를 리턴한다.
	// 매개변수로 http.Dir("./images")의 반환값을 넘겨줬으므로 전달받은 요쳥 경로가 '/cat.jpg' 이라면 ./images/cat.jpg 파일을 반환한다.
	// 만약 해당 경로의 파일이 존재하지 않으면 404 Not Found를 반환한다.
	fileHandler := http.FileServer(http.Dir("./images"))
	// http.Handle(pattern string, handler Handler) 함수를 통해 새 핸들러를 등록한다.
	// 요청 경로가 첫번째 매개변수("/images/")로 시작되는 모든 요청은 두번째 매개변수로 넘긴 핸들러의 ServeHTTP 메서드를 실행시킨다.
	// http.StripPrefix 함수는 요청 경로에서 매개 변수로 받은 접두사("/images/")를 제거한 다음 두번째 매개변수로 넘긴 핸들러를 호출한다.
	// 따라서 /images/cat.jpg 경로로 요청을 보내면 fileHandler에게는 /cat.jpg의 경로만 전달된다.
	http.Handle("/images/", http.StripPrefix("/images/", fileHandler))

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}