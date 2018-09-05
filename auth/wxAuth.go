package main

import (
	"context"
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/stevenkitter/wxapi/models"
	proto "github.com/stevenkitter/wxapi/proto/wxAuth"
)

//WxAuth 微信授权模块
type WxAuth struct {
	db *gorm.DB
}

//SaveTicket 实现
func (wx *WxAuth) SaveTicket(ctx context.Context, req *proto.WxAuthTicketSaveRequest, rsp *proto.WxAuthTicketSaveResponse) error {
	if govalidator.IsNull(req.AppID) {
		return errors.New("appid can not be null")
	}
	if req.CreateTime == 0 {
		return errors.New("createTime can not be null")
	}
	if govalidator.IsNull(req.InfoType) {
		return errors.New("infoType can not be null")
	}
	if govalidator.IsNull(req.ComponentVerifyTicket) {
		return errors.New("componentVerifyTicket can not be null")
	}

	res := wx.db.Where(models.WXTicket{AppID: req.AppID}).
		Assign(models.WXTicket{
			AppID:                 req.AppID,
			CreateTime:            req.CreateTime,
			InfoType:              req.InfoType,
			ComponentVerifyTicket: req.ComponentVerifyTicket,
		}).
		FirstOrCreate(&models.WXTicket{})
	if res.Error == nil {
		rsp.Message = "ticket saved ok"
	}
	return res.Error
}

//GetTicket 实现
func (wx *WxAuth) GetTicket(ctx context.Context, req *proto.WxAuthTicketGetRequest, rsp *proto.WxAuthTicketGetResponse) error {
	if govalidator.IsNull(req.AppID) {
		return errors.New("appid can not be null")
	}
	var ticket models.WXTicket
	res := wx.db.Find(&ticket, "app_id = ?", req.AppID)
	if res.Error == nil {
		rsp.ComponentVerifyTicket = ticket.ComponentVerifyTicket
	}
	return res.Error
}
