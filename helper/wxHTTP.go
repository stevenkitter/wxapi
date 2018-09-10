package helper

import (
	"github.com/stevenkitter/wxapi/wx"
)

//WXGet wx get
func WXGet(url string, params map[string]interface{}) (map[string]interface{}, error) {
	res, err := Get(wx.WXBaseURL+url, params)
	if err != nil {
		return nil, err
	}
	var v map[string]interface{}
	err = ResponseToInterface(res, &v)
	return v, err
}

//WXPost wx post
func WXPost(url string, data map[string]interface{}) (map[string]interface{}, error) {
	res, err := PostJSON(wx.WXBaseURL+url, data)
	if err != nil {
		return nil, err
	}
	var v map[string]interface{}
	err = ResponseToInterface(res, &v)
	return v, err
}
