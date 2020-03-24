package server

import (
	"../contract"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloWorldHandler struct {}

func (h *HelloWorldHandler) HelloWorld(args *contract.HelloWorldRequest, reply *contract.HelloWorldResponse) error {
	reply.Message = "Hello " + args.Name
	return nil
}

// http.ResponseWriter 및 http.Request에는 io.ReadWriteCloser 인터페이스를 구현하는 타입이 없기 때문에 직접 정의해야 한다.
type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (int, error) { return c.in.Read(p) }
func (c *HttpConn) Write(p []byte) (int, error) { return c.out.Write(p) }
func (c *HttpConn) Close() error { return nil }

func StartServer(port int) {
	handler := &HelloWorldHandler{}
	rpc.Register(handler)

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to listen on given port %v", port))
	}
	defer l.Close()

	// 03. rpc_http/server/server.go 와는 다르게 두 번째 매개변수에 nil이 아닌 실제 핸들러를 전달하였다.
	// 참고로 일반 함수를 http.HandlerFunc 타입으로 type casting 하는데 성공하면 자동으로 ServeHTTP 메서드가 생기므로 http.Handler 인터페이스로 사용할 수 있다.
	// 따라서 요청이 들어올 때 마다 두 번째 매개변수로 넘긴 핸들러 객체의 ServeHTTP(= httpHandler) 함수가 실행된다.
	http.Serve(l, http.HandlerFunc(httpHandler))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	// NewServerCodec() 함수는 RPC 응답 작성 및 요청 읽기를 구현한 rpc.ServerCodec 타입의 객체를 반환한다.
	// 위에서 정의한 ReadWriteCloser 인터페이스를 구현하는 타입인 HttpConn 객체를 매개변수로 념긴다.
	serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
	// rpc.ServeRequest() 함수에 위의 함수로부터 얻은 rpc.ServerCodec 객체를 매개변수로 넘겨 응답을 작성하고 에러 처리 후 연결을 종료한다.
	err := rpc.ServeRequest(serverCodec)
	if err != nil {
		log.Printf("Error while serving JSON request: %v", err)
		http.Error(w, "Error while serving JSON request, details have been logged.", 500)
		return
	}
}