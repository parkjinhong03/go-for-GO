package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

// 새 쿠키 생성 후 클라이언트에 쿠키를 저장시키는 엔드포인트
func writeCookie(c echo.Context) error {
	cookie := new(http.Cookie) // 새 쿠키는 new(http.Cookie)로 생성된다.
	cookie.Name = "username"
	cookie.Value = "jon"
	cookie.Expires = time.Now().Add(24 * time.Hour) // Cookie의 인스턴스 객체의 public attribute에 접근하여 변경할 수 있다.
	c.SetCookie(cookie) // SetCookie() 메서드를 사용하여 Cookie 정보를 HTTP Response에 포함할 수 있다.
	return c.String(http.StatusOK, "write a cookie")
}

// 클라이언트로부터 쿠키를 받아와 인증을 진행하는 엔드포인트
func readCookie(c echo.Context) error {
	cookie, err := c.Cookie("username") // c.Cookie(string) 메서드를 이용하여 HTTP Request의 쿠키 받아온다.
	if err != nil {
		return err
	}
	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)
	return c.String(http.StatusOK, "read a cookie!")
}

func main() {
	e := echo.New()

	e.GET("/writeCookie", writeCookie)
	e.GET("/readCookie", readCookie)

	e.Start(":8080")
}