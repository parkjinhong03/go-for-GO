package user

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"os"
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
	ud.db = ud.db.Begin()
}
