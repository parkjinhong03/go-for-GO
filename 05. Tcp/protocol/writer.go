// 명령을 문자열로 변환시키기 위한 writer 이다.

package protocol

import (
	"fmt"
	"io"
)

// io.Writer 인터페이스의 객체를 필드에 포함시켜 작성 값을 저장시키는 용도로 사용한다.
type CommandWriter struct {
	writer io.Writer
}

func NewCommandWriter(w io.Writer) *CommandWriter {
	return &CommandWriter{
		writer: w,
	}
}

func (w *CommandWriter) writeString(msg string) (err error) {
	// io.Writer.Write 메서드를 이용하여 매개변수로 받은 바이트 배열(문자열을 변환시킨)을 저장시킨다.
	_, err = w.writer.Write([]byte(msg))
	return
}

func (w *CommandWriter) Write(command interface{}) (err error) {
	// switch type 구문을 이용하여 매개변수로 받은 command의 타입을 검사하고 그에 따른 처리를 진행한다.
	switch v := command.(type) {
	case SendCommand:
		err = w.writeString(fmt.Sprintf("SEND %s\n", v.Message))
	case NameCommand:
		err = w.writeString(fmt.Sprintf("NAME %s\n", v.Name))
	case MessageCommand:
		err = w.writeString(fmt.Sprintf("MESSAGE %s %s\n", v.Name, v.Message))
	default:
		err = UndefinedCommand
	}

	return
}