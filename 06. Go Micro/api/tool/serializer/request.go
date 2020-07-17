package serializer

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strings"
)

type ToLogrusField struct {}

func (s ToLogrusField) Request(r *http.Request) logrus.Fields {
	return logrus.Fields{
		"method": r.Method,
		"path": r.URL.Path,
		"client_ip": s.getClientIP(r),
		"X-Request-Id": r.Header.Get("X-Request-Id"),
		"header": s.encodeHeaderToString(r.Header),
	}
}

func (s ToLogrusField) encodeHeaderToString(h http.Header) string {
	out, err := json.Marshal(h)
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}

func (s ToLogrusField) getClientIP(r *http.Request) string {
	var ip string

	if forwarded := r.Header.Get("X-FORWARDED-FOR"); forwarded != "" {
		ip = forwarded
		return ip
	}

	ip = r.RemoteAddr
	if strings.Contains(ip, "[::1]") {
		ip = strings.Replace(ip, "[::1]", "localhost", 1)
	}
	ip = strings.Split(ip, ":")[0]
	return ip
}