package features

import "github.com/cucumber/godog"

func iHaveNoSearchCriteria() error {
	return nil
}

func iCallTheSearchEndpointSearch() error {
	return nil
}

func iShouldReceiveABadRequestMessage() error {
	return nil
}

func iHaveValidSearchCriteria() error {
	return godog.ErrPending
}

func iShouldReceiveAListOfKittens() error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I have no search criteria \(검색 기준이 없다\)$`, iHaveNoSearchCriteria)
	s.Step(`^I call the search endpoint \(search 엔드 포인트를 호출한다\)$`, iCallTheSearchEndpointSearch)
	s.Step(`^I should receive a bad request message \(잘못된 요청이라는 메세지를 받는다\)$`, iShouldReceiveABadRequestMessage)
	s.Step(`^I have valid search criteria \(유효한 검색 기준이 있다\)$`, iHaveValidSearchCriteria)
	s.Step(`^I should receive a list of kittens \(새끼 고양이의 목록을 받는다\)$`, iShouldReceiveAListOfKittens)
}