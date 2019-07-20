package crypt

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 生成md5摘要
func MD5(s string) string {
	data := []byte(s)
	m := md5.New()
	_, _ = m.Write(data)

	return hex.EncodeToString(m.Sum(nil))
}
