package dao

import (
	"auth/dao/user"
	"auth/model"
	"github.com/jinzhu/gorm"
)

type AuthDAOCreator struct {
	db *gorm.DB
}

func NewAuthDAOCreator(db *gorm.DB) *AuthDAOCreator {
	if !db.HasTable(&model.Auth{}) {
		db.CreateTable(&model.Auth{})
	}
	db.AutoMigrate(&model.Auth{})
	return &AuthDAOCreator{
		db: db,
	}
}

type UserDAOService interface {
	Insert(*model.Auth) (result *model.Auth, err error)
}

func (dc *AuthDAOCreator) GetDefaultAuthDAO() UserDAOService {
	tx := dc.db.Begin()
	return &user.DefaultDAO{
		DB: tx,
	}
}
