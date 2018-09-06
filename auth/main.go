package main

import (
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	micro "github.com/micro/go-micro"
	"github.com/stevenkitter/wxapi/db"
	"github.com/stevenkitter/wxapi/global"
	proto "github.com/stevenkitter/wxapi/proto/wxAuth"
)

const (
	serverName = "wx.srv.auth"
)

func main() {
	//mysql
	mysql, err := db.InitMysql()
	Migrate(mysql) //migrate
	if err != nil {
		log.Panic(err)
	}
	defer mysql.Close()
	service := micro.NewService(
		micro.Name(serverName),
		micro.Version(global.Version),
	)
	service.Init()
	wxAuth := &WxAuth{
		db: mysql,
	}
	proto.RegisterWxAuthHandler(service.Server(), wxAuth)
	log.Println("Run the server")
	if err := service.Run(); err != nil {
		log.Panicln(err)
	}
}
