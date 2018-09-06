package main

import (
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	micro "github.com/micro/go-micro"
	"github.com/stevenkitter/wxapi/db"
	"github.com/stevenkitter/wxapi/global"
	"github.com/stevenkitter/wxapi/micro/trace"
	proto "github.com/stevenkitter/wxapi/proto/wxAuth"
	"github.com/stevenkitter/wxapi/tracing"
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
	//tracing
	tracer, closer, err := tracing.InitTracing(serverName)
	defer closer.Close()
	if err != nil {
		log.Panic(err)
	}
	service := micro.NewService(
		micro.Name(serverName),
		micro.Version(global.Version),
		micro.WrapHandler(trace.NewHandlerWrapper(tracer)),
		micro.WrapClient(trace.NewClientWrapper(tracer)),
		micro.WrapSubscriber(trace.NewSubscriberWrapper(tracer)),
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
