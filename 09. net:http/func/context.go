package net

import "net/http"


// Params 필드에 라우터에서 해석한 URL 매개변수로 담고, 핸들러 내부에는 Context 값이 전달되게 한다.
type Context struct {
	Params map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request *http.Request
}

type HandlerFunc func(*Context)