package user

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"strings"
	"user/model"
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

	var code int
	var msg string
	if me, ok := r.Error.(*mysql.MySQLError); ok { // 만약 mysql 연결 db가 아닐 경우에는?
		code = int(me.Number)
		msg = me.Message
	}

	//code, err = parser.ParseDBError(r.Error.Error())
	//if err != nil { err = parser.InvalidError; return }

	switch code {
	case DuplicatedErrorCode:
		switch attr := strings.Split(msg, "'")[3]; attr {
		case KeyAuthId:
			err = AuthIdDuplicatedError
		case KeyEmail:
			err = EmailDuplicatedError
		default:
			err = r.Error
		}
	case DataTooLongErrorCode:
		switch attr := strings.Split(msg, "'")[1]; attr {
		case ColumnName:
			err = NameTooLongError
		case ColumnEmail:
			err = EmailTooLongError
		case ColumnPhoneNumber:
			err = PhoneNumberTooLongError
		default:
			err = r.Error
		}
	default:
		err = r.Error
	}

	return
}

func (d *defaultDAO) CheckIfEmailExist(email string) (exist bool, err error) {
	user := model.User{}
	exist = false
	r := d.db.Where("email = ?", email).Find(&user)

	if r.Error != nil && r.Error != gorm.ErrRecordNotFound {
		err = r.Error
	}
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

	var code int
	var msg string
	if me, ok := r.Error.(*mysql.MySQLError); ok {
		code = int(me.Number)
		msg = me.Message
	}

	//code, err := parser.ParseDBError(r.Error.Error())
	//if err != nil { err = parser.InvalidError; return }

	switch code {
	case DuplicatedErrorCode:
		switch attr := strings.Split(msg, "'")[3]; attr {
		case KeyMsgId:
			err = MessageDuplicatedError
		default:
			err = r.Error
		}
	case DataTooLongErrorCode:
		switch attr := strings.Split(msg, "'")[1]; attr {
		case ColumnMsgId:
			err = MessageTooLongError
		default:
			err = r.Error
		}
	default:
		err = r.Error
	}

	return
}