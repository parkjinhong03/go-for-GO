package writers

import (
	"compress/gzip"
	"net/http"
)

// 응답 작성을 위해 마지막에 호출할 Write() 메서드를 가지고있는 gzip 인코딩을 위한 Response Writer 객체 정의
type GzipResponseWriter struct {
	Gw *gzip.Writer
	http.ResponseWriter
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	if _, ok := w.Header()["Content-Type"]; !ok {
		// 응답 헤더의 Content-Type가 설정되어 있지 않은 경우, 인코딩되지 않은 응답 내용에서 컨텐트 타입을 유추하여 적용한다.
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}

	// gzip.Writer의 Write() 메서드를 이용하여 응답 본문을 gzip으로 인코딩하여 응답하는것으로 연결을 마무리한다.
	return w.Gw.Write(b)
}