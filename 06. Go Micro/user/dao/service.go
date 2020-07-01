package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"user/dao/user"
	"user/model"
)

type UserDAOService interface {
	InsertUser(*model.User) (result *model.User, err error)
	InsertMessage(*model.ProcessedMessage) (result *model.ProcessedMessage, err error)
	CheckIfEmailExist(email string) (exist bool, err error)
	Commit() (db *gorm.DB)
	Rollback() (db *gorm.DB)
}

type UserDAOCreator struct {
	db *gorm.DB
}

func NewUserDAOCreator(db *gorm.DB) (udc *UserDAOCreator) {
	udc = &UserDAOCreator{
		db: db,
	}

	if !db.HasTable(&model.User{}) {
		db.CreateTable(&model.User{})
	}
	if !db.HasTable(&model.ProcessedMessage{}) {
		db.CreateTable(&model.ProcessedMessage{})
	}

	db.AutoMigrate(&model.User{}, &model.ProcessedMessage{})
	return
}

func (udc *UserDAOCreator) GetDefaultUserDAO() UserDAOService {
	tx := udc.db.Begin()
	return user.NewDefaultDAO(tx)
}

func (udc *UserDAOCreator) GetTestUserDAO(mock *mock.Mock) UserDAOService {
	return user.NewTestDAO(mock)
}