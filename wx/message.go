package wx

import (
	"encoding/xml"
)

//EncMessage 微信加密消息固定格式
type EncMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"-"` // 开发者微信号
	Encrypt      string   // 加密的消息报文
	MsgSignature string   // 报文签名
	TimeStamp    string   // 时间戳
	Nonce        string   // 随机数
}

//ReceiveMessage 微信每隔10分钟发过来一个数据ticket，需要解密
type ReceiveMessage struct {
	XMLName               xml.Name `xml:"xml"`
	AppID                 string   `xml:"AppId"`                 //第三方平台appid
	CreateTime            int64    `xml:"CreateTime"`            //时间戳
	InfoType              string   `xml:"InfoType"`              //component_verify_ticket
	ComponentVerifyTicket string   `xml:"ComponentVerifyTicket"` //Ticket内容
}

//APIComponentTokenRequest api
type APIComponentTokenRequest struct {
	ComponentAppid        string `json:"component_appid"`
	ComponentAppsecret    string `json:"component_appsecret"`
	ComponentVerifyTicket string `json:"component_verify_ticket"`
}

//APIComponentTokenResponse api
type APIComponentTokenResponse struct {
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int    `json:"expires_in"`
}

//PreAuthCode pac
type PreAuthCode struct {
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int    `json:"expires_in"`
}
