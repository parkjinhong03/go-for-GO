package dataservice

import (
	"MSA.example.com/1/model"
	"errors"
	"github.com/jinzhu/gorm"
)

type userDAO struct {
	db *gorm.DB
}

func GetUserDAO(db *gorm.DB) *userDAO {
	db.LogMode(false)
	db.AutoMigrate(
		&model.Users{},
	)
	if !db.HasTable(&model.Users{}) {
		db.CreateTable(&model.Users{})
	}

	return &userDAO{db: db}
}

func (u *userDAO) Find(id uint32) (*model.Users, error) {
	user := model.Users{}
	if u.db.Where("id = ?", id).Find(&user).RowsAffected != 0 {
		return nil, errors.New("A user with that ID already exists")
	}
	return &user, nil
}

func (u *userDAO) FindByUserId(userId string) (*model.Users, error) {
	user := model.Users{}
	if u.db.Where("user_id = ?", userId).Find(&user).RowsAffected != 0 {
		return nil, errors.New("A user with that ID already exists")
	}
	return &user, nil
}

func (u *userDAO) Insert(user *model.Users) (*model.Users, error) {
	var r *model.Users
	txFunc := func(tx *gorm.DB) error {
		if tx = tx.Create(user); tx.Error == nil {
			r = tx.Value.(*model.Users)
		}
		return tx.Error
	}
	return r, u.db.Transaction(txFunc)
}


func (u *userDAO) Remove(id uint32) (int64, error) {
	db := u.db.Where("id = ?", id).Delete(model.Users{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *userDAO) UpdateStatus(user *model.Users, status string) (*model.Users, error) {
	var r *model.Users
	txFunc := func(tx *gorm.DB) error {
		if tx = tx.Model(user).Update("status", status); tx.Error == nil {
			r = tx.Value.(*model.Users)
		}
		return tx.Error
	}
	return r, u.db.Transaction(txFunc)
}