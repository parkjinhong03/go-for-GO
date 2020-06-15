package user

import (
	"auth/model"
	"auth/tool/parser"
	"github.com/jinzhu/gorm"
)

type defaultDAO struct {
	db *gorm.DB
}

func NewDefaultDAO(db *gorm.DB) *defaultDAO {
	return &defaultDAO{
		db: db,
	}
}

func (d *defaultDAO) Insert(u *model.Auth) (result *model.Auth, err error) {
	u.Status = CreatePending
	r := d.db.Create(u)
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

func (d *defaultDAO) Commit() *gorm.DB {
	return d.db.Commit()
}

func (d *defaultDAO) Rollback() *gorm.DB {
	return d.db.Rollback()
}