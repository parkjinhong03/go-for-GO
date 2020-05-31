package model

import (
	"github.com/jinzhu/gorm"
)

type Users struct {
	gorm.Model
	UserId string `gorm:"NOT NULL;Type:varchar(100);unique_index"`
	UserPwd string `gorm:"NOT NULL;Type:TEXT"`
	Status string `gorm:"NOT NULL;Type:varchar(20)"`
}