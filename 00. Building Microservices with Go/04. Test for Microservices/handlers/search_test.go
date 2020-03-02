// 이 파일은 search.go 핸들러를 유닛 테스트하기 위한 파일로 동일한 폴더에 _test를 뒤에 붙혀 만드는 것이 원칙이다.
// 참고로 이 테스트 코드는 실제 웹 서버 코드를 작성하기 전에 작성한 실패하는 단위 테스트이다.
// 또한 물리적인 웹 서버를 만들지 많아도 핸들러를 실행할 수 있는 코드로, 웹 서버를 통한 테스트보다 훨씬 빠르게 실행될 것 이다.

package handlers

import (
	"../data"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 테스트 하는 동안에 사용해야 하는 데이터 저장소 모의 객체를 전역 변수로 선언한다.
var mockStore = data.NewMockStore(nil)

// 테스트 메서드의 이름은 Test로 시작해 바로 대문자 또는 숫자가 나오는 특수한 이름이어야 한다.
// 해당 테스트는 검색 기준이 요청과 함께 전송됐는지 확인하는 테스트이다.
func TestSearchHandlerReturnsBadRequestWhenNoSearchCriteriaIsSent(t *testing.T) {
	handler := SearchHandler{}

	// net/http/httptest 패키지의 NewRequest, NewResponse 메서드는 의존성을 제거하기 위한 편리한 메서드이다.
	// 즉, 실제로 요청을 받지 않았어도 마치 요청을 받은 것만 같은(의존성 제거) 상황을 표출하기 위해 사용하는 메서드이다.
	// 이를 통해 의존성이 있는 객체인 http.Request와 http.ResponseWriter의 모의 객체 버전을 생성한 것 이다.
	request := httptest.NewRequest("GET", "/search", nil)
	response := httptest.NewRecorder()

	// 실제 요청시에만 실행되는 ServeHTTP 메서드를 인위적으로 호출하여 핸들러를 실행시킨다.
	// 그리고 핸들러 실행 결과(응답)는 참조 변수로 넘긴 response에 저장되어 이 객체를 단정문(assertion)을 작성할 수 있다.
	handler.ServeHTTP(response, request)

	// 검색 기준을 요청에 포함시키지 않았는데도 응답 상태 코드가 BadRequest(400)이 아니라면 t.Errorf 함수를 호출해 테스트를 실패시킨다.
	// 하지만 테스트 코드 실행중 t.Fail 함수가 호출되지 않는다면 해당 테스트는 성공으로 끝난다.
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected BadRequest got %v", response.Code)
	}
}

// 해당 테스트는 쿼리에 빈 문자열이 있는 요청을 받았을때 400 Bad Request를 반환하는가를 검사한다.
func TestSearchHandlerReturnsBadRequestWhenBlackSearchCriteriaIsSent(t *testing.T) {
	r, rw, handler := setupTest(searchRequest{})

	handler.ServeHTTP(rw, r)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected BadRequest got %v", rw.Code)
	}
}

// 해당 테스트는 검색 기준을 옳바르게 넘겨주었을 때, 원하는 메서드가 호출되어 원하는 값을 반환하는지 테스트 하는 메서드이다.
func TestSearchHandlerReturnsKittensWithValidQuery(t *testing.T) {
	r, rw, handler := setupTest(searchRequest{Query: "Fat Freddy's Cat"})
	// mock.On 메서드를 사용하여 일정 메서드가 호출되었는지 확인할 수 있다.
	// 다음 코드는  "Fat Freddy Cat" 매개변수를 넘겨 Search 메서드를 호출할 경우, []data.Kitten 타입의 객체를 반환한다는 단정문이다.
	mockStore.On("Search", "Fat Freddy's Cat").Return([]data.Kitten{})

	handler.ServeHTTP(rw, r)
	// mock.AssertExpectations 메서드는 위에서 mock.On 메서드로 등록한 단정문대로 핸들러가 실행되지 않았을 경우 테스트 에러를 발생시킨다.
	mockStore.AssertExpectations(t)

	response := searchResponse{}
	_ = json.Unmarshal(rw.Body.Bytes(), &response)

	// testify/assert 패키지의 Equal 함수는 2, 3번쨰로 넘긴 매개변수가 서로 다를 경우 테스트 에러를 발생시킨다.
	assert.Equal(t, 1, len(response.Kittens))
	assert.Equal(t, http.StatusOK, rw.Code)
}

func setupTest(d interface{}) (*http.Request, *httptest.ResponseRecorder, SearchHandler) {
	h := SearchHandler{
		DataStore: mockStore,
	}
	rw := httptest.NewRecorder()

	if d == nil {
		return httptest.NewRequest("GET", "/search", nil), rw, h
	}

	body, _ := json.Marshal(d)
	return httptest.NewRequest("GET", "/search", bytes.NewReader(body)), rw, h
}