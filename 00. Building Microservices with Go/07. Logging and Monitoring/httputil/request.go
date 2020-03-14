package httputil

import (
	"encoding/json"
	"net/http"
	"strings"
)

// RequestSerializer는 http.Request를 로깅에 필요한 정보만 뽑아서 json으로 Serialize하는 객체이다.
type RequestSerializer struct {
	*http.Request
}

// rs.serialize 메서드의 반환 객체(serializedRequest)를 마샬링하여 문자열로 반환하는 함수이다.
func (rs *RequestSerializer) ToJSON() string {
	bytes, _ := json.Marshal(rs.serializeRequest())
	return string(bytes)
}

type serializedHeader struct {
	Key string
	Value string
}

type serializedRequest struct {
	Method string
	Host string
	Path string
	Headers []serializedHeader
}

// http.Request 객체중에서 로깅에 필요한 필드만 골라서 serializedRequest 객체에 저장하여 반환하는 함수이다.
func (rs *RequestSerializer) serializeRequest() serializedRequest {
	var headers []serializedHeader
	for k, v := range rs.Header {
		// strings.Join 함수를 이용하여 슬라이스로 되어있는 Header의 Value를 ", "를 사이사이에 두어 하나의 문자열로 변환시킬 수 있다.
		headers = append(headers, serializedHeader{Key: k, Value: strings.Join(v, ", ")})
	}

	return serializedRequest{
		Method:  rs.Method,
		Host:    rs.Host,
		Path:    rs.URL.Path,
		Headers: headers,
	}
}