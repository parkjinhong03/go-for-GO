package protocol_test

import (
	"bytes"
	"chat.server.com/protocol"
	"testing"
)

func TestWriteCommand(t *testing.T) {
	commands := []struct {
		command interface{}
		result 	string
	}{
		{
			command: protocol.SendCommand{Message: "Test"},
			result:  "SEND Test\n",
		},
		{
			command: protocol.NameCommand{Name: "Test"},
			result:  "NAME Test\n",
		},
	}

	// Writer의 write 메서드 호출 후 결과값을 저장하기 위해 bytes.Buffer 타입의 객체를 하나 선언한다.
	buf := new(bytes.Buffer)
	cmdWriter := protocol.NewCommandWriter(buf)

	for _, c := range commands {
		// 버퍼에 저장된 값을 초기화하기 위해 Reset 메서드를 호출해야 한다.
		buf.Reset()
		if err := cmdWriter.Write(c.command); err != nil {
			t.Errorf("Unable to write command %v", err)
		}
		// String 메서드를 이용하여 위에서 Write 메서드를 호출하여 저장한 문자열 값을 얻을 수 있다.
		if buf.String() != c.result {
			t.Errorf("Command output is not same %v %v", c.result, buf.String())
		}
	}
}