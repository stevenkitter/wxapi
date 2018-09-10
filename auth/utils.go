package main

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stevenkitter/wxapi/helper"
	"github.com/stevenkitter/wxapi/models"
	"github.com/stevenkitter/wxapi/wx"
)

//Ticket mysql
func Ticket(db *gorm.DB, appID string) (string, error) {
	var ticket models.WXTicket
	res := db.Find(&ticket, "app_id = ?", appID)
	if res.Error != nil {
		return "", res.Error
	}
	return ticket.ComponentVerifyTicket, nil
}

//RequestToken req
func RequestToken(data map[string]interface{}) (string, float64, error) {
	result, err := helper.WXPost("/api_component_token", data)
	if err != nil {
		return "", 0, err
	}
	if result["errcode"] != nil {
		return "", 0, errors.New(result["errmsg"].(string))
	}
	return result["component_access_token"].(string), result["expires_in"].(float64), nil
}

//Token token
func Token(db *gorm.DB, appID string) (string, error) {
	var token models.WXComponentAccessToken
	res := db.Where("component_appid = ?", appID).Find(&token)
	if res.Error != nil && !res.RecordNotFound() {
		return "", res.Error
	}
	if time.Now().Unix() >= token.ExpiresIn {
		ticket, err := Ticket(db, appID)
		if err != nil {
			return "", err
		}
		data := map[string]interface{}{
			"component_appid":         wx.AppID,
			"component_appsecret":     wx.AppSecrect,
			"component_verify_ticket": ticket,
		}
		tokenStr, expiresIn, err := RequestToken(data)
		if err != nil {
			return "", err
		}
		mysqlDB := db.Where(models.WXComponentAccessToken{ComponentAppid: wx.AppID}).
			Assign(models.WXComponentAccessToken{
				ComponentAppid:       wx.AppID,
				ComponentAccessToken: tokenStr,
				ExpiresIn:            time.Now().Unix() + int64(expiresIn),
			}).
			FirstOrCreate(&models.WXComponentAccessToken{})
		if mysqlDB.Error != nil {
			return "", mysqlDB.Error
		}
		return tokenStr, nil
	}

	return token.ComponentAccessToken, nil
}

//RequestPreAuthCode req
func RequestPreAuthCode(db *gorm.DB, appID string) (string, float64, error) {
	token, err := Token(db, appID)
	if err != nil {
		return "", 0, err
	}
	data := map[string]interface{}{
		"component_appid": appID,
	}
	res, err := helper.WXPost("/api_create_preauthcode?component_access_token="+token, data)

	if err != nil {
		return "", 0, err
	}
	if res["errcode"] != nil {
		return "", 0, errors.New(res["errmsg"].(string))
	}
	return res["pre_auth_code"].(string), res["expires_in"].(float64), nil
}

//PreAuthCode pre_auth_code
func PreAuthCode(db *gorm.DB, appID string) (string, error) {
	var code models.WXPreAuthCode
	res := db.Where("component_appid = ?", appID).Find(&code)
	if res.Error != nil && !res.RecordNotFound() {
		return "", res.Error
	}
	if time.Now().Unix() >= code.ExpiresIn {
		preCode, expiresIn, err := RequestPreAuthCode(db, appID)
		if err != nil {
			return "", err
		}
		mysqlDB := db.Where(models.WXPreAuthCode{ComponentAppid: appID}).
			Assign(models.WXPreAuthCode{
				ComponentAppid: appID,
				PreAuthCode:    preCode,
				ExpiresIn:      time.Now().Unix() + int64(expiresIn),
			}).
			FirstOrCreate(&models.WXPreAuthCode{})
		if mysqlDB.Error != nil {
			return "", mysqlDB.Error
		}
		return preCode, nil
	}
	return code.PreAuthCode, nil
}
