package middleware

import (
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type Middleware struct {
	ajc, ujc jaegercfg.Configuration
}

func New(ajc, ujc jaegercfg.Configuration) Middleware {
	return Middleware{
		ajc: ajc,
		ujc: ujc,
	}
}