package user

import (
	"auth/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

var AuthArr []model.Auth

type testDAO struct {
	mock *mock.Mock
}

func NewTestDAO(mock *mock.Mock) *testDAO {
	return &testDAO{
		mock: mock,
	}
}

func (td *testDAO) Insert(auth *model.Auth) (result *model.Auth, err error) {
	td.mock.Called(auth)

	for _, a := range AuthArr {
		if a.UserId != (*auth).UserId { continue }
		err = IdDuplicateError
		return
	}

	auth.Status = CreatePending
	auth.ID = uint(len(AuthArr) + 1)
	AuthArr = append(AuthArr, *auth)
	result = auth
	return
}

func (td *TestDAO) Commit() *gorm.DB {
	td.Mock.Called()
	return nil
}

func (td *TestDAO) Rollback() *gorm.DB {
	td.Mock.Called()
	return nil
}