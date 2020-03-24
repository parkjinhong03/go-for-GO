// search.feature 시나리오에서 각각의 단계마다 어떤 상테 또는 어떤 행동을 해야지 해당 단계를 통과할 수 있는지를 함수로 정의한 파일이다.
// 함수의 시그니처는 func() error 이고, error로 nil을 반환하면 해당 단계는 통과되는 것 이다.
// 시나리오와 함수는 godog.Suite.Step 메서드를 이용하여 연결시킬 수 있다.

package features

import (
	"../data"
	"../handlers"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/gherkin"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

// godog.Suite 객체에 정의된 메서드의 서명을 마음대로 바꿀 수 없으므로, 해당 메서드에서 다른 객체에 접근하기 위해선 접근할 객체들을 전역변수로 선언해야 한다.
// 또힌 각각의 단계(함수)마다 연관성을 지키기위해 보존되어야 하는 데이터들은 전역변수로 선언한다.
var (
	response *httptest.ResponseRecorder
 	criteria handlers.SearchRequest
 	handler handlers.SearchHandler
 	store *data.MongoStore
	err error
)

func iHaveNoSearchCriteria() error {
	criteria = handlers.SearchRequest{}
	return nil
}

func iCallTheSearchEndpointSearch() error {
	response = httptest.NewRecorder()
	var request *http.Request

	if criteria.Query == "" {
		request = httptest.NewRequest("get", "/", nil)
	} else {
		body, _ := json.Marshal(criteria)
		request = httptest.NewRequest("get", "/", bytes.NewReader(body))
	}

	handler.ServeHTTP(response, request)
	return nil
}

func iShouldReceiveABadRequestMessage() error {
	if response.Code != http.StatusBadRequest {
		return fmt.Errorf("should have received a bad response")
	}
	return nil
}

func iHaveValidSearchCriteria() error {
	criteria = handlers.SearchRequest{
		Query: "Fat Freddy's Cat",
	}
	return nil
}

func iShouldReceiveAListOfKittens() error {
	var body handlers.SearchResponse
	decoder := json.NewDecoder(response.Body)
	err := decoder.Decode(&body)

	if len(body.Kittens) < 1 || err != nil {
		return fmt.Errorf("should have receive a list of kittens")
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	// godog.Suite.Step 메서드를 이용하여 시나리오의 각 단계와 위에서 만든 함수를 연결할 수 있다.
	s.Step(`^I have no search criteria \(검색 기준이 없다\)$`, iHaveNoSearchCriteria)
	s.Step(`^I call the search endpoint \(search 엔드 포인트를 호출한다\)$`, iCallTheSearchEndpointSearch)
	s.Step(`^I should receive a bad request message \(잘못된 요청이라는 메세지를 받는다\)$`, iShouldReceiveABadRequestMessage)
	s.Step(`^I have valid search criteria \(유효한 검색 기준이 있다\)$`, iHaveValidSearchCriteria)
	s.Step(`^I should receive a list of kittens \(새끼 고양이의 목록을 받는다\)$`, iShouldReceiveAListOfKittens)

	// Feature(전체 시나리오)에 대한 테스트를 실행시키기 전에 실행할 함수를 등록하는 함수이다.
	s.BeforeFeature(func(*gherkin.Feature) {
		setHandler()
		clearDB()
		setupDB()
	})
}

// 테스트 환경애서 모의 객체로 ServeHTTP를 실행하기 전에 런타임 에러가 나지 않게 handler와 store의 값을 채워주는 함수
func setHandler() {
	serverURL := "localhost"
	if os.Getenv("DOCKER_IP") != "" {
		serverURL = os.Getenv("DOCKER_IP")
	}

	for i:=0; i<=10; i++ {
		store, err = data.NewMongoStore(serverURL)
		if err == nil {
			handler = handlers.SearchHandler{
				DataStore: store,
			}
			return
		}

		time.Sleep(1 * time.Second)
	}

	panic("Can't connect session to mongoDB server")
}

func clearDB() {
	store.DeleteAllKittens()
}

func setupDB() {
	store.InsertKittens(
		[]data.Kitten{
			{
				Id:     1,
				Name:   "Felix",
				Weight: 12.0,
			}, {
				Id:     2,
				Name:   "Fat Freddy's Cat",
				Weight: 20.0,
			}, {
				Id:     3,
				Name:   "Garfield",
				Weight: 35.0,
			},
		})
}