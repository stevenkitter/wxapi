package main

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/stevenkitter/wxapi/api/wx"
	wxauth "github.com/stevenkitter/wxapi/proto/wxAuth"
)

//AuthWX AuthWX
type AuthWX struct{}

//Receive weixin income ticket
func (auth *AuthWX) Receive(c *gin.Context) {
	res, err := WXInComeHandler(c)
	if err != nil {
		c.String(200, "success")
		return
	}
	receivedMessage := &wx.ReceiveMessage{}
	err = xml.Unmarshal(res, receivedMessage)
	if err != nil {
		c.String(200, "success")
		return
	}
	c.String(200, "success")
}

//GetTicket GetTicket
func (auth *AuthWX) GetTicket(c *gin.Context) {
	appID := c.Param("appId")
	res, err := cl.GetTicket(context.TODO(), &wxauth.WxAuthTicketGetRequest{
		AppID: appID,
	})
	if err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(200, res)
}

//WXInComeHandler WXInComeHandler
func WXInComeHandler(c *gin.Context) ([]byte, error) {
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	encryptType := c.Query("encrypt_type")
	msgSignature := c.Query("msg_signature")
	if encryptType != "aes" {
		return nil, errors.New("encryptType is wrong")
	}
	//解析过来的数据
	var encMessage wx.EncMessage
	err := c.ShouldBind(&encMessage)
	if err != nil {
		return nil, err
	}
	//判断签名正确与否
	if !wx.CheckSignature(timestamp, nonce, encMessage.Encrypt, msgSignature) {
		return nil, errors.New("checkSignature is wrong")
	}
	aeskey, _ := base64.StdEncoding.DecodeString(wx.EncodingAESKey + "=")
	res, err := wx.DecryptMsg(encMessage.Encrypt, aeskey, wx.AppID)
	if err != nil {
		return nil, err
	}
	return res, nil
	//获取需要的数据

}
