package crypt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMD5(t *testing.T) {
	s := "goutil"
	require.Equal(t, "3fa6e676e7d5c558e9a49b599cbc975f", MD5(s))
}

func BenchmarkMD5(b *testing.B) {
	s := "goutil"
	for i := 0; i < b.N; i++ {
		MD5(s)
	}
}

func TestMd5Stream(t *testing.T) {
	s := "goutil"
	result, err := Md5Stream(strings.NewReader(s))
	require.NoError(t, err)
	require.Equal(t, "3fa6e676e7d5c558e9a49b599cbc975f", result)
}
