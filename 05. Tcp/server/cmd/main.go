package main

import (
	"chat.server.com/server"
	"log"
)

func main() {
	s := server.NewTcpChatServer()
	if err := s.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
	s.StartServer()
}