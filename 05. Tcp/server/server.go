package server

import "net"

// 서버의 동작을 명확하게 정의하기 위해 정의한 인터페이스
type ChatServer interface {
	// Listen 메서드는 외부에서 들어오는 연결을 수신한다.
	Listen(address string) error
	// StartServer와 CloseServer 메서드는 서버를 구동하고 종료시킨다.
	StartServer()
	CloseServer()

	// 클라이언트와의 연결을 수락하고 해당 클라이언트에 대한 정보를 반환한다.
	accept(conn net.Conn) *client
	// 연결시킨 클라이언트에 대한 비즈니스 로직을 실행시킨다.
	serve(client *client)
	// 다른 모든 클라이언트에게 명령을 보낸다.
	broadCast(command interface{})
	// 매개변수로 받은 클라이언트와의 연결을 종료시킨다.
	remove(client *client)
}