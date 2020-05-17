package natsEncoder

import (
	"encoding/json"
	"io"
)

type jsonEncoder struct {
	*json.Encoder
	proxy io.Writer
}

func NewJsonEncoder(proxy io.Writer) *jsonEncoder {
	return &jsonEncoder{
		Encoder: json.NewEncoder(proxy),
		proxy:   proxy,
	}
}