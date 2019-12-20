// 이번 화는 Form Data, Query Parameters, Path Parameters 데이터들을 검색하는 방법을 알아본다.
package main

import (
	"github.com/labstack/echo"
	"net/http"
)

// Form Data 에서 데이터를 검색하는 엔드포인트
// curl -X POST  http://localhost:1323 -d 'name=Joe'
func FormHandler (c echo.Context) error {
	name := c.FormValue("name")
	return c.String(http.StatusOK, name)
}

// Query Parameter 에서 데이터를 검색하는 엔드포인트
// curl -X GET http://localhost:1323\?name\=Joe
func queryHandler(c echo.Context) error {
	name := c.QueryParam("name")
	return c.String(http.StatusOK, name)
}

// Path Parameters 에서 데이터를 검색하는 엔드포인트
// curl http://localhost:1323/users/Joe
func pathHandler(c echo.Context) error {
	name := c.Param("name")
	return c.String(http.StatusOK, name)
}

func main() {
	e := echo.New()

	e.GET("/form", FormHandler)
	e.GET("/query", queryHandler)
	// echo에서는 동적 URI를 표시할 때 :를 사용한다.
	// 사용자의 요청이 들어왔을 때 수행되는 ServeHTTP 메서드에서 :name 부분의 값을 파싱해서 echo.Context의 pnames와 pvalues에 저장한다.
	e.GET("/path/:name", pathHandler)

	e.Logger.Fatal(e.Start(":8080"))
}