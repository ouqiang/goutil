package crypt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSHA1(t *testing.T) {
	s := "goutil"
	require.Equal(t, "05cb8eaebc8f5d3514449b4f6929ae0fb6a036a1", SHA1(s))
}

func BenchmarkSHA1(b *testing.B) {
	s := "goutil"
	for i := 0; i < b.N; i++ {
		SHA1(s)
	}
}
