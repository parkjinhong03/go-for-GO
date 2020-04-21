package server

import (
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	"sort"
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
			log.Printf("Completed connecting with new client(id: %s)! tatal client: %d", conn.ID(), len(s.clients))
		}
		return
	})
	s.Server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		if err = s.remove(conn); err == nil {
			log.Printf("Completed disconnecting with client(id: %s)! tatal client: %d", conn.ID(), len(s.clients))
		}
		return
	})
	s.Server.OnEvent("/", "NAME", func(conn socketio.Conn, name string) {
		idx := s.searchIdx(conn.ID())
		s.clients[idx].name = name
		for _, i := range s.clients { fmt.Println(i.name) }
	})
	return s, nil
}

func (s *socketServer) searchIdx(id string) (idx int) {
	idx = sort.Search(len(s.clients)-1, func(i int) bool {
		return s.clients[i].id == id
	})
	return
}

func (s *socketServer) accept(conn socketio.Conn) error {
	c := &client{
		conn: conn,
		id: conn.ID(),
		name: "undefined",
	}
	s.clients = append(s.clients, c)
	return nil
}

func (s *socketServer) remove(conn socketio.Conn) error {
	idx := s.searchIdx(conn.ID())
	s.clients = append(s.clients[:idx], s.clients[idx+1:]...)
	return nil
}