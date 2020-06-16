package model

import "github.com/jinzhu/gorm"

type Auth struct {
	gorm.Model
	UserId string `gorm:"Type:varchar(20);NOT NULL;UNIQUE_INDEX"`
	UserPw string `gorm:"Type:varchar(100);NOT NULL"`
	Status string `gorm:"Type:varchar(20);NOT NULL"`
}
