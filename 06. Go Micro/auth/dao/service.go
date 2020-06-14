package dao

import (
	"auth/dao/user"
	"auth/model"
	"github.com/jinzhu/gorm"
)

type authDAOCreator struct {
	db *gorm.DB
}

func NewAuthDAOCreator(db *gorm.DB) *authDAOCreator {
	return &authDAOCreator{
		db: db,
	}
}

type UserDAOService interface {
	Insert(*model.Auth) (result *model.Auth, err error)
}

func (dc *authDAOCreator) GetDefaultAuthDAO() UserDAOService {
	tx := dc.db.New()
	return &user.DefaultDAO{
		DB: tx,
	}
}
