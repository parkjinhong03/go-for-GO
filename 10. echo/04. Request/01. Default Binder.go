// Binder란, Request의 Content-Type에 따라 다른 Request Body를 Go에서 사용할 수 있는 타입으로 바꿔준다.
// 뒤에서 다뤄볼 Custom Binder를 선언하지 않고 사용하는 Binder를 Default Binder라고 하는데, 이 Default Binder에 대해 알아보겠다.
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
		// 여기서 호출한 Bind() 메서드는 따로 Binder를 지정하지 않았으므로 DefaultBinder의 Bind() 메서드가 호출될 것 이다.
		if err = c.Bind(u); err!=nil {
			return
		}
		// 바인딩 된 정보가 담긴 인스턴스 u의 필드 값을 HTTP 프로토콜에서 전송할 수 있는 json 형식으로 변환하여 반환한다.
		return c.JSON(http.StatusOK, u)
	})

	e.Logger.Fatal(e.Start(":8080"))
}