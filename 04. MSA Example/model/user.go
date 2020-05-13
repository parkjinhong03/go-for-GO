package model

import (
	"github.com/jinzhu/gorm"
)

type Users struct {
	gorm.Model
	UserId string `gorm:"NOT NULL;Type:varchar(100)"`
	UserPwd string `gorm:"NOT NULL;Type:TEXT"`
}