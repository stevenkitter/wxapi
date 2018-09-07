package wx

// 微信使用的加密解密
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

//EncryptMsg 加密报文
func EncryptMsg(msg []byte, aesKey []byte, AppID string) (b64Enc string, err error) {
	// 拼接完整报文
	src := SpliceFullMsg(msg, AppID)

	// AES CBC 加密报文
	dst, err := AESCBCEncrypt(src, aesKey, aesKey[:aes.BlockSize])
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(dst), nil
}

//DecryptMsg 解密报文
func DecryptMsg(b64Enc string, aesKey []byte, AppID string) (msg []byte, err error) {
	enc, err := base64.StdEncoding.DecodeString(b64Enc)
	if err != nil {
		return nil, err
	}

	// AES CBC 解密报文
	src, err := AESCBCDecrypt(enc, aesKey, aesKey[:aes.BlockSize])
	if err != nil {
		return nil, err
	}

	_, _, msg, AppID2 := ParseFullMsg(src)
	if AppID2 != AppID {
		return nil, fmt.Errorf("expected AppID %s, but %s", AppID, AppID2)
	}

	return msg, nil
}

// SpliceFullMsg 拼接完整报文，
// AES加密的buf由16个字节的随机字符串、4个字节的msg_len(网络字节序)、msg和$AppIDAppID组成，
// 其中msg_len为msg的长度，$AppIDAppID为公众帐号的AppIDAppID
func SpliceFullMsg(msg []byte, AppID string) (fullMsg []byte) {
	// 16个字节的随机字符串
	randBytes := RandBytes(16)

	// 4个字节的msg_len(网络字节序)
	msgLen := len(msg)
	lenBytes := []byte{
		byte(msgLen >> 24 & 0xFF),
		byte(msgLen >> 16 & 0xFF),
		byte(msgLen >> 8 & 0xFF),
		byte(msgLen & 0xFF),
	}

	return bytes.Join([][]byte{randBytes, lenBytes, msg, []byte(AppID)}, nil)
}

// ParseFullMsg 从完整报文中解析出消息内容，
// AES加密的buf由16个字节的随机字符串、4个字节的msg_len(网络字节序)、msg和$AppIDAppID组成，
// 其中msg_len为msg的长度，$AppIDAppID为公众帐号的AppIDAppID
func ParseFullMsg(fullMsg []byte) (randBytes []byte, msgLen int, msg []byte, AppID string) {
	randBytes = fullMsg[:16]

	msgLen = (int(fullMsg[16]) << 24) |
		(int(fullMsg[17]) << 16) |
		(int(fullMsg[18]) << 8) |
		int(fullMsg[19])
	// log.Tracef("msgLen=[% x]=(%d %d %d %d)=%d", fullMsg[16:20], (int(fullMsg[16]) << 24),
	// 	(int(fullMsg[17]) << 16), (int(fullMsg[18]) << 8), int(fullMsg[19]), msgLen)

	msg = fullMsg[20 : 20+msgLen]

	AppID = string(fullMsg[20+msgLen:])

	return
}

// RandBytes 产生 size 个长度的随机字节
func RandBytes(size int) (r []byte) {
	r = make([]byte, size)
	_, err := rand.Read(r)
	if err != nil {
		// 忽略错误，不影响其他逻辑，仅仅打印日志
		log.Printf("rand read error: %s", err)
	}
	return r
}

// AESCBCEncrypt 采用 CBC 模式的 AES 加密
func AESCBCEncrypt(src, key, iv []byte) (enc []byte, err error) {
	log.Printf("src: %s", src)
	src = PKCS7Padding(src, len(key))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(src, src)
	enc = src

	log.Printf("enc: % x", enc)
	return enc, nil
}

// AESCBCDecrypt 采用 CBC 模式的 AES 解密
func AESCBCDecrypt(enc, key, iv []byte) (src []byte, err error) {
	log.Printf("enc: % x", enc)
	if len(enc) < len(key) {
		return nil, fmt.Errorf("the length of encrypted message too short: %d", len(enc))
	}
	if len(enc)&(len(key)-1) != 0 { // or len(enc)%len(key) != 0
		return nil, fmt.Errorf("encrypted message is not a multiple of the key size(%d), the length is %d", len(key), len(enc))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(enc, enc)
	src = PKCS7UnPadding(enc)

	log.Printf("src: %s", src)
	return src, nil
}

// PKCS7Padding PKCS#7填充，Buf需要被填充为K的整数倍，
// 在buf的尾部填充(K-N%K)个字节，每个字节的内容是(K- N%K)
func PKCS7Padding(src []byte, k int) (padded []byte) {
	padLen := k - len(src)%k
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padding...)
}

// PKCS7UnPadding 去掉PKCS#7填充，Buf需要被填充为K的整数倍，
// 在buf的尾部填充(K-N%K)个字节，每个字节的内容是(K- N%K)
func PKCS7UnPadding(src []byte) (padded []byte) {
	padLen := int(src[len(src)-1])
	return src[:len(src)-padLen]
}
