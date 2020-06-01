package dataservice

import "github.com/jinzhu/gorm"

type userInformDAO struct {
	db *gorm.DB
}

func NewUserInformDAO(db *gorm.DB) *userInformDAO {
	return &userInformDAO{
		db: db,
	}
}