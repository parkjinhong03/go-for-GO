// rpc 통신 형식의 서버를 구동하고 클라이언트를 생성하여 요청을 보내는 예제 코드이다.
// 각각의 패키지(server, client)에 대한 설명은 해당 패키지에 더 자세하게 나와 있다.

package main

import (
	"./client"
	"./server"
	"fmt"
)

func main() {
	// 고루틴을 이용하여 서버를 localhost:8080에 바인딩한다.
	go server.StartServer()

	// client.CreateClient() 함수를 이용하여 위에서 구동한 서버에 대한 클라이언트 생성
	c := client.CreateClient()
	defer c.Close()

	// client.PerformRequest() 함수를 이용하여 서버에 보낼 요청을 생성한다.
	reply := client.PerformRequest(c, "RPC")
	fmt.Println(reply.Message)
}