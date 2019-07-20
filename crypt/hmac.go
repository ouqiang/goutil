package crypt

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

// HMacMD5 生成hmac hash
func HMacMD5(s string, key []byte) string {
	mac := hmac.New(md5.New, key)
	_, _ = mac.Write([]byte(s))

	return hex.EncodeToString(mac.Sum(nil))
}
