package client

import "chat.server.com/protocol"

// 클라이언트의 동작을 명확하게 정의하기 위해 정의한 클라이언트 인터페이스
type ChatClient interface {
	// 서버와의 연결을 생성해주는 메서드
	Dial(address string) error
	// 서버에게 명령어를 전송할 때 사용하는 메서드
	Send(command interface{}) error
	// 서버에게 메세지 송신 명령어를 전송할 때 사용하는 메서드
	SendMessage(message string) error
	// 서버에게 이름 설정 명령어를 전송할 때 사용하는 메서드
	SetName(name string) error
	// Start&Close 메서드를 이용하여 클라이언트를 시작&중지시킬 수 있다.
	Start()
	Close()
	// 서버로부터 명령어를 수신하기 위한 채널을 반환해주는 메서드
	Incoming() chan protocol.MessageCommand
}