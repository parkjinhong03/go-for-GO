// 01. 마이크로 서비스의 소개의 마지막 예제는 JSON-RPC 표준 직렬화 및 역직렬화이다.
// net/rpc/jsonrpc 패키지를 사용하여 해당 기능을 구현해볼 것 이다.

package main

import (
	"./client"
	"./server"
	"fmt"
)

const port = 8080

func main() {
	go server.StartServer(port)

	reply := client.PerfumeRequest(port)
	fmt.Println(reply.Message)
}