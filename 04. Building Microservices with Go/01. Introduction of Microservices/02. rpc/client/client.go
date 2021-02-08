// http 프로토콜을 사용하지 않으므로 요청을 보낼 때 단순하게 curl을 사용할 수 없다.
// 따라서 이번 예제에서 rpc 서버의 연결 클라이언트를 구현해볼 것 이다.

package client

import (
	"../contract"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/rpc"
)

const port = 8080

// 해당 함수는 rpc 서버와 연결을 생성하고 그 연결을 반환하는 기능을 한다.
func CreateClient() *rpc.Client {
	client, err := rpc.Dial("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	return client
}

// 해당 함수는 연결될 클라이언트를 이용하여 서버로 요청을 보낸 후 응답을 리턴하는 기능을 한다.
func PerformRequest(client *rpc.Client, name string) contract.HelloWorldResponse {
	// 서버에 등록된 함수를 호출할 때 필요한 매개변수를 위해 요청 인스턴스를 생성한다.
	args := &contract.HelloWorldRequest{Name: name}
	// 서버에 등록된 함수의 호출 결과 값을 담을 응답 인스턴스를 생성한다.
	var reply contract.HelloWorldResponse

	// rpc 서버를 바인딩할 때 HelloWorldHandler.HelloWorld() 함수를 등록했기 때문에 해당 이름을 매개변수로 넘겨 함수를 호출한다.
	err := client.Call("HelloWorldHandler.HelloWorld", args, &reply)
	if err != nil {
		log.Fatal("error: ", err)
	}

	return reply
}