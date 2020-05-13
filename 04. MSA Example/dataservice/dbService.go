package dataservice

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func ConnectDB(db string) (*gorm.DB, error) {
	pwd := os.Getenv("DB_PASSWORD")
	if pwd == "" {
		return nil, errors.New("unable to parse DB_PASSWORD from environment variable")
	}
	conn, err := gorm.Open("mysql", fmt.Sprintf("root:%s@tcp(localhost)/%s?charset=utf8&parseTime=True&loc=Local", pwd, db))
	if err != nil { return nil, err }

	return conn, nil
}