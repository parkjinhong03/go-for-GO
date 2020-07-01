package user

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"user/model"
)

type testDAO struct {
	mock *mock.Mock
}

func NewTestDAO(mock *mock.Mock) *testDAO {
	return &testDAO{
		mock: mock,
	}
}

func (t *testDAO) InsertUser(u *model.User) (result *model.User, err error) {
	args := t.mock.Called(u)
	return args.Get(0).(*model.User), args.Error(1)
}

func (t *testDAO) InsertMessage(msg *model.ProcessedMessage) (result *model.ProcessedMessage, err error) {
	args := t.mock.Called(msg)
	return args.Get(0).(*model.ProcessedMessage), args.Error(1)
}

func (t *testDAO) CheckIfEmailExist(email string) (exist bool, err error) {
	args := t.mock.Called(email)
	return args.Bool(0), args.Error(1)
}

func (t *testDAO) Commit() (db *gorm.DB) {
	args := t.mock.Called()
	return args.Get(0).(*gorm.DB)
}

func (t *testDAO) Rollback() (db *gorm.DB) {
	args := t.mock.Called()
	return args.Get(0).(*gorm.DB)
}