package main

import (
	"github.com/jinzhu/gorm"
	"github.com/stevenkitter/wxapi/models"
)

//Migrate migrate
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.WXTicket{})
}
