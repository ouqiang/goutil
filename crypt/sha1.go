package crypt

import (
	"crypto/sha1"
	"encoding/hex"
)

// SHA1 生成sha1摘要
func SHA1(s string) string {
	m := sha1.New()
	_, _ = m.Write([]byte(s))

	return hex.EncodeToString(m.Sum(nil))
}
