package crypt

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	aesKey = []byte(`xpLFVxGOaGozef4L`)
)

func TestAes(t *testing.T) {
	msg := []byte(`goutil`)
	result, err := AesCFBEncrypt(msg, aesKey)
	require.NoError(t, err)
	t.Logf("encrypt: %s", base64.StdEncoding.EncodeToString(result))

	result, err = AesCFBDecrypt(result, aesKey)
	require.NoError(t, err)
	require.Equal(t, msg, result)
	t.Logf("decrypt: %s", result)
}

func BenchmarkAesCFBEncrypt(b *testing.B) {
	msg := []byte(`goutil`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AesCFBEncrypt(msg, aesKey)
	}
}

func BenchmarkAesCFBDecrypt(b *testing.B) {
	m, _ := base64.StdEncoding.DecodeString(`M/7YtKkG0ped967YusafsDFgFqH68A==`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AesCFBDecrypt(m, aesKey)
	}
}
