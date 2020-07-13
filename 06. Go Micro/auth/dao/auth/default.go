package user

import (
	"auth/model"
	"auth/tool/hash"
	"auth/tool/parser"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type defaultDAO struct {
	db *gorm.DB
}

func NewDefaultDAO(db *gorm.DB) *defaultDAO {
	return &defaultDAO{
		db: db,
	}
}

func (d *defaultDAO) InsertAuth(u *model.Auth) (result *model.Auth, err error) {
	if !contains([]string{CreatePending, Created, Rejected, Remove}, u.Status) {
		err = InvalidStatusError
		return
	}

	if u.UserPw, err = hash.BcryptGenerate(u.UserPw, bcrypt.DefaultCost); err != nil {
		err = BcryptGenerateError
		return
	}

	r := d.db.Create(u)
	if r.Error == nil {
		result = r.Value.(*model.Auth)
		return
	}

	var code int
	var msg string
	if me, ok := r.Error.(*mysql.MySQLError); ok {
		code = int(me.Number)
		msg = me.Message
	}

	//code, err := parser.DBErrorParse(r.Error.Error())
	//if err != nil {
	//	err = parser.InvalidError
	//	return
	//}

	switch code {
	case DuplicateErrorCode:
		switch attr := strings.Split(msg, "'")[3]; attr {
		case KeyUserId:
			err = UserIdDuplicatedError
		default:
			err = r.Error
		}
	case DataTooLongErrorCode:
		switch attr := strings.Split(msg, "'")[1]; attr {
		case ColumnUserId:
			err = UserIdTooLongError
		default:
			err = r.Error
		}
	default:
		err = r.Error
	}
	return
}

func (d *defaultDAO) UpdateStatus(id uint, status string) error {
	if !contains([]string{CreatePending, Created, Rejected, Remove}, status) {
		return InvalidStatusError
	}

	var auth model.Auth
	var count uint
	auth.ID = id

	d.db.Where("id = ?", id).Find(&auth).Count(&count)
	if count == 0 {
		return NonexistentUserError
	}

	r := d.db.Model(&auth).Update("status", status)
	return r.Error
}

func (d *defaultDAO) InsertMessage(m *model.ProcessedMessage) (result *model.ProcessedMessage, err error) {
	r := d.db.Create(m)
	if r.Error == nil { result = r.Value.(*model.ProcessedMessage); return }

	code, err := parser.DBErrorParse(r.Error.Error())
	if err != nil { err = parser.InvalidError; return }

	switch code {
	case DuplicateErrorCode:
		err = MsgIdDuplicateError
	case DataTooLongErrorCode:
		//err = StatusTooLongError
	default:
		err = r.Error
	}
	return
}

func (d *defaultDAO) CheckIfUserIdExist(id string) (exist bool, err error) {
	auth := new(model.Auth)
	result := d.db.Where("user_id = ?", id).Find(auth)
	if err = result.Error; err != nil { return }
	if result.RowsAffected == 0 { exist = false } else { exist = true }
	return
}

func (d *defaultDAO) Commit() *gorm.DB {
	return d.db.Commit()
}

func (d *defaultDAO) Rollback() *gorm.DB {
	return d.db.Rollback()
}

func contains(arrStr []string, index string) bool {
	for _, str := range arrStr {
		if str == index { return true }
	}
	return false
}