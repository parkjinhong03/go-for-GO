package protocol

import "errors"

var UndefinedCommand = errors.New("undefined command")

// SendCommand는 새 메세지를 클라이언트에게 전송할 때 사용된다.
type SendCommand struct {
	Message string
}

// NameCommand는 클라이언트의 사용 이름을 설정할 때 사용된다.
type NameCommand struct {
	Name string
}

// MessageCommand는 새 메세지를 다른 사용자에게 알릴 때 사용된다.
type MessageCommand struct {
	Name string
	Message string
}