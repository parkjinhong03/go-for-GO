package dataservice

import (
	"MSA.example.com/1/model"
	"github.com/jinzhu/gorm"
)

type userInformDAO struct {
	db *gorm.DB
}

func NewUserInformDAO(db *gorm.DB) *userInformDAO {
	userInform := &model.UserInform{}
	db.LogMode(false)
	db.AutoMigrate(userInform)
	if !db.HasTable(userInform) {
		db.CreateTable(userInform)
	}

	return &userInformDAO{db: db}
}

func (ui *userInformDAO) Insert(userInform *model.UserInform) (*model.UserInform, error) {
	var r *model.UserInform
	txFunc := func(tx *gorm.DB) (err error) {
		if tx = tx.Create(userInform); tx.Error == nil {
			r = tx.Value.(*model.UserInform)
		}
		return tx.Error
	}
	return r, ui.db.Transaction(txFunc)
}