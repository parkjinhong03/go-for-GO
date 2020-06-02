package model

import "github.com/jinzhu/gorm"

type UserInform struct {
	gorm.Model
	UserPk		 uint   `gorm:"Type:INT(10) UNSIGNED;"`
	Name         string `gorm:"NOT NULL;Type:VARCHAR(20)"`
	PhoneNumber  string `gorm:"NOT NULL;Type:CHAR(11)"`
	Email        string `gorm:"NOT NULL;Type:VARCHAR(30)"`
	Introduction string `gorm:"Type:VARCHAR(100)"`
	NumOfBlog    int    `gorm:"NOT NULL;Type:SMALLINT UNSIGNED;DEFAULT:0"`
}