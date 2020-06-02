package dataservice

import (
	"MSA.example.com/1/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserDataService interface {
	Find(id uint32) (user *model.Users, err error)
	FindByUserId(userId string) (user *model.Users, err error)
	Insert(user *model.Users) (result *model.Users, err error)
	Remove(id uint32) (rowAffected int64, err error)
	UpdateStatus(user *model.Users, status string) (result *model.Users, err error)
}

type UserInformDataService interface {
	Insert(userInform *model.UserInform) (result *model.UserInform, err error)
}