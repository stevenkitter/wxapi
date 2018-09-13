package wx

import (
	"crypto/sha1"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	//Token token
	Token = "K81ec6fF37"
	//AppID appid
	AppID = "wxdd9779d0ca45ea77"
	//EncodingAESKey encodingaeskey
	EncodingAESKey = "7WfXuJfsGHYqt5eSPH8Gg7B9Y115vU8dx4Z48rZbzH1"
)

var (
	//AppSecrect wx app secrect
	AppSecrect = os.Getenv("WXAppSecrect")
)

//CheckSignature 验证签名
func CheckSignature(timestamp, nonce, encrypt, sign string) bool {
	tmpArr := []string{Token, timestamp, nonce, encrypt}
	sort.Strings(tmpArr)
	tmpStr := strings.Join(tmpArr, "")
	actual := fmt.Sprintf("%x", sha1.Sum([]byte(tmpStr)))
	return actual == sign
}
