package net

import "net/http"

type Server struct {
	*Router
	middlewares []Middleware
	startHandler HandlerFunc
}

func NewServer() *Server {
	r := &Router{Handlers: map[string]map[string]HandlerFunc{}}
	s := &Server{Router:r}
	s.middlewares = []Middleware{
		LogHandler,
		RecoverHandler,
		staticHandler,
		ParseJsonBodyHandler,
		ParseFormHandler,
	}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{
		Params: map[string]interface{}{},
		ResponseWriter: w,
		Request: 		r,
	}
	for k, v := range r.URL.Query() {
		c.Params[k] = v[0]
	}
	s.startHandler(c)
}

func (s *Server) Run(addr string) {
	// startHandler를 라우터 핸들러 함수로 지정
	s.startHandler = s.Router.handler()

	// 등록된 미들웨어를 라우터 핸들러 앞에 하나씩 추가
	for i:=len(s.middlewares)-1; i>=0; i-- {
		s.startHandler = s.middlewares[i](s.startHandler)
	}

	// 웹 서버 시작
	if err:=http.ListenAndServe(addr, s); err!=nil {
		panic(err)
	}
}