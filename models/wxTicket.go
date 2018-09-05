package models

import "github.com/jinzhu/gorm"

//WXTicket mysql 表结构
type WXTicket struct {
	gorm.Model
	AppID                 string `gorm:"unique_index;not null"`
	CreateTime            int64
	InfoType              string
	ComponentVerifyTicket string
}
