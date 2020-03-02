# 기능: 사용자인 내가 search 엔드 포인트를 호출하면, 나는 새끼 고양이의 목록을 받을 것이다.
Feature: As a user when I call the search endpoint, I would like to receive a list of kittens

  Scenario: Invalid query (유효하지 않은 쿼리)
    Given I have no search criteria (검색 기준이 없다)
    When I call the search endpoint (search 엔드 포인트를 호출한다)
    Then I should receive a bad request message (잘못된 요청이라는 메세지를 받는다)

  Scenario: Valid query (유효한 쿼리)
    Given I have valid search criteria (유효한 검색 기준이 있다)
    When I call the search endpoint (search 엔드 포인트를 호출한다)
    Then I should receive a list of kittens (새끼 고양이의 목록을 받는다)