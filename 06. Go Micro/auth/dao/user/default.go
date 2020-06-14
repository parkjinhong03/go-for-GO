package user

import (
	"auth/model"
	"github.com/jinzhu/gorm"
)

type DefaultDAO struct {
	*gorm.DB
}

func (d *DefaultDAO) Insert(u *model.Auth) (*model.Auth, error) {
	r := d.Create(u)
	if r.Error != nil {
		return nil, r.Error
	}
	return r.Value.(*model.Auth), nil
}