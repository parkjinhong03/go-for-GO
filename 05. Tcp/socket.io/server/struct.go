package server

import (
	"github.com/googollee/go-socket.io"
	"sync"
)

type socketServer struct {
	Server  *socketio.Server
	clients []*client
	mutex   *sync.Mutex
}

type client struct {
	conn socketio.Conn
	name string
	id   string
}