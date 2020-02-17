package handlers

import (
	"../entities"
	"encoding/json"
	"net/http"
)

type helloWorldHandler struct {}

func NewHelloWorldHandler() *helloWorldHandler {
	return &helloWorldHandler{}
}

func (h *helloWorldHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	response := entities.HelloWorldResponse{Message: "hello "}
	encoder := json.NewEncoder(rw)
	encoder.Encode(response)
}