package natsEncoder

import (
	"MSA.example.com/1/protocol"
	"MSA.example.com/1/proxy"
	"encoding/json"
	"errors"
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

func (e *jsonEncoder) Encode(v interface{}) error {
	switch e.proxy.(type) {
	case *proxy.AuthServiceProxy:
		if _, ok := v.(protocol.AuthSignUpProtocol); !ok {
			return errors.New("this object cannot be encoded in your proxy")
		}
	}

	return e.Encoder.Encode(v)
}