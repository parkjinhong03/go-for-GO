package main

import "C"
import (
	"./func"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
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

	r.HandleFunc("GET", "/users/:user_id", func(c *net.Context) {
		u := User{Id:c.Params["user_id"].(string)}
		c.RenderXml(u)
	})

	r.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *net.Context) {
		u := User{Id:c.Params["user_id"].(string), AddressId:c.Params["address_id"].(string)}
		c.RenderJson(u)
	})

	r.HandleFunc("GET", "/login", func(c *net.Context) {
		c.RenderTemplate("/public/login.html", map[string]interface{}{
			"message": "로그인이 필요합니다.",
		})
	})

	r.HandleFunc("POST", "/login", func(c *net.Context) {
		if CheckLogin(c.Params["username"].(string), c.Params["password"].(string)) {
			http.SetCookie(c.ResponseWriter, &http.Cookie{
				Name: "X_AUTH",
				Value: Sign(VerifyMessage),
				Path:"/",
			})
			c.Redirect("/")
		}
		c.RenderTemplate("/public/login.html", map[string]interface{}{
			"message": "id 또는 password가 일치하지 않습니다.",
		})
	})

	r.Use(AuthHandler)

	r.Run(":8080")
}

func CheckLogin(username, password string) bool {
	const (
		USERNAME = "tester"
		PASSWORD = "12345"
	)

	return username == USERNAME && password == PASSWORD
}

func Sign(message string) string {
	secretKey := []byte("golang-book-secret-key2")
	if len(secretKey) == 0 {
		return ""
	}

	mac := hmac.New(sha1.New, secretKey)
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}