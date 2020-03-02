package data

import (
	"github.com/stretchr/testify/mock"
)

func NewMockStore(init []Kitten) *MockStore {
	var Kittens []Kitten

	if init == nil {
		Kittens = []Kitten{
			{
				Id:     1,
				Name:   "Felix",
				Weight: 12.3,
			}, {
				Id:     2,
				Name:   "Fat Freddy's Cat",
				Weight: 20.0,
			}, {
				Id:     3,
				Name:   "Garfield",
				Weight: 35.0,
			},
		}
	} else {
		Kittens = init
	}

	return &MockStore{Kittens: Kittens}
}

// 모의 객체 프레임워크 Testify를 사용하여 데이터 저장소의 모의 객체를 만든다.
type MockStore struct {
	// mock.Mock 타입을 임베디드 필드로 포함하고 있기 때문에 해당 구조체의 모든 메서드를 사용할 수 있다.
	mock.Mock
	Kittens []Kitten
}

// 데이터 저장소로 사용하기 위해 Search 메서드를 정의한다.
func (m *MockStore) Search(name string) []Kitten {
	// mock.Called 메서드를 이용하여 search_test.go 에서 mock.On 메서드로 등록한 단정문을 실행시켰다는 것을 알려줄 수 있다.
	m.Mock.Called(name)

	var kittens []Kitten

	for _, k := range m.Kittens {
		if k.Name == name {
			kittens = append(kittens, k)
		}
	}

	// 만약 []Kittens 타입이 아닌 객체를 반환하였을 경우, mock.On 메서드에서 
	return kittens
}