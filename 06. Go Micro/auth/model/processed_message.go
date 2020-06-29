package model

import "github.com/jinzhu/gorm"

type ProcessedMessage struct {
	gorm.Model
	MsgId string `gorm:"Type:char(32);NOT NULL;UNIQUE_INDEX"`
}
