package main

import (
	"github.com/labstack/echo"
	"net/http"
)

// Request Body를 Bind하여 tag를 기반으로 데이터를 담을 구조체를 정의한다.
type User struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"name" query:"name"`
}

func main() {
	e := echo.New()

	e.GET("/users", func(c echo.Context) (err error) {
		u := new(User)
		// Bind(interface i)메서드를 사용하여서 Request body의 데이터를 위에서 선언한 User 구조체로 바꿔준다.
		if err = c.Bind(u); err!=nil {
			return
		}
		// 바인딩 된 정보가 담긴 인스턴스 u의 필드 값을 HTTP 프로토콜에서 전송할 수 있는 json 형식으로 변환하여 반환한다.
		return c.JSON(http.StatusOK, u)
	})

	e.Logger.Fatal(e.Start(":8080"))
}