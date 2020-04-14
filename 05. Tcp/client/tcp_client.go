package client

import (
	"chat.server.com/protocol"
	"io"
	"log"
	"net"
)

type TcpChatClient struct {
	conn net.Conn							// 서버와의 연결을 저장하기 위한 필드
	writer *protocol.CommandWriter			// 서버에게 명령을 전송하기 위한 writer 필드
	reader *protocol.CommandReader			// 서버에게 명령을 수신하기 위한 reader 필드
	name string								// 클라이언트의 이름을 저장하기 위한 필드
	incoming chan protocol.MessageCommand	// 서버가 명령을 전송하였다는 것을 알려주기 위한 채널 필드
}

func NewTcpChatClient() *TcpChatClient {
	return &TcpChatClient{
		incoming: make(chan protocol.MessageCommand),
	}
}

func (tc *TcpChatClient) Dial(address string) (err error) {
	// net.Dial 함수를 이용하여 서버의 연결을 시도할 수 있다.
	conn, err := net.Dial("tcp", address)
	if err != nil { return }

	tc.conn = conn
	tc.writer = protocol.NewCommandWriter(conn)
	tc.reader = protocol.NewCommandReader(conn)
	return
}

func (tc *TcpChatClient) Start() {
	for {
		// reader를 이용해서 서버로부터 받은 명령을 읽어들일 수 있다.
		v, err := tc.reader.Read()
		if err == io.EOF {
			continue
		} else if err != nil {
			log.Printf("Some error occurs while reading command from server err: %v\n", err)
			return
		}

		switch cmd := v.(type) {
		// 명령어의 타입이 MessageCommand일 경우, incoming 채널에 해당 명령어를 송신한다.
		case protocol.MessageCommand:
			tc.incoming <- cmd
		default:
			log.Println("Undefined command comes from server")
		}
	}
}

func (tc *TcpChatClient) Close() {
	_ = tc.conn.Close()
}

func (tc *TcpChatClient) Send(v interface{}) error {
	// writer를 이용하여 서버에게 명령을 전송할 수 있다.
	return tc.writer.Write(v)
}

func (tc *TcpChatClient) SendMessage(message string) error {
	return tc.Send(protocol.SendCommand{Message: message})
}

func (tc *TcpChatClient) SetName(name string) error {
	return tc.Send(protocol.NameCommand{Name: name})
}

func (tc *TcpChatClient) Disconnect() (err error) {
	err = tc.Send(protocol.DisconnectCommand{})
	_ = tc.conn.Close()
	return
}

func (tc *TcpChatClient) Incoming() <-chan protocol.MessageCommand {
	return tc.incoming
}