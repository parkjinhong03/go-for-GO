package user

import (
	"auth/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type TestDAO struct {
	mock.Mock
}

func (td *TestDAO) Insert(a *model.Auth) (*model.Auth, error) {
	td.Mock.Called(a)
	return &model.Auth{}, nil
}

func (td *TestDAO) Commit() *gorm.DB {
	td.Mock.Called()
	return nil
}

func (td *TestDAO) Rollback() *gorm.DB {
	td.Mock.Called()
	return nil
}