package server

import (
	"../contract"
	"fmt"
	"github.com/labstack/gommon/log"
	"net"
	"net/http"
	"net/rpc"
)

type HelloWorldHandler struct {}

func (h *HelloWorldHandler) HelloWorld(args *contract.HelloWorldRequest, reply *contract.HelloWorldResponse) error {
	reply.Message = "Hello " + args.Name
	return nil
}

func StartServer(port int) {
	handler := &HelloWorldHandler{}
	rpc.Register(handler)
	// rpc.HandleHTTP() 함수를 이용하여 쉽게 HTTP 프로토콜을 사용할 수 있다.
	// 해당 메서드는 서버에게 두 개의 API를 Default Path로 생성해준다.
	// 1. "/_goRPC_"   -> 클라이언트를 통해 해당 경로로 요청을 보내면 rpc 서버에서 등록한 핸들러 함수가 실행된다.
	// 2. "/debug/rpc" -> 웹 브라우저에서 요청을 보내면 위의 엔드포인트에 대한 호춯 횟수를 확인할 수 있다.
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to listen on given port: %s", err))
	}
	defer l.Close()

	log.Printf("Server starting on port %v\n", port)

	// http.Serve() 메서드를 호출하여 서버 바인딩을 하여 요청 받을 준비를 완료한다.
	// 해당 메서드는 이전에서 REST 규칙의 API를 구현하였을 때 사용한 메서드와 동일한 기능의 메서드이다.
	http.Serve(l, nil)
}