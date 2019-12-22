// 이번에는 Go의 구조체를 HTTP 프로토콜 위에서 통신될 수 있는 데이터(json)으로 인코딩 한 후
// 그 인코딩한 값과 상태 코드를 사용자에게 응답하는 여러가지 방법을 알아볼 것 이다.
package main

import (
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
)

type User struct {
	Name string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

// 1. Context#JSON(code int, i interface{})
// 이 메서드를 이용하여 Go 데이터 자료형을 JSON으로 encode하고 status code와 같이 반환할 수 있다.
func JsonHandler(c echo.Context) error {
	u := &User{
		Name:  "myName",
		Email: "myEmail@naver.com",
	}
	return c.JSON(http.StatusOK ,u)
}

// 2. built-in json 패키지의 사용
// Context#JSON는 내부적으로 json.Marshal을 사용하기 때문에 구조체 데이터가 커짐에 따라 효율이 떨어진다.
// 따라서 그러한 경우, 직접 json 패키지를 이용하여 구조체를 인코딩하여 효율을 높일 수 있다.
func StreamJsonHandler(c echo.Context) error {
	u := &User{
		Name:  "myName",
		Email: "myEmail@naver.com",
	}

	// Response body가 JSON임을 Response Header의 Content-Type에 직접 명시
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	// Response start-line에 status code가 200임을 직접 명시
	c.Response().WriteHeader(http.StatusOK)
	// c.response를 Writer로 가지고있는 Encoder 객체 생성 후 u 인스턴스 json으로 인코딩 진행
	return json.NewEncoder(c.Response()).Encode(u)
}

// 3. Context#JSONPretty(code int, i interface{}, indent string)
// 이 메서드는 indent를 이용하여 가독성이 높은 json을 만들 때 사용한다.
// 참고로 Context#JSON과 같이 indent를 따로 지정하지 않았다면 기본적으로 "  "가 적용된다.
// 또한, 요청시 URL의 query에 indent 값을 준다면 Context#JSON 메서드로도 pretty json이 생성된다.
func PrettyJsonHandler(c echo.Context) error {
	u := &User{
		Name:  "myName",
		Email: "myEmail@naver.com",
	}

	return c.JSONPretty(http.StatusOK, u, "    ")
	/*
    {
	    "Name": "myName",
	    "Email": "myEmail@naver.com"
    }
	*/
}


func main() {
	e := echo.New()

	e.GET("/json", JsonHandler)
	e.GET("/jsonStream", StreamJsonHandler)
	e.GET("/prettyJson", PrettyJsonHandler)

	e.Logger.Fatal(e.Start(":8080"))
}