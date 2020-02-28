package data

import "github.com/stretchr/testify/mock"

// 모의 객체 프레임워크 Testify를 사용하여 데이터 저장소의 모의 객체를 만든다.
type MockStore struct {
	// mock.Mock 타입을 임베디드 필드로 포함하고 있기 때문에 해당 구조체의 모든 메서드를 사용할 수 있다.
	mock.Mock
}

// 데이터 저장소로 사용하기 위해 Search 메서드를 정의한다.
func (m *MockStore) Search(name string) []Kitten {
	// Called 메서드를 호출하면 모의 저장소 객체 mock를 초기 설정했을 때 제공한 인수의 목록을 반환한다.
	args := m.Mock.Called(name)

	// Called 메서드에서 반환한 목록의 0번째 인덱스를 추출하여 []Kitten 타입으로 assertion 한 후 반환한다.
	return args.Get(0).([]Kitten)
}