package main

import (
	"MSA.example.com/1/tool/message"
	"fmt"
	"log"
	"net/http"
)

func main () {
	natsM, err := message.GetDefaultNatsByEnv()
	if err != nil {
		log.Fatalf("unable to connect nats message server, err:%v\n", err)
	}
	handler := &apiGateWayHandler{natsM: natsM}

	http.Handle("/api/", http.StripPrefix("/api/", handler))
	log.Println("Server is starting on port 8080...")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

type apiGateWayHandler struct {
	natsM message.NatsMessage
}

func (h *apiGateWayHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	_, _ = fmt.Fprint(rw, 1)
	return
}