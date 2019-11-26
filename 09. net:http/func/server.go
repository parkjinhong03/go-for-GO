package net

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