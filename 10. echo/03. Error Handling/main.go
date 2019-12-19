package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

// 에러 발생을 위한 엔드포인트
func raiseError(c echo.Context) error {
	// echo.NewHTTPError() 메서드를 사용하여 *HttpError 타입의 객체를 반환하면 명시된 code와 message로 error handling을 할 수 있다.
	return echo.NewHTTPError(500, "An error has occurred.")
}

// HTTPErrorHandler는 모든 에러가 발생한 상황에서 실행되는 메서드이다.
// 대부분의 경우 기본 오류 처리로 충분하지만 에러에 따른 또 다른 조치를 취해야 하는 경우에는 유용할 수 있다.
func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	var message interface{}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message
	}

	fmt.Printf("%d | %s\n", code, message)
}

func main() {
	e := echo.New()

	// 위에서 정의한 customHTTPErrorHandler를 인스턴스 e의 HTTPErrorHandler로 등록함
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.GET("/error", raiseError)

	e.Start(":8080")
}