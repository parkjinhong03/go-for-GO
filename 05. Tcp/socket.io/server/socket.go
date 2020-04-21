package server

import (
	"github.com/googollee/go-socket.io"
	"log"
	"sync"
)

var s *socketServer

func NewSocketServer() (*socketServer, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	s = &socketServer{
		Server:  server,
		clients: []*client{},
		mutex:   &sync.Mutex{},
	}

	s.Server.OnConnect("/", func(conn socketio.Conn) (err error) {
		if err = s.accept(conn); err == nil {
			log.Printf("Completed connecting new client(id: %s)! tatal client: %d", conn.ID(), len(s.clients))
		}
		return
	})
	return s, nil
}

func (s *socketServer) accept(conn socketio.Conn) error {
	c := &client{
		conn: conn,
		id: conn.ID(),
	}
	s.clients = append(s.clients, c)
	return nil
}
