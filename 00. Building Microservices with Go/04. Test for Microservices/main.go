package main

import (
	"./data"
	"./handlers"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
)

func main() {
	serverURI := "localhost"
	if os.Getenv("DOCKER_IP") != "" {
		serverURI = os.Getenv("DOCKER_IP")
	}

	store, err := data.NewMongoStore(serverURI)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.SearchHandler{DataStore: store}
	log.Fatal(http.ListenAndServe(":8080", &handler))
}