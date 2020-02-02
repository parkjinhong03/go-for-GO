package client

import (
	"../contract"

	"fmt"
	"github.com/labstack/gommon/log"
	"net/rpc"
)

func CreateClient(port int) *rpc.Client {
	// HTTP 프로토콜의 연결을 생성해야하므로 rpc.DialHTTP() 메서드를 사용한다.
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	return client
}


func PerfumeRequest(client *rpc.Client, name string) contract.HelloWorldResponse {
	args := &contract.HelloWorldRequest{Name: name}
	var reply contract.HelloWorldResponse

	err := client.Call("HelloWorldHandler.HelloWorld", args, &reply)
	if err != nil {
		log.Fatal("error: ", err)
	}

	return reply
}