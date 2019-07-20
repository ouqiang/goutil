package crypt

import (
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
