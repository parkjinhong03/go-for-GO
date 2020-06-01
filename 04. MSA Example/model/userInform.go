package model

import "github.com/jinzhu/gorm"

type userInform struct {
	gorm.Model
	Users Users `gorm:"foreignkey:UserPk"`
	Name string `gorm:"NOT NULL;Type:VARCHAR(20)"`
	PhoneNumber int `gorm:"NOT NULL;Type:INT(11)"`
	Email string `gorm:"NOT NULL;Type:VARCHAR(20)"`
	Introduction string `gorm:"VARCHAR(100)"`
	NumOfBlog int `gorm:"NOT NULL;Type:UNSIGNED SMALLINT;DEFAULT:0"`
}