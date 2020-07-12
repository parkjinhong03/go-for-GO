package user

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

	ud = NewDefaultDAO(db)
}
