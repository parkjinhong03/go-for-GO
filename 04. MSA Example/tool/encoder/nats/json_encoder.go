package natsEncoder

import (
	"MSA.example.com/1/protocol"
	"MSA.example.com/1/tool/proxy"
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
		switch v.(type) {
		case protocol.AuthSignUpRequestProtocol:
		default:
			return errors.New("this object cannot be encoded in your AuthServiceProxy")
		}
	case *proxy.ApiGatewayProxy:
		switch v.(type) {
		case protocol.ApiGatewaySignUpResponseProtocol:
		default:
			return errors.New("this object cannot be encoded in your ApiGatewayProxy")
		}
	case *proxy.UserServiceProxy:
		switch v.(type) {
		case protocol.UserRegistryPublishProtocol:
		default:
			return errors.New("this object cannot be encoded in your UserServiceProxy")
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