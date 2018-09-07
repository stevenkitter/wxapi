package main

import (
	"log"

	"github.com/gin-gonic/gin"
	micro "github.com/micro/go-micro"
	wxauth "github.com/stevenkitter/wxapi/proto/wxAuth"
)

var (
	cl wxauth.WxAuthService
)

func main() {
	service := micro.NewService(micro.Name("wx.api.client"))
	service.Init()
	// setup Greeter Server Client
	cl = wxauth.NewWxAuthService("wx.srv.auth", service.Client())

	wxau := new(AuthWX)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/wx", wxau.Receive)
	router.GET("/ticket/:appId", wxau.GetTicket)
	// Register Handler
	if err := router.Run(":8080"); err != nil {
		log.Println(err)
	}
}
