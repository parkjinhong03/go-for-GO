package dataservice

import (
	"MSA.example.com/1/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserDataService interface {
	Find(id uint) (user *model.Users, exist bool)
	FindByUserId(userId string) (user *model.Users, exist bool)
	Insert(user *model.Users) (result *model.Users, err error)
	Remove(id uint) (rowAffected int64, err error)
	UpdateStatus(id uint, status string) (result *model.Users, err error)
}

type UserInformDataService interface {
	Insert(userInform *model.UserInform) (result *model.UserInform, err error)
}