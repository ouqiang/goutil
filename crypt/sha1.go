package crypt

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

// SHA1 生成sha1摘要
func SHA1(s string) string {
	m := sha1.New()
	_, _ = m.Write([]byte(s))

	return hex.EncodeToString(m.Sum(nil))
}

// SHA1Stream 流式处理
func SHA1Stream(reader io.Reader) (string, error) {
	m := sha1.New()
	_, err := io.Copy(m, reader)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(m.Sum(nil)), nil
}

func SHA1StreamSum(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	return SHA1Stream(file)
}
