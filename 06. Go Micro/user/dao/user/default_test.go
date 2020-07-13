package user

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"user/model"
)

var ud *defaultDAO

func init() {
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	addr := fmt.Sprintf("%s:%s@/UserTestDB?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd)
	db, err := gorm.Open("mysql", addr)
	if err != nil { log.Fatal(err) }

	db.DropTableIfExists(model.User{})
	db.CreateTable(model.User{})
	db.AutoMigrate(model.User{})
	db.LogMode(false)

	_ = db.Close()
}

type insertUserTest struct {
	AuthId       uint
	Name         string
	PhoneNumber  string
	Email        string
	Introduction string
	ExpectError	 error
}

func (i insertUserTest) Exec() (*model.User, error) {
	return ud.InsertUser(&model.User{
		AuthId:       i.AuthId,
		Name:         i.Name,
		PhoneNumber:  i.PhoneNumber,
		Email:        i.Email,
		Introduction: i.Introduction,
	})
}

type checkIfEmailExistTest struct {
	Email       string
	ExpectExist bool
	ExpectError error
}

func (c checkIfEmailExistTest) Exec() (bool, error) {
	return ud.CheckIfEmailExist(c.Email)
}

func setUpEnv() {
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	addr := fmt.Sprintf("%s:%s@/UserTestDB?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd)
	db, err := gorm.Open("mysql", addr)
	if err != nil { log.Fatal(err) }
	db.LogMode(false)
	ud = NewDefaultDAO(db)
}

func TestDefaultUserDAOInsertUser(t *testing.T) {
	setUpEnv()
	ud.db = ud.db.Begin()

	tests := []insertUserTest{
		{
			AuthId:       1,
			Name:         "박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719@naver.com",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  nil,
		}, {
			AuthId:       1,
			Name:         "박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719@hanmail.net",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  AuthIdDuplicatedError,
		}, {
			AuthId:       2,
			Name:         "박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719@naver.com",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  EmailDuplicatedError,
		}, {
			AuthId:       3,
			Name:         "박진홍박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719@hanmail.net",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  NameTooLongError,
		}, {
			AuthId:       4,
			Name:         "박진홍",
			PhoneNumber:  "01088378347123",
			Email:        "jinhong0719@hanmail.net",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  PhoneNumberTooLongError,
		}, {
			AuthId:       5,
			Name:         "박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719jinhong0719@hanmail.net",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  EmailTooLongError,
		},
	}

	for _, test := range tests {
		_, err := test.Exec()
		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v\n)", test)
	}

	ud.Rollback()
	_ = ud.db.Close()
}

func TestDefaultUserDAOCheckIfEmailExist(t *testing.T) {
	setUpEnv()
	ud.db = ud.db.Begin()

	inits := []insertUserTest{
		{
			AuthId:       1,
			Name:         "박진홍",
			PhoneNumber:  "01088378347",
			Email:        "jinhong0719@naver.com",
			Introduction: "안녕하세요 저의 이름은 박진홍입니다.",
			ExpectError:  nil,
		},
	}

	for _, init := range inits {
		_, err := init.Exec()
		if err != nil { log.Fatal(err) }
	}

	tests := []checkIfEmailExistTest{
		{
			Email: "jinhong0719@naver.com",
			ExpectExist: true,
			ExpectError: nil,
		}, {
			Email: "jinhong0719@gmail.com",
			ExpectExist: false,
			ExpectError: nil,
		},
	}

	for _, test := range tests {
		exist, err := test.Exec()
		assert.Equalf(t, test.ExpectExist, exist, "exist assertion error (test case: %v)\n", test)
		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)\n", test)
	}

	ud.Rollback()
	_ = ud.db.Close()
}