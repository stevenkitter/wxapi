package main

import (
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	micro "github.com/micro/go-micro"
	logger "github.com/sirupsen/logrus"
	"github.com/stevenkitter/wxapi/db"
	"github.com/stevenkitter/wxapi/global"
	"github.com/stevenkitter/wxapi/log"
	"github.com/stevenkitter/wxapi/micro/trace"
	proto "github.com/stevenkitter/wxapi/proto/wxAuth"
	"github.com/stevenkitter/wxapi/tracing"
)

const (
	serverName = "wx.srv.auth"
)

var (
	serverLogger *logger.Entry
)

func main() {
	//mysql
	mysql, err := db.InitMysql()
	Migrate(mysql) //migrate
	if err != nil {
		serverLogger.Panic(err)
	}
	defer mysql.Close()
	//logger
	serverLogger = log.InitLogger(serverName)
	//tracing
	tracer, closer, err := tracing.InitTracing(serverName)
	defer closer.Close()
	if err != nil {
		serverLogger.Panic(err)
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

	serverLogger.Info("Run the server")
	if err := service.Run(); err != nil {
		serverLogger.Panic(err)
	}
}

func init() {
	logger.SetFormatter(&logger.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logger.InfoLevel)
}
