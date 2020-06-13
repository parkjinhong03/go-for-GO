package dao

import (
	"github.com/jinzhu/gorm"
)

type daoCreator struct {
	db *gorm.DB
}

func NewDaoCreator(db *gorm.DB) *daoCreator {
	return &daoCreator{
		db: db,
	}
}
