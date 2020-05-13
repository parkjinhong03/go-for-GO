package userdata

import (
	"MSA.example.com/1/model"
	"github.com/jinzhu/gorm"
)

type userDAO struct {
	db *gorm.DB
}

func GetUserDAO(db *gorm.DB) *userDAO {
	if !db.HasTable(&model.Users{}) {
		db.CreateTable(&model.Users{})
	}

	return &userDAO{db: db}
}