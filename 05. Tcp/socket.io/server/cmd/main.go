package main

import (
	"log"
	"net/http"
	"socket.io/test/server"
)

func main() {
	s, err := server.NewSocketServer()
	if err != nil {
		log.Fatal(err)
	}

	go s.Server.Serve()
	defer s.Server.Close()

	http.Handle("/socket.io/", s.Server)
	log.Printf("Socket server starting on port %v", ":8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}