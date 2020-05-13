package dataservice

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func ConnectDB(pwd, db string) (*gorm.DB, error) {
	conn, err := gorm.Open("mysql", fmt.Sprintf("root:%s@tcp(localhost)/%s?charset=utf8&parseTime=True&loc=Local", pwd, db))
	if err != nil { return nil, err }

	return conn, nil
}