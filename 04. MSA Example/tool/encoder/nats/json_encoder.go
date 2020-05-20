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
	default:
		return errors.New("it is a proxy or incorrect proxy that this object cannot handle")
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = e.proxy.Write(b)
	return err
}