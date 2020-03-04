// 해당 파일은 go의 벤치마크 테스트를 이용하기 위한 파일이다.
// 벤치마크란 각 함수마다의 실행속도를 비교하여 어느 함수의 기능이 더 좋은지를 비교할 수 있는 테스트이다.
// "go test -bench=. -benchmem" 명령어를 이용하면 아래 나온 규칙을 지킨 함수들을 벤치마크 테스트 한다.
//	1. 함수 이름은 Benchmark로 시작해야 한다.
//  2. Benchmark 뒤에는 대문자가 와야 한다.
//  3. *testing.B 매개변수를 받아햐 한다.

package handlers

import (
	"../data"
	"bytes"
	"net/http/httptest"
	"testing"
)

// 서로 다른 두 벤치마크 함수에서 중복된 코드를 묶어서 다음과 같이 함수로 만들 수 있다.
func testInit() SearchHandler {
	mockStore = data.NewMockStore(nil)
	mockStore.On("Search", "Fat Freddy's Cat").Return([]data.Kitten{
		{
			Name: "Fat Freddy's Cat",
		},
	})

	return SearchHandler{DataStore: mockStore}
}

// 아래 나온 두 벤치마크 테스트 함수는 둘 다 모의 객체를 이용하여 가상 요청을 보내는 함수이다.
// 하지만 하나는 반환값을 변수에 담고, 나머지는 바로 매개변수로 넘긴다.
// 두 개의 함수를 벤치마크한 결과, 실제로 두 함수의 속도는 20000 ns/op 정도로 속도 차이가 거의 없다.
func BenchmarkSearchHandler1(b *testing.B) {
	handler := testInit()

	// 한번의 프로세스로는 속도를 비교하기 어려우므로 Go가 측정하여 설정해준 b 객체의 N 필드 값만큼 반복하여 실행한다.
	for i:=0; i<b.N; i++ {
		request := httptest.NewRequest("POST", "/search", bytes.NewReader([]byte(`{"query":"Fat Freddy's Cat"}`)))
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
	}
}

func BenchmarkSearchHandler2(b *testing.B) {
	handler := testInit()

	for i:=0; i<b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/search", bytes.NewReader([]byte(`{"query":"Fat Freddy's Cat"}`))))
	}
}
