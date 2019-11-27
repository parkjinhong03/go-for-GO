package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var renderer *render.Render

func init() {
	// 렌더러 생성
	renderer = render.New()
}

func main() {
	// 라우터 생성
	router := httprouter.New()

	// 핸들러 정의
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// 렌더러를 사용하여 템플릿 렌더링
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title":"Simple Chat!"})
	})

	// negroni 미들웨어 생성
	n := negroni.Classic()

	// negroni에 router를 핸들러로 등록
	n.UseHandler(router)

	// 웹 서버 실행
	n.Run(":8000")
}
