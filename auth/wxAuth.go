package main

import (
	"context"
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/stevenkitter/wxapi/helper"
	"github.com/stevenkitter/wxapi/models"

	proto "github.com/stevenkitter/wxapi/proto/wxAuth"
	wxPakage "github.com/stevenkitter/wxapi/wx"
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
func (wx *WxAuth) GetTicket(ctx context.Context, req *proto.ComponentAppidRequest, rsp *proto.WxAuthTicketGetResponse) error {
	if govalidator.IsNull(req.ComponentAppid) {
		return errors.New("appid can not be null")
	}
	ticket, err := Ticket(wx.db, req.ComponentAppid)
	if err == nil {
		rsp.ComponentVerifyTicket = ticket
	}
	return err
}

//GetComponentAccessToken req token
func (wx *WxAuth) GetComponentAccessToken(ctx context.Context, req *proto.ComponentAccessTokenRequest, rsp *proto.ComponentAccessTokenResponse) error {
	if govalidator.IsNull(req.ComponentAppid) {
		return errors.New("componentAppid can not be null")
	}
	token, err := Token(wx.db, req.ComponentAppid)
	if err != nil {
		return err
	}
	rsp.ComponentAccessToken = token
	return nil
}

//GetPreAuthCode req pre_auth_code
func (wx *WxAuth) GetPreAuthCode(ctx context.Context, req *proto.ComponentAppidRequest, rsp *proto.PreAuthCodeResponse) error {
	if govalidator.IsNull(req.ComponentAppid) {
		return errors.New("componentAppid can not be null")
	}
	preCode, err := PreAuthCode(wx.db, req.ComponentAppid)
	if err != nil {
		return err
	}
	rsp.PreAuthCode = preCode
	return nil
}

//GetAuthorizationInfo authorizationInfo
func (wx *WxAuth) GetAuthorizationInfo(ctx context.Context, req *proto.AuthorizationInfoRequest, rsp *proto.AuthorizationInfoResponse) error {
	return nil
}

//RefreshAuthorizerAccessToken refresh token
func (wx *WxAuth) RefreshAuthorizerAccessToken(ctx context.Context, req *proto.RefreshAuthorizerAccessTokenRequest, rsp *proto.RefreshAuthorizerAccessTokenResponse) error {
	return nil
}

//GetAuthorizerInfo info
func (wx *WxAuth) GetAuthorizerInfo(ctx context.Context, req *proto.AuthorizerInfoRequest, rsp *proto.AuthorizerInfoResponse) error {
	return nil
}

//GetAuthorizerOption option
func (wx *WxAuth) GetAuthorizerOption(ctx context.Context, req *proto.AuthorizerOptionRequest, rsp *proto.AuthorizerOptionResponse) error {
	return nil
}

//SetAuthorizerOption set option
func (wx *WxAuth) SetAuthorizerOption(ctx context.Context, req *proto.SetAuthorizerOptionRequest, rsp *proto.SetAuthorizerOptionResponse) error {
	return nil
}

//GetAuthURL auth url
func (wx *WxAuth) GetAuthURL(ctx context.Context, req *proto.GetAuthURLRequest, rsp *proto.GetAuthURLResponse) error {
	if govalidator.IsNull(req.ComponentAppid) {
		return errors.New("componentAppid can not be null")
	}
	if govalidator.IsNull(req.RedirectURL) {
		return errors.New("redirectURL can not be null")
	}
	if !govalidator.IsNull(req.BizAppid) && !govalidator.IsNull(req.AuthType) {
		return errors.New("bizAppid and authType can not be set the same time")
	}
	preAuthCode, err := PreAuthCode(wx.db, req.ComponentAppid)
	if err != nil {
		return err
	}
	var params = map[string]interface{}{
		"component_appid": req.ComponentAppid,
		"pre_auth_code":   preAuthCode,
		"redirect_uri":    req.RedirectURL,
	}
	if req.Tag == 0 {
		//web
		if !govalidator.IsNull(req.AuthType) {
			params["auth_type"] = req.AuthType
		} else {
			params["biz_appid"] = req.BizAppid
		}

		url := wxPakage.WXAuthURL + "?" + helper.JoinParams(params)
		rsp.AuthURL = url
		return nil
	}
	//app
	params["action"] = "bindcomponent"
	params["no_scan"] = "1"

	if !govalidator.IsNull(req.AuthType) {
		params["auth_type"] = req.AuthType
	} else {
		params["biz_appid"] = req.BizAppid
	}
	url := wxPakage.WXAPPAuthURL + "?" + helper.JoinParams(params) + "#wechat_redirect"
	rsp.AuthURL = url

	return nil
}
