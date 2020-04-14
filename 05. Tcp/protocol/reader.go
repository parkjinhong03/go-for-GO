// 스트림으로 부터 명령 타입 및 구문을 분석하기 위한 reader 이다.

package protocol

import (
	"bufio"
	"io"
	"strings"
)

type CommandReader struct {
	reader *bufio.Reader
}

func NewCommandReader(r io.Reader) *CommandReader {
	return &CommandReader{
		reader: bufio.NewReader(r),
	}
}

func (r *CommandReader) Read() (v interface{}, err error) {
	// reader.ReadString 메서드를 호출하면 reader에 저장되어 있는 문자열의 첫 번째 ' '까지 삭제한 후 반환한다.
	command, err := r.reader.ReadString(' ')
	if err != nil {
		return
	}

	// 위에서 파싱한 명령어 이름에 따라 각각의 처리를 하는 switch case 구문
	switch command {
	case "SEND ":
		var message string
		message, err = r.reader.ReadString('\n')
		if err != nil { v = nil; return }
		v = SendCommand{Message: strings.Split(message, "\n")[0]}

	case "NAME ":
		var name string
		name, err = r.reader.ReadString('\n')
		if err != nil { v = nil; return }
		v = NameCommand{Name: strings.Split(name, "\n")[0]}

	case "MESSAGE ":
		var name, message string
		name, err = r.reader.ReadString(' ')
		if err != nil { v = nil; return }
		message, err = r.reader.ReadString('\n')
		if err != nil { v = nil; return }
		v = MessageCommand{
			Name:    strings.Split(name, " ")[0],
			Message: strings.Split(message, "\n")[0],
		}

	default:
		_, _ = r.reader.ReadString('\n')
		v = UndefinedCommand
	}

	return
}