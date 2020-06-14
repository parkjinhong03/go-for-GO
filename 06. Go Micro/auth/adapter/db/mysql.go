package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func ConnMysql() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPwd := os.Getenv("DB_PASSWORD")
	dialect := fmt.Sprintf("%s:%s@/AuthDB?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd)

	return gorm.Open("mysql", dialect)
}