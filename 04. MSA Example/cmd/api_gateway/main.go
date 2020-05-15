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
			_, _ = fmt.Fprintf(rw, "Please format body as json, err:%v", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.validate.Struct(&request)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// json 마샬링 후 이벤트 발생 및 수신 코드 추가

		rw.WriteHeader(http.StatusOK)
		return
	}

	_, _ = fmt.Fprint(rw, "404 page not found")
	rw.WriteHeader(http.StatusNotFound)
	return
}