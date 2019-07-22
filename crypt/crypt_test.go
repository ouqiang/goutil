package crypt

import (
	"os"
	"strings"
	"testing"
)

var (
	rsaEncrypt *Rsa
)

func TestMain(m *testing.M) {
	publicKeyReader := strings.NewReader(rsaPublicKeyFile)
	privateKeyReader := strings.NewReader(rsaPrivateKeyFile)
	var err error
	rsaEncrypt, err = NewRsa(publicKeyReader, privateKeyReader)
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}
