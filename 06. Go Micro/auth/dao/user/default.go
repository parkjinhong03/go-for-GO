package user

import "github.com/jinzhu/gorm"

type DefaultDAO struct {
	*gorm.DB
}