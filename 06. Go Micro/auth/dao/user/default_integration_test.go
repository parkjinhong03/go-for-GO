package user

import (
	"auth/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

var ud *defaultDAO

func connectDB () *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	addr := fmt.Sprintf("%s:%s@/UserTestDB?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd)
	db, err := gorm.Open("mysql", addr)
	if err != nil { log.Fatal(err) }
	return db
}

func init() {
	db := connectDB()

	db.DropTableIfExists(model.Auth{}, model.ProcessedMessage{})
	db.CreateTable(model.Auth{}, model.ProcessedMessage{})
	db.AutoMigrate(model.Auth{}, model.ProcessedMessage{})
	db.LogMode(false)
	_ = db.Close()
}

func setUpEnv() {
	db := connectDB()

	db.LogMode(false)
	ud = NewDefaultDAO(db)
	ud.db = ud.db.Begin()
}

