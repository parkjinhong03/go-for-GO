// 03. rpc_http 에서는 rpc 방식의 통신에서 HTTP 프로토콜을 사용하는 것을 구현한다.
// 코드의 대부분은 02. rpc와 비슷하므로 동일한 코드의 주석은 제외하였다.

package main

import (
	"./client"
	"./server"
	"fmt"
)

const port = 8080

func main()  {
	go server.StartServer(port)

	c := client.CreateClient(port)
	defer c.Close()

	reply := client.PerfumeRequest(c, "rpc_http")

	fmt.Println(reply.Message)
}