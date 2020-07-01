package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	AuthId uint `gorm:"Type:mediumint(10) UNSIGNED;NOT NULL;UNIQUE_INDEX"`
	Name string `gorm:"Type:char(4);NOT NULL"`
	PhoneNumber string `gorm:"char(11);NOT NULL"`
	Email string `gorm:"varchar(30);NOT NULL;UNIQUE_INDEX"`
	Introduction string `gorm:"varchar(100);NOT NULL"`
}