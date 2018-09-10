package main

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"log"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	wxauth "github.com/stevenkitter/wxapi/proto/wxAuth"
	"github.com/stevenkitter/wxapi/response"
	"github.com/stevenkitter/wxapi/wx"
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
	saveResult, err := cl.SaveTicket(context.TODO(), &wxauth.WxAuthTicketSaveRequest{
		AppID:                 receivedMessage.AppID,
		CreateTime:            receivedMessage.CreateTime,
		InfoType:              receivedMessage.InfoType,
		ComponentVerifyTicket: receivedMessage.ComponentVerifyTicket,
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println(saveResult.Message)
	}

	c.String(200, "success")
}

//GetTicket GetTicket
func (auth *AuthWX) GetTicket(c *gin.Context) {
	appID := c.Param("appId")
	res, err := cl.GetTicket(context.TODO(), &wxauth.ComponentAppidRequest{
		ComponentAppid: appID,
	})
	if err != nil {
		c.JSON(400, response.ResFAIL(err.Error()))
		return
	}
	c.JSON(200, response.Response{
		Code: 200,
		Data: res.ComponentVerifyTicket,
	})
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

//GetAuthURLJSON json
type GetAuthURLJSON struct {
	RedirectURL string `json:"redirectURL" binding:"required"`
	AuthType    string `json:"authType"`
	BizAppid    string `json:"bizAppid"`
	Tag         int64  `json:"tag"`
}

//GetAuthURL auth url
func (auth *AuthWX) GetAuthURL(c *gin.Context) {
	var json GetAuthURLJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(400, response.ResFAIL(err.Error()))
		return
	}
	if !govalidator.IsNull(json.BizAppid) && !govalidator.IsNull(json.BizAppid) {
		c.JSON(400, response.ResFAIL("bizAppid and authType 互斥"))
		return
	}
	res, err := cl.GetAuthURL(context.TODO(), &wxauth.GetAuthURLRequest{
		ComponentAppid: wx.AppID,
		RedirectURL:    json.RedirectURL,
		Tag:            json.Tag,
		AuthType:       json.AuthType,
		BizAppid:       json.BizAppid,
	})
	if err != nil {
		c.JSON(400, response.ResFAIL(err.Error()))
		return
	}
	c.JSON(200, response.Response{
		Code: 200,
		Data: res.AuthURL,
	})
}
