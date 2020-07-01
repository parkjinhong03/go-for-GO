package user

import "github.com/stretchr/testify/mock"

type testDAO struct {
	mock *mock.Mock
}

func NewTestDAO(mock *mock.Mock) *testDAO {
	return &testDAO{
		mock: mock,
	}
}
