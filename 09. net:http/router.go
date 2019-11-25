package main

import "net/http"

type router struct {
	// 키: http 메서드
	// 값: URL 패턴별로 실행할 HandleFunc
	handlers map[string]map[string]http.HandlerFunc
}

// 라우터에 핸들러를 등록하기 위한 메서드 정의
func (r *router) HandleFunc(method, pattern string, h http.HandlerFunc) {
	// http 메서드로 등록된 맵이 있는지 확인
	m, ok := r.handlers[method]
	if !ok {
		// 등록된 맵이 없으면 새 맵을 생성
		m = make(map[string]http.HandlerFunc)
		r.handlers[method] = m
	}

	// http 메서드로 등록된 맵에 URL 패턴과 핸들러 함수 등록
	m[pattern] = h
}

// http.Handler 인터페이스로 사용하기 위한 ServeHTTP(http.ResponseWriter, *http.Request) 메서드 정의ß
