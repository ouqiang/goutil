package crypt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHMacMD5(t *testing.T) {
	s := "goutil"
	key := []byte("test")
	require.Equal(t, "127e4783c022e9ca01ece00c617dc240", HMacMD5(s, key))
}

func BenchmarkHMacMD5(b *testing.B) {
	s := "goutil"
	key := []byte("test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HMacMD5(s, key)
	}
}
