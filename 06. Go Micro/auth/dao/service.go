package dao

import (
	"auth/dao/auth"
	"auth/model"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type AuthDAOCreator struct {
	db *gorm.DB
}

func NewAuthDAOCreator(db *gorm.DB) (adc *AuthDAOCreator) {
	adc = &AuthDAOCreator{
		db: db,
	}
	if db == nil {
		return
	}
	if !db.HasTable(&model.Auth{}) {
		db.CreateTable(&model.Auth{})
	}
	db.AutoMigrate(&model.Auth{}, &model.ProcessedMessage{})
	return
}

type AuthDAOService interface {
	InsertAuth(*model.Auth) (result *model.Auth, err error)
	UpdateStatus(id uint, status string) (err error)
	CheckIfUserIdExist(id string) (exist bool, err error)
	Commit() *gorm.DB
	Rollback() *gorm.DB
	InsertMessage(*model.ProcessedMessage) (result *model.ProcessedMessage, err error)
}

func (dc *AuthDAOCreator) GetDefaultAuthDAO() AuthDAOService {
	tx := dc.db.Begin()
	return user.NewDefaultDAO(tx)
}

func (dc *AuthDAOCreator) GetTestAuthDAO(mock *mock.Mock) AuthDAOService {
	return user.NewTestDAO(mock)
}