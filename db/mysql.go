package db

import (
	"os"

	"github.com/jinzhu/gorm"
)

//InitMysql 初始化mysql
func InitMysql() (*gorm.DB, error) {
	mysqlURL := os.Getenv("MysqlURL")
	if mysqlURL == "" {
		mysqlURL = "115.159.222.199:3306"
	}
	return gorm.Open("mysql", "wx:rWk1hvqMT62K2JYH@tcp("+mysqlURL+")/wx?charset=utf8&parseTime=True&loc=Local")
}
