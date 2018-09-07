package main

import (
	"log"

	"github.com/gin-gonic/gin"
	wxauth "github.com/stevenkitter/wxapi/proto/wxAuth"
)

var (
	cl wxauth.WxAuthService
)

func main() {
	// setup Greeter Server Client
	cl = wxauth.NewWxAuthService("wx.srv.auth", nil)

	wxau := new(AuthWX)
	router := gin.Default()
	router.GET("/wx", wxau.Receive)
	router.GET("/ticket/:appId", wxau.GetTicket)
	// Register Handler
	if err := router.Run(":8080"); err != nil {
		log.Println(err)
	}
}
