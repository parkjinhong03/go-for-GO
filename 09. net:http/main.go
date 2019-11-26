package main

import (
	"./func"
	"fmt"
)

func main() {
	r := net.NewServer()

	r.HandleFunc("GET", "/", func(c *net.Context) {
		fmt.Fprintln(c.ResponseWriter, "welcome!", c.Params)
	})

	r.HandleFunc("GET", "/about", func(c *net.Context) {
		fmt.Fprintln(c.ResponseWriter, "about!", c.Params)
	})

	r.HandleFunc("GET", "/users/:user_id", func(c *net.Context) {
		if c.Params["user_id"] == "0" {
			panic("id is zero")
		}
		fmt.Fprintf(c.ResponseWriter, "recieve user %v\n", c.Params["user_id"])
	})

	r.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *net.Context) {
		fmt.Fprintf(c.ResponseWriter, "retrieve user %v's address %v\n", c.Params["user_id"], c.Params["address_id"])
	})

	r.HandleFunc("POST", "/users", func(c *net.Context) {
		fmt.Fprintln(c.ResponseWriter, "welcome!", c.Params)
	})

	r.HandleFunc("POST", "/users/:user_id/addresses", func(c *net.Context) {
		fmt.Fprintf(c.ResponseWriter, "create user %v's address", c.Params["user_id"])
	})

	r.Run(":8080")
}