// 이번에는 Go의 구조체를 XML로 인코딩하여 그 값과 상태 코드를 함께 응답하는 여러 방법들을 알아볼 것 이다.
// 저번의 Json encoding과 방법이 매우 비슷하므로 쉽게 이해할 수 있을 것 이다.
package main

import (
	"encoding/xml"
	"github.com/labstack/echo"
	"net/http"
)

type XmlUser struct {
	Name string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

// 1. Content#XML(code int, i interface{})
// 이 메서드는 Go의 구조체를 XML로 Encoding하여 status code와 함께 반환해준다.
func XmlHandler(c echo.Context) error {
	u := &XmlUser{
		Name:  "myName",
		Email: "myEmail@naver.com",
	}

	return c.XML(http.StatusOK, u)
}

// 2. built-in xml 패키지의 사용
// Context#XML 또한 내부적으로 xml.Marshal을 사용하기 때문에 구조체 데이터가 커짐에 따라 효율이 떨어진다.
// 따라서 그러한 경우, 직접 xml 패키지를 이용하여 구조체를 인코딩하여 효율을 높일 수 있다.
func StreamXmlHandler(c echo.Context) error {
	u := &XmlUser{
		Name:  "myName",
		Email: "myEmail@naver.com",
	}

	// Response body가 XML임을 Response Header의 Content-Type에 직접 명시
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)
	// Response start-line에 status code가 200임을 직접 명시
	c.Response().WriteHeader(http.StatusOK)
	// c.response를 Writer로 가지고있는 Encoder 객체 생성 후 u 인스턴스 xml로 인코딩 진행
	return xml.NewEncoder(c.Response()).Encode(u)
}

// 3. Context#XMLPretty(code int, i interface{}, indent string)
// 이 메서드는 indent를 이용하여 가독성이 높은 xml을 만들 때 사용한다.
// 참고로 Context#XML 같이 indent를 따로 지정하지 않았다면 기본적으로 "  "가 적용된다.
// 또한, 요청시 URL의 query에 indent 값을 준다면 Context#XML 메서드로도 pretty xml이 생성된다.
func PrettyXmlHandler(c echo.Context) error {
	u := &XmlUser{
		Name:  "myName",
		Email: "myEmail@naver.com",
	}

	return c.XMLPretty(http.StatusOK, u, "    ")
	/*
	<XmlUser>
	    <name>myName</name>
	    <email>myEmail@naver.com</email>
	</XmlUser>
	*/
}

func main() {
	e := echo.New()

	e.GET("/xml", XmlHandler)
	e.GET("/xmlStream", StreamXmlHandler)
	e.GET("/prettyXml", PrettyXmlHandler)

	e.Logger.Fatal(e.Start(":8080"))
}