package main

import (
	"./func"
	"crypto/hmac"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Id string
	AddressId string
}

const VerifyMessage = "verified"

func AuthHandler(next net.HandlerFunc) net.HandlerFunc {
	ignore := []string{"/login", "public/index.html"}
	return func(c *net.Context) {
		for _, s := range ignore {
			if strings.HasPrefix(c.Request.URL.Path, s) {
				next(c)
				return
			}
		}
		if v, err := c.Request.Cookie("X_AUTH"); err == http.ErrNoCookie {
			c.Redirect("/login")
			return
		} else if err != nil {
			c.RenderErr(http.StatusInternalServerError, err)
			return
		} else if Verify(VerifyMessage, v.Value) {
			next(c)
			return
		}

		c.Redirect("/login")
	}
}

func Verify(message, sig string) bool {
	return hmac.Equal([]byte(sig), []byte(Sign(message)))
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