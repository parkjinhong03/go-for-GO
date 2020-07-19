package logrusfield

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
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