package user

import (
	"auth/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type TestDAO struct {
	mock.Mock
	auths []model.Auth
}

func (td *TestDAO) Insert(auth *model.Auth) (result *model.Auth, err error) {
	td.Mock.Called(auth)

	for _, a := range td.auths {
		if a != *auth { continue }
		err = IdDuplicateError
		return
	}

	auth.Status = CreatePending
	auth.ID = uint(len(td.auths) + 1)
	td.auths = append(td.auths, *auth)
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