package throttling

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
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

// 최대 요청의 수가 1인 서비스에 2개의 요청을 보냈을 때 429를 반환하는지에 대한 테스트 함수
func TestReturnsBusyWhenConnectionsExhausted(t *testing.T) {
	// 원하는 타이밍에 핸들러를 종료시킬 수 있게 하기 위해서 context.WithCancel 함수를 사용한다.
	ctx1, cancel1 := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())
	handler := NewLimitMiddleware(1, newTestHandler())
	rw1, r1 := setupTest(ctx1)
	rw2, r2 := setupTest(ctx2)

	// 10ms 뒤에 context.WithCancel 함수로부터 받을 함수를 호출시켜 핸들러에게 종료 신호를 보낸다.
	// 이렇게 10ms 라는 시간을 둔 이유는 첫 번째 요청이 바로 끝나버리면 해당 핸들러가 429를 반환하는지 모르기 때문이다.
	time.AfterFunc(10 * time.Millisecond, func() {
		cancel1()
		cancel2()
	})

	// 각 핸들러가 완료되기 전에 테스트가 끝나버려지지 않기(동기화를) 위해서 sync.WaitGroup 객체를 선언한다.
	waitGroup := sync.WaitGroup{}
	// 우리는 두 개의 핸들러를 실행시키기 때문에 WaitGroup.Add 메서드를 이용하여 2개의 대기 그룹을 추가시킨다.
	waitGroup.Add(2)

	// 두 개의 요청이 동시에 실행되게 하기 위하여 각각 고루틴을 이용하여 병행처리를 한다.
	go func() {
		handler.ServeHTTP(rw1, r1)
		// WaitGroup.Done 메서드를 이용하여 핸들러가 완료되었다는 의미로 대기 그룹을 하나 감소시킬 수 있다.
		waitGroup.Done()
	}()
	go func() {
		handler.ServeHTTP(rw2, r2)
		waitGroup.Done()
	}()

	// WaitGroup.Wait 메서드를 이용하여 대기 그룹이 모두 감소될 때 까지 프로세스를 대기시킬 수 있다.
	waitGroup.Wait()

	// 모든 핸들러가 종료되면, 그제서야 단정문을 작성하여 테스트를 완료시킨다.
	if rw1.Code != http.StatusOK || rw2.Code != http.StatusTooManyRequests {
		t.Fatalf("\nOne request should be busy \nrequest1: %d \nrequest2: %d", rw1.Code, rw2.Code)
	}
}