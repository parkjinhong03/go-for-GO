// rpc(remote procedure call)란 네트워크로부터 떨어져 있는 컴퓨터에서 코드를 실행하는 방식이다
// 이번 예제에서는 간단하게 표준 RPC 패키지늬 사용법을 살펴본다.

package server

import (
	"../contract"
	"fmt"
	"github.com/labstack/gommon/log"
	"net"
	"net/rpc"
)

const port = 8080

// 일단 REST API와 마찬가지로 RPC를 위해서도 핸들러를 정의해야 한다.
type HelloWorldHandler struct {}

// http.Handler와 이 핸들러의 차이점은 정의된 인터페이스가 딱히 없어서 이를 준수할 필요가 없다는 것이다.
// 구조채 팔드 및 관련 메서드를 가지고 있다면, 이것을 RPC 서버에 등록할 수 있다.
func (h *HelloWorldHandler) HelloWorld(args *contract.HelloWorldRequest, reply *contract.HelloWorldResponse) error {
	reply.Message = "Hello " + args.Name
	return nil
}

func main() {
	helloWorld := &HelloWorldHandler{}
	// rpc.Register() 함수를 이용해 위에서 선언한 인스턴스를 기본 RPC 서버에 등록한다.
	// 이를 통해 클라이언트가 해당 메서드를 호출할 수 있게 한다.
	rpc.Register(helloWorld)

	// net.Listen() 함수에 매개변수로 프로토콜과 포트를 넘겨 서버를 해당 IP주소에 바인딩할 수 있다.
	// 반환 값으로 Listener 인터페이스를 구현하는 인스턴스를 리턴한다.
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to listen on given port: %s", err))
	}
	// Close() 함수는 리스너를 닫아 연결을 기다리고 있던 Accept 동작에서 빠져나와 에러를 리턴한다.
	// defer 키워드를 사용하여 함수가 끝나기 직전에 l.Close() 함수를 실행시킨다.
	defer l.Close()

	// rpc 서버는 모든 연결에 대한 수신을 대기하지 않고 하나의 요청에 대해 응답이 완료되면 어플리케이션을 종료한다.
	// 따라서 무한 루프를 이용하여 하나의 요청이 완료되어도 끝나지 않고 다음 요청을 받아드릴 수 있는 상태로 만들 수 있다.
	for {
		// Accept() 함수는 다음 클라이언트가 연결될 때 까지 기다린다.
		// 클라이언트가 연결되면 해당 클라이언트와의 TCP 연결이 반환된다.
		conn, _ := l.Accept()
		// 그리고 rpc.ServeConn() 함수를 호출하여 등록시킨 핸들러를 실행시키며 연결을 마무리한다.
		go rpc.ServeConn(conn)
	}
}