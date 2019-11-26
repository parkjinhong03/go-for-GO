package net

import (
	"net/http"
	"strings"
)

type Router struct {
	// 키: http 메서드
	// 값: URL 패턴별로 실행할 HandleFunc
	Handlers map[string]map[string]http.HandlerFunc
}

// 라우터에 핸들러를 등록하기 위한 메서드 정의
func (r *Router) HandleFunc(method, pattern string, h http.HandlerFunc) {
	// http 메서드로 등록된 맵이 있는지 확인
	m, ok := r.Handlers[method]
	if !ok {
		// 등록된 맵이 없으면 새 맵을 생성
		m = make(map[string]http.HandlerFunc)
		r.Handlers[method] = m
	}

	// http 메서드로 등록된 맵에 URL 패턴과 핸들러 함수 등록
	m[pattern] = h
}

// http.Handler 인터페이스로 사용하기 위한 ServeHTTP(http.ResponseWriter, *http.Request) 메서드 정의
// ServeHTTP 메서드는 웹 요청의 http 메서드와 URL 경로를 분석해서 그에 맞는 핸들러를 찾아 동작시킨다.
func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// http 메서드에 맞는 모든 handlers를 반복하여 요청 URL에 해당하는 handler를 찾음
	for patterns, handler := range r.Handlers[req.Method] {
		if ok, _ := match(patterns, req.URL.Path); ok {
			// Context 생성

			// 요청 URL에 해당하는 handler 실행
			handler(w, req)
			return
		}
	}
	http.NotFound(w, req)
	return
}

func match(pattern, path string) (bool, map[string]string) {
	// 패턴과 패스가 정확히 일치하면 바로 true를 반환
	if pattern == path {
		return true, nil
	}

	// 패턴과 패스를 "/" 단위로 구분
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	// 패턴과 패스를 "/"로 구분한 후 부분 문자열 집합의 개수가 다르면 false를 반환
	if len(patterns) != len(paths) {
		return false, nil
	}

	// 패턴에 일치하는 URL 매개변수를 담기 위한 params 맵 생성
	params := make(map[string]string)

	// "/"로 구분된 패턴/패스의 각 문자열을 하나씩 비교
	for i:=0; i<len(patterns); i++ {
		switch {
			case patterns[i] == paths[i]:
				// 패턴과 패스의 부분 문자열이 일치하면 바로 다음 루프 실행
			case len(patterns[i]) > 0 && patterns[i][0] == ':':
				// 패턴이 ':' 문자로 시작하면 params에 URL params를 담은 후 다음 루프 실행≈
				params[patterns[i][1:]] = paths[i]
			default:
				// 일치하는 경우가 없다면 false를 반환
				return false, nil
		}
	}

	// true와 params를 반
	return true, params
}