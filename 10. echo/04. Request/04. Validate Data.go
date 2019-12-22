// echo는 프레임워크 자체로는 요청 데이터의 유효성 검사 기능을 가지고 있지 않다.
// 따라서 서드 파티 라이브러리로 custom validator를 만든 후 echo에 연동 시켜줘야 한다.

package main

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"net/http"
	"reflect"
)

type (
	ValidateUser struct {
		Name  string `json:"name" form:"name" query:"name"`
		Email string `json:"email" form:"email" query:"email"`
	}

	// validator 패키지의 Validate 타입의 멤버를 필드로 가지고 있는 CustomValidator 생성
	CustomValidator struct {
		validator *validator.Validate
	}
)

// CustomValidator를 Echo#Validator 인터페이스로 사용하기 위해 유효성을 검사하는 Validate 메서드 정의
func (cv *CustomValidator) Validate(i interface{}) error {
	// 만약 ValidateUser 구조체에서 한 필드라도 빈 값이 있으면 Bad Request를 반환한다.
	val := reflect.ValueOf(i)
	for i:=0; i<val.Elem().NumField(); i++ {
		name := val.Elem().Type().Field(i).Name
		if reflect.Indirect(val).FieldByName(name).String() == "" {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
	}
	return nil
}

func main() {
	e := echo.New()
	e.Use()
	// CustomValidater.Validate() 메서드를 사용하기 위해 인스턴스 e의 Validator를 위해서 선언한 CustomValidater로 바꿈
	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/users", func(c echo.Context) (err error) {
		u := new(ValidateUser)
		if err = c.Bind(u); err!=nil {
			return
		}
		// c.Validate() 호출시 앞에서 설정해줬던 echo.Validator 인터페이스에 정의된 객체의 Validate() 메서드를 호출한다.
		if err = c.Validate(u); err!=nil {
			return
		}
		return c.JSON(http.StatusOK, u)
	})
	e.Logger.Fatal(e.Start(":8080"))
}