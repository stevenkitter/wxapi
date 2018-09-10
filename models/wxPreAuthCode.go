package models

import (
	"github.com/jinzhu/gorm"
)

//WXPreAuthCode pre auth code
type WXPreAuthCode struct {
	gorm.Model
	ComponentAppid string `gorm:"unique_index;not null"`
	PreAuthCode    string
	ExpiresIn      int64
}
