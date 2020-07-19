package logrusfield

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func ForReturn(v interface{}, c int, err error) logrus.Fields {
	var es string
	if err != nil {
		es = err.Error()
	}

	b, err := json.Marshal(v)
	if err != nil {
		b = []byte{}
	}

	return logrus.Fields{
		"json":    string(b),
		"outcome": c,
		"error":   es,
	}
}

func ForHandleRequest(r *http.Request, cip string) logrus.Fields {
	return logrus.Fields{
		"method":       r.Method,
		"path":         r.URL.Path,
		"client_ip":    cip,
		"X-Request-Id": r.Header.Get("X-Request-Id"),
		"header":       encodeHeader(r.Header),
	}
}

func encodeHeader(h http.Header) string {
	b, err := json.Marshal(h)
	if err != nil {
		b = []byte{}
	}

	return string(b)
}