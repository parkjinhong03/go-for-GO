package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type helloWorldRequest4 struct {
	Name string `json:"name"`
}

type helloWorldResponse7 struct {
	Message string `json:"message"`
}

type validationHandler1 struct {
	next http.Handler
}

func newValidationHandler1(next http.Handler) http.Handler {
	return validationHandler1{next: next}
}

type validationContextKey string

func (h validationHandler1) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request helloWorldRequest4

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)

	if err != nil {
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(r.Context(), validationContextKey("name"), request.Name)
	r = r.WithContext(ctx)

	h.next.ServeHTTP(rw, r)
}

type helloWorldHandler6 struct {}

func newHelloWorldHandler1() http.Handler {
	return helloWorldHandler6{}
}

func (h helloWorldHandler6) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	name := r.Context().Value(validationContextKey("name")).(string)
	response := helloWorldResponse7{Message: "Hello " + name}

	encoder := json.NewEncoder(rw)
	encoder.Encode(response)
}

func main() {
	port := 8080

	handler := newValidationHandler1(newHelloWorldHandler1())

	http.Handle("/helloworld", handler)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}