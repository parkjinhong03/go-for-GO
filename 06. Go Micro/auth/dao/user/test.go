package user

import (
	"auth/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type testDAO struct {
	mock *mock.Mock
}

func NewTestDAO(mock *mock.Mock) *testDAO {
	return &testDAO{
		mock: mock,
	}
}

func (td *testDAO) Insert(auth *model.Auth) (*model.Auth, error) {
	args := td.mock.Called(auth)
	return args.Get(0).(*model.Auth), args.Error(1)
}

func (td *testDAO) CheckIfUserIdExists(id string) (bool, error) {
	args := td.mock.Called(id)
	return args.Bool(0), args.Error(1)
}

func (td *testDAO) Commit() *gorm.DB {
	return td.mock.Called().Get(0).(*gorm.DB)
}

func (td *testDAO) Rollback() *gorm.DB {
	return td.mock.Called().Get(0).(*gorm.DB)
}