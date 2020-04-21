package server

import (
	"github.com/googollee/go-socket.io"
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
	return s, nil
}
