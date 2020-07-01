package user

import (
	"github.com/jinzhu/gorm"
	"user/model"
	"user/tool/parser"
)

// 통합 테스트 필요
type defaultDAO struct {
	db *gorm.DB
}

func NewDefaultDAO(db *gorm.DB) *defaultDAO {
	return &defaultDAO{
		db: db,
	}
}

func (d *defaultDAO) InsertUser(user *model.User) (result *model.User, err error) {
	r := d.db.Create(user)
	if r.Error == nil { return r.Value.(*model.User), nil }

	code, err := parser.ParseDBError(r.Error.Error())
	if err != nil { err = parser.InvalidError; return }

	switch code {
	case EmailDuplicatedErrorCode:
		err = EmailDuplicatedError
	default:
		err = r.Error
	}

	return
}

func (d *defaultDAO) CheckIfEmailExist(email string) (exist bool, err error) {
	user := model.User{}
	exist = false
	r := d.db.Where("email = ?", email).Select(&user)
	if err = r.Error; err != nil { return }
	if r.RowsAffected != 0 { exist = true }
	return
}

func (d *defaultDAO) Commit() *gorm.DB {
	return d.db.Commit()
}

func (d *defaultDAO) Rollback() *gorm.DB {
	return d.db.Rollback()
}

func (d *defaultDAO) InsertMessage(pm *model.ProcessedMessage) (result *model.ProcessedMessage, err error) {
	r := d.db.Create(pm)
	if r.Error == nil { return r.Value.(*model.ProcessedMessage), nil }

	code, err := parser.ParseDBError(r.Error.Error())
	if err != nil { err = parser.InvalidError; return }

	switch code {
	case MessageDuplicatedErrorCode:
		err = MessageDuplicatedError
	default:
		err = r.Error
	}

	return
}