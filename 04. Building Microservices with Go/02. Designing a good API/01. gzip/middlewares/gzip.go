package middlewares

import (
	"../entities"
	"../writers"
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// HelloWorld 핸들러를 실행하기 전에 gzip을 이용하여 ResponseWriter를 작성하기 위한 설정을 하는 미들웨어이다.
type GzipMiddleWare struct {
	next http.Handler
}

func NewGzipMiddleWare(next http.Handler) *GzipMiddleWare {
	return &GzipMiddleWare{next: next}
}

// HelloWorld 핸들러보다 먼저 실행되는 미들웨어로, 이 미들웨어의 로직에 따라 다음에 실행될 핸들러의 처리 방식이 달라진다.
func (h *GzipMiddleWare) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request entities.HelloWorldRequest
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request)
	context.WithValue(r.Context(), "name", request.Name)

	// 클라이언트가 제공하는 인코딩 형식들을 헤더로부터 빼온다.
	encodings := r.Header.Get("Accept-Encoding")

	// Accept-Encoding에 gzip이 포함되어 있으면 응답을 gzip 인코딩으로 하도록 설정한다.
	if strings.Contains(encodings, "gzip") {
		h.serveGzipped(rw, r)
		// deflate 인코딩에 관한 설정은 아직 하지 않았으므로 500 에러를 발생시킨다.
	} else if strings.Contains(encodings, "deflate") {
		panic("Deflate not implemented")
		// 그 둘도 아니면 아무 인코딩 없이 일반적으로 응답을 처리한다.
	} else {
		h.servePlain(rw, r)
	}
}

func (h *GzipMiddleWare) serveGzipped(rw http.ResponseWriter, r *http.Request) {
	// json.NewEncoder() 함수와 비슷하게 gzip.NewWriter() 함수를 사용하여 응답 작성 객체를 생성한다.
	gzw := gzip.NewWriter(rw)
	defer gzw.Close()

	// 응답 헤더ㅡㄹ 'Content-Encoding: gzip'으로 설정한다.
	rw.Header().Set("Content-Encoding", "gzip")
	// gzip 인코당을 하여 응답을 마무리하는 Response Writer로 대체하여 HelloWorldHandler를 실행시킨다.
	h.next.ServeHTTP(writers.GzipResponseWriter{Gw: gzw, ResponseWriter: rw}, r)
}

func (h *GzipMiddleWare) servePlain(rw http.ResponseWriter, r *http.Request) {
	h.next.ServeHTTP(rw, r)
}