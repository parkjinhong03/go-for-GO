package main

import (
	"MSA.example.com/1/entities"
	"MSA.example.com/1/middleware"
	"MSA.example.com/1/tool/message"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

func main () {
	natsM, err := message.GetDefaultNatsByEnv()
	if err != nil {
		log.Fatalf("unable to connect nats message server, err:%v\n", err)
	}

	authH := middleware.NewCorrelationMiddleware(&authServiceHandler{
		natsM: natsM,
		validate: validator.New(),
	})

	http.Handle("/api/auth/", http.StripPrefix("/api/auth/", authH))
	http.Handle("/api/auth", http.StripPrefix("/api/auth", authH))
	log.Println("Server is starting on port 8080...")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

type authServiceHandler struct {
	natsM message.NatsMessage
	validate *validator.Validate
}

func (h *authServiceHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	
}