package handlers

import (
	"../entities"
	"../middlewares"
	"encoding/json"
	"net/http"
)

type helloWorldHandler struct {}

func NewHelloWorldHandler() *helloWorldHandler {
	return &helloWorldHandler{}
}

func (h *helloWorldHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	name := r.Context().Value(middlewares.ValidationContextKey("name")).(string)

	response := entities.HelloWorldResponse{Message: "hello " + name}
	encoder := json.NewEncoder(rw)
	encoder.Encode(response)
}