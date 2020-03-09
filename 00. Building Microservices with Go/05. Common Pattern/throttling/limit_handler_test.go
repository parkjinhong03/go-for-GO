package throttling

import (
	"context"
	"net/http"
	"net/http/httptest"
)

// LimitMiddleware 후에 실행시킬 모의 핸들러를 생성하여 반환시키는 함수
func newTestHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		// r.Context()는 WithCancel 컨텍스트로, 첫 번째 요청을 대기 상태로 두어야 하기 때문에 해당 구문을 실행시켜야 한다.
		<-r.Context().Done()
	}
}

// 테스트를 하기 위한 모의 객체를 생성한 후 매개 변수로 받은 컨텍스트를 적용해서 반환하는 함수이다.
func setupTest(ctx context.Context) (*httptest.ResponseRecorder, *http.Request) {
	r := httptest.NewRequest("GET", "/", nil)
	// 핸들러에서 접근해야 하는 컨텍스트를 WithContext 매세드를 이용하여 요청에 적용시킬 수 있다.
	r = r.WithContext(ctx)
	return httptest.NewRecorder(), r
}

