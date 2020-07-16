package serializer

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type ToLogrusField struct {}

func (s ToLogrusField) Request(r *http.Request) logrus.Fields {
	return logrus.Fields{
		"method": r.Method,
		"path": r.URL.Path,
		"client_ip": r.getClientIP(),
		"headers": r.serializeRequestHeader(),
	}
}