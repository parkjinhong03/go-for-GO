package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type CustomUser struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"name" query:"name"`
}

// 커스텀 바인더를 생성하기 위해 빈 구조체 정의
type CustomBinder struct {}

// 위에서 만들 구조체에 Bind(interface{}, echo.Context) err 서명의 메서드를 추가하여 echo.Binder 인터페이스로 사용할 수 있게 함
func (cb CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
	// Custom Binder를 위한 echo 프레임워크의 Default Binder 생성
	db := new(echo.DefaultBinder)

	// Default Binder의 Bind() 메서드를 호출한 후, 기본적인 Bind로 디코딩할 수 있는 Content-Type 이면 CustomBinder Bind() 종료
	// 참고로 가능한 Content-Type 으로는 application/json, application/xml, application/x-www-form-urlencoded data 가 있다.
	if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
		fmt.Println("Default Binder로 디코딩할 수 있는 Content-Type의 요청입니다.")
		return
	}

	fmt.Println("Default Binder로 디코딩할 수 없는 Content-Type의 요청입니다.")
	return nil
}

func main() {
	e := echo.New()
	// 인스턴스 e의 Binder를 DefaultBinder에서 위에서 정의한 CustomBinder로 교체
	e.Binder = &CustomBinder{}

	e.GET("/users", func(c echo.Context) error {
		u := new(CustomUser)

		// c.Bind(u) 호출시 내부에서 CustomBind.Bind(u, c) 메서드를 호출함
		if err := c.Bind(u); err!=nil {
			return err
		}
		return c.JSON(http.StatusOK, u)
	})

	e.Logger.Fatal(e.Start(":8080"))
}