package main

import (
	"MSA.example.com/1/entities"
	"MSA.example.com/1/middleware"
	"MSA.example.com/1/protocol"
	"MSA.example.com/1/proxy"
	natsEncoder "MSA.example.com/1/tool/encoder/nats"
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
		log.Fatalf("unable to connect nats message server, err: %v\n", err)
	}
	validate := validator.New()

	authH := middleware.NewCorrelationMiddleware(&authServiceHandler{
		natsM: natsM,
		validate: validate,
		natsE: natsEncoder.NewJsonEncoder(proxy.NewAuthServiceProxy(natsM)),
	})

	http.Handle("/api/auth/", http.StripPrefix("/api/auth/", authH))
	http.Handle("/api/auth", http.StripPrefix("/api/auth", authH))
	log.Println("Server is starting on port 8080...")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

type authServiceHandler struct {
	natsM message.NatsMessage
	validate *validator.Validate
	natsE natsEncoder.Encoder
}

func (h *authServiceHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Println(r.URL.Path)
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	d := json.NewDecoder(r.Body)
	switch r.URL.Path {
	case "signup", "signup/":
		request := entities.AuthSignUpEntities{}
		err := d.Decode(&request)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.validate.Struct(&request)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		err = h.natsE.Encode(protocol.AuthSignUpProtocol{
			RequestId:     r.Header.Get("X-Request-ID"),
			UserId:        request.UserId,
			UserPwd:       request.UserPwd,
			Name:          request.Name,
			PhoneNumber:   request.PhoneNumber,
			Introduction:  request.Introduction,
			Email:         request.Email,
			ReturnChannel: "auth.signup.return",
			InputChannel:  "auth.signup",
		})
		if err != nil {
			log.Printf("some error occurs while encoding to message, err: %v\n", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		return
	}

	rw.WriteHeader(http.StatusNotFound)
	return
}