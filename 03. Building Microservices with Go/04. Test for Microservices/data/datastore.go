package data

// 어떤 객체가 데이터 저장소로 쓰이기 위해 있어야 하는 메서드를 정의하는 인터페이스
type Store interface {
	// Search 메서드는 매개 변수로 전달된 문자열과 같은 이름의 새끼 고양이의 목록을 반환한다.
	Search(name string) []Kitten
}