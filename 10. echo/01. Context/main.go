package main

import (
	"github.com/labstack/echo"
)

// echo의 Context interface를 임베디드 필드로 받아 custom context 생성
type CustomContext struct {
	echo.Context
}

// 새로 만든 custom context에 foo를 출력하는 Foo() 메서드 추가
func (c *CustomContext) Foo() {
	println("foo")
}

// 새로 만든 custom context에 bar를 출력하는 Bar() 메서드 추가
func (c *CustomContext) Bar() {
	println("bar")
}

func main() {
	e := echo.New()

	// Context를 디폴트인 echo.Context 대신 Custom Context인 CustomContext로 변경 미들웨어 등록
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c}
			return next(cc)
		}
	})

	// Handler 함수에서 CustomContext를 이용해서 라우팅
	e.GET("/", func(c echo.Context) error {
		cc := c.(*CustomContext)
		cc.Foo()
		cc.Bar()
		return cc.String(200, "OK")
	})

	e.Start(":8080")
}