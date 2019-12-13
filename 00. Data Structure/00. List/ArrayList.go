package main

const ListLen = 100
type LData interface{}

// 배열기반 리스트를 정의한 구조체
type List struct {
	arr [ListLen]LData	// 리스트의 저장소인 배열
	numOfData int		// 저장된 데이터의 수
	curPosition int		// 데이터 참조위치를 기록
}
