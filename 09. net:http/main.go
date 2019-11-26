package main

import (
	"./func"
	"fmt"
	"net/http"
)

func main() {
	r := net.Router{Handlers: make(map[string]map[string]net.HandlerFunc)}

	r.HandleFunc("GET", "/", net.LogHandler(func(c *net.Context) {
		fmt.Fprintln(c.ResponseWriter, "welcome!")
	}))

	r.HandleFunc("GET", "/about", func(c *net.Context) {
		fmt.Fprintln(c.ResponseWriter, "about")
	})

	r.HandleFunc("GET", "/users/:user_id", func(c *net.Context) {
		fmt.Fprintf(c.ResponseWriter, "retrieve user %v\n", c.Params["user_id"])
	})

	r.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *net.Context) {
		fmt.Fprintf(c.ResponseWriter, "retrieve user %v's address %v\n", c.Params["user_id"], c.Params["address_id"])
	})

	r.HandleFunc("POST", "/users", func(c *net.Context) {
		fmt.Fprintln(c.ResponseWriter, "create user")
	})

	r.HandleFunc("POST", "/users/:user_id/addresses", func(c *net.Context) {
		fmt.Fprintf(c.ResponseWriter, "create user %v's address", c.Params["user_id"])
	})

	http.ListenAndServe(":8080", r)
}