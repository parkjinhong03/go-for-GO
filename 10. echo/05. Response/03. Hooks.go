// 이번 화에서는 요청이 들어왔을 때, 응답이 완료되기 직전에 실행되는 함수와 응답이 완료된 바로 실행되는 함수를 만드는 방법에 대해 알아본다.
package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func anyHandler(c echo.Context) error {
	// Content#Response#Before(func())
	// Before 메서드의 인자값으로 전달된 클로저 함수는, 이 핸들러로 요청이 들어온 후 응답이 완료되기 바로 전에 실행된다.
	c.Response().Before(func() {
		fmt.Println("before response!!")
	})
	// Content#Response#After(func())
	// After 메서드의 인자값으로 전달된 클로저 함수는, 이 핸들러로 요청이 들어온 후 응답이 완료되고 바로 다음에 실행된다.
	c.Response().After(func() {
		fmt.Println("after response!!")
	})

	// 아래 함수를 기준으로 실행되지 전과 후, 각각 위의 클로저 함수들이 실행된다.
	// 만약 요청 응답을 c.NoContent(http.StatusOK)로 수행할 경우에는 Content#Response#After(func())로 등록한 함수가 실행되지 않는다.
	return c.String(http.StatusOK, "hook test!")
}

func main() {
	e := echo.New()

	e.GET("/", anyHandler)
	e.Logger.Fatal(e.Start(":8080"))
}