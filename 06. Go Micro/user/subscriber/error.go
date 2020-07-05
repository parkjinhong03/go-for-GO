package subscriber

import "errors"

var (
	ErrorBadRequest = errors.New("bad request")
	ErrorForbidden = errors.New("request forbidden")
	ErrorMsgDuplicated = errors.New("massage duplicated")
)