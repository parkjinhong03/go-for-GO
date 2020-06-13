package dao

import (
	"auth/dao/user"
	"auth/model"
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

type UserDAOService interface {
	Insert(*model.User) (result *model.User, err error)
}

func (dc *daoCreator) GetUserDefaultDAO() UserDAOService {
	tx := dc.db.New()
	return &user.DefaultDAO{
		DB: tx,
	}
}
