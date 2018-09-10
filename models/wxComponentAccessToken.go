package models

import (
	"github.com/jinzhu/gorm"
)

//WXComponentAccessToken token
type WXComponentAccessToken struct {
	gorm.Model
	ComponentAppid       string `gorm:"unique_index;not null"`
	ComponentAccessToken string
	ExpiresIn            int64
}
