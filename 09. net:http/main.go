package main

import (
	"./func"
	"fmt"
	"time"
)

type User struct {
	Id string
	AddressId string
}

func main() {
	r := net.NewServer()

	r.HandleFunc("GET", "/", func(c *net.Context) {
		c.RenderTemplate("/public/index.html", map[string]interface{}{"time":time.Now()})
	})

	r.HandleFunc("GET", "/about", func(c *net.Context) {
		fmt.Fprintln(c.ResponseWriter, "about!", c.Params)
	})

	r.HandleFunc("GET", "/users/:user_id", func(c *net.Context) {
		u := User{Id:c.Params["user_id"].(string)}
		c.RenderXml(u)
	})

	r.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *net.Context) {
		u := User{Id:c.Params["user_id"].(string), AddressId:c.Params["address_id"].(string)}
		c.RenderJson(u)
	})

	r.HandleFunc("POST", "/users", func(c *net.Context) {
		fmt.Fprintln(c.ResponseWriter, "welcome!", c.Params)
	})

	r.HandleFunc("POST", "/users/:user_id/addresses", func(c *net.Context) {
		fmt.Fprintf(c.ResponseWriter, "create user %v's address", c.Params["user_id"])
	})

	r.Run(":8080")
}