package main

import (
	"MSA.example.com/1/middleware"
	natsEncoder "MSA.example.com/1/tool/encoder/nats"
	"MSA.example.com/1/tool/message"
	"MSA.example.com/1/tool/proxy"
	"MSA.example.com/1/usecase/apiGatewayUsecase"
	"github.com/eapache/go-resiliency/breaker"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"time"
)

func main () {
	natsM, err := message.GetDefaultNatsByEnv()
	if err != nil {
		log.Fatalf("unable to connect nats message server, err: %v\n", err)
	}
	validate := validator.New()
	natsE := natsEncoder.NewJsonEncoder(proxy.NewAuthServiceProxy(natsM, validate))
	breakeR := breaker.New(3, 1, time.Minute)

	authH := middleware.NewCorrelationMiddleware(
		apiGatewayUsecase.NewAuthServiceHandler(natsM, validate, natsE, breakeR),
	)

	http.Handle("/api/auth/", http.StripPrefix("/api/auth/", authH))
	http.Handle("/api/auth", http.StripPrefix("/api/auth", authH))
	log.Println("Server is starting on port 8080...")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
