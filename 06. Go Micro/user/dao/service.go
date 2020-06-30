package dao

import (
	"github.com/jinzhu/gorm"
	"user/model"
)

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

