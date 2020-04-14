package server

import (
	"chat.server.com/protocol"
	"io"
	"log"
	"net"
	"sync"
)

type TcpChatServer struct {
	listener net.Listener	// 연결을 수신하기 위해 서버를 대기시키 위한 필드
	clients []*client		// 연결되어있는 모든 클라이언트에 대한 정보를 담기 위한 필드
	mutex *sync.Mutex		// race condition을 다루기 위한 mutex 필드
}

type client struct {
	conn net.Conn						// 클라이언트와 연결에 대한 정보를 담아두기 위한 필드
	name string							// 클라이언트의 사용 이름을 저장시키기 위한 필드
	writer *protocol.CommandWriter		// 클라이언트 측에 명령을 보내기 위한 필드
}

func NewTcpChatServer() *TcpChatServer {
	return &TcpChatServer{
		mutex: &sync.Mutex{},
	}
}

func (ts *TcpChatServer) Listen(address string) (err error) {
	var l net.Listener
	// net.Listen 함수를 이용하여 해당 주소에 서버를 대기시키기 위한 객체를 얻을 수 있다.
	if l, err = net.Listen("tcp", address); err != nil {
		return
	}

	log.Printf("Listening on %v\n", address)
	ts.listener = l
	return
}

func (ts *TcpChatServer) StartServer() {
	for {
		// 클라이언트로 부터 새로운 연결 요청이 있을때까지 대기한다.
		conn, err := ts.listener.Accept()
		if err != nil {
			log.Printf("Unable to Accept err: %v\n", err)
			return
		}

		// 요청이 들어오면 해당 클라이언트와의 연결을 저장하고, 고루틴으로 비즈니스 로직(명령어)을 처리한다.
		client := ts.accept(conn)
		go ts.serve(client)
	}
}

func (ts *TcpChatServer) CloseServer() {
	_ = ts.listener.Close()
}

func (ts *TcpChatServer) accept(conn net.Conn) *client {
	log.Printf("Accepting new connection from %v... (current clients: %v)\n", conn.RemoteAddr().String(), len(ts.clients))

	// race condition을 막기 위한 mutex Lock & Unlock
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	client := &client{
		conn:   conn,
		writer: protocol.NewCommandWriter(conn),
	}

	ts.clients = append(ts.clients, client)
	log.Printf("Complete accepting new connection from %v! (current clients: %v)\n", conn.RemoteAddr().String(), len(ts.clients))
	return client
}

func (ts *TcpChatServer) serve(client *client) {
	reader := protocol.NewCommandReader(client.conn)
	defer ts.remove(client)

	for {
		// 위에서 선언한 Reader를 이용해 클라이언트가 보낸 명령어를 읽어들인다.
		v, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Unable to parse string command, err: %v\n", err)
			return
		}
		switch cmd := v.(type) {
		case protocol.NameCommand:
			client.name = cmd.Name
		case protocol.SendCommand:
			go ts.broadCast(protocol.MessageCommand{
				Name:    client.name,
				Message: cmd.Message,
			})
		}
	}
}

func (ts *TcpChatServer) remove(client *client) {
	// race condition을 막기 위한 mutex Lock & Unlock
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	for i, c := range ts.clients {
		if c != client { continue }
		// 해당 클라이언트를 현재 연결중인 클라이언트들의 목록에서 제거한다.
		ts.clients = append(ts.clients[:i], ts.clients[i+1:]...)
		_ = client.conn.Close()
		log.Printf("Closed connection from %v. (current clients: %v)\n", client.conn.RemoteAddr().String(), len(ts.clients))
		return
	}

	log.Println("something error... (on TcpChatServer.remove)")
}

func (ts *TcpChatServer) broadCast(v interface{}) {
	// 서버의 목록에 있는 모든 클라이언트들에게 명령어를 보낸다.
	for _, client := range ts.clients {
		_ = client.writer.Write(v)
	}
}