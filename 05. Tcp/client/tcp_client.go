package client

import (
	"chat.server.com/protocol"
	"net"
)

type TcpChatClient struct {
	conn *net.Conn							// 서버와의 연결을 저장하기 위한 필드
	writer protocol.CommandWriter			// 서버에게 명령을 전송하기 위한 writer 필드
	reader protocol.CommandReader			// 서버에게 명령을 수신하기 위한 reader 필드
	name string								// 클라이언트의 이름을 저장하기 위한 필드
	incoming chan protocol.MessageCommand	// 서버가 명령을 전송하였다는 것을 알려주기 위한 채널 필드
}

func NewTcpChatClient() *TcpChatClient {
	return &TcpChatClient{
		incoming: make(chan protocol.MessageCommand),
	}
}