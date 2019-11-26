package main

import (
	"./router"
	"fmt"
	"net/http"
)

func main() {
	r := router.Router{Handlers: make(map[string]map[string]http.HandlerFunc)}

	r.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "welcome!")
	})

	r.HandleFunc("GET", "/about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "about")
	})

	r.HandleFunc("GET", "/users/:user_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "retrieve user")
	})

	r.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "retrieve user's address")
	})

	r.HandleFunc("POST", "/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "create user")
	})

	r.HandleFunc("POST", "/users/:user_id/addresses", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "create user's address")
	})

	http.ListenAndServe(":8080", r)
}