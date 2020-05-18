package customError

import (
	"github.com/nats-io/nats.go"
)

type ProxyWriteError struct {
	Err 		error
	ReturnMsg	*nats.Msg
}

func (pwe *ProxyWriteError) Error() string {
	return pwe.Err.Error()
}