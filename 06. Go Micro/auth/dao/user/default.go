package user

import (
	"auth/model"
	"auth/tool/parser"
	"github.com/jinzhu/gorm"
)

type DefaultDAO struct {
	*gorm.DB
}

func (d *DefaultDAO) Insert(u *model.Auth) (result *model.Auth, err error) {
	u.Status = CreatePending
	r := d.Create(u)
	if r.Error == nil { result = r.Value.(*model.Auth); return }

	code, err := parser.DBErrorParse(r.Error.Error())
	if err != nil { err = parser.InvalidError; return }

	switch code {
	case IdDuplicateErrorCode:
		err = IdDuplicateError
	default:
		err = UnknownError
	}
	return
}