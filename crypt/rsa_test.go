package crypt

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	rsaPublicKeyFile = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1NygfCvp72abPmF47ArO
o2L7zAuAsSPQiPF1JUqB2Y5vrlfPiyq8QXYIqMPHvz2isr+8T6RR6MzGDp5vCo5U
sSVfQiS8CptsEZzLiXawhQl5bzgn21J721nG/+c/+Wc3Un9vZN22zPEbxKUR7P1k
M5UveU2J5gbLk6+b5UsHTDOu3Ity39TKAfjuKRvxVFbSkkNHoxZX2OmKHwWBLZht
K9Og7BtvH9Mw0B9WKUsRfutN+uE6/Nbaeri6d/l5kxGu0eJbH5yTFSbuqAH8Gu9G
Coghz7pOO2CB4kD8ksykQkM7ugW10SelsmJ1O76JBkCGEuTi2sg3QIkCqWIko9z4
oQIDAQAB
-----END PUBLIC KEY-----
`
	rsaPrivateKeyFile = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA1NygfCvp72abPmF47ArOo2L7zAuAsSPQiPF1JUqB2Y5vrlfP
iyq8QXYIqMPHvz2isr+8T6RR6MzGDp5vCo5UsSVfQiS8CptsEZzLiXawhQl5bzgn
21J721nG/+c/+Wc3Un9vZN22zPEbxKUR7P1kM5UveU2J5gbLk6+b5UsHTDOu3Ity
39TKAfjuKRvxVFbSkkNHoxZX2OmKHwWBLZhtK9Og7BtvH9Mw0B9WKUsRfutN+uE6
/Nbaeri6d/l5kxGu0eJbH5yTFSbuqAH8Gu9GCoghz7pOO2CB4kD8ksykQkM7ugW1
0SelsmJ1O76JBkCGEuTi2sg3QIkCqWIko9z4oQIDAQABAoIBAQCSeTMx3uH+K+P9
3ALiolkpEw7jjLLEsjloUobi309VDm+JT8FwKTsb7IXu47HKxjXzvH5va0o14NsU
6i7M274pm+bd0/tVbTfDMtrYP8Ud9rQKNWEvyaWS/kcyhsV98GmgKjLT/MEya2NJ
QLGCG6dc23asiQ+wKtLhUID9GlDuPIzONXvzMlYC/Acbp3pLtaSuFH+yQYg7DS1s
8IjJVJGFtNdw8qy0TILLNaMXXizbzY+E6XZfpBBF1YnbVGfnE52VkesCRzoaPXlw
9+TGC5M9xN0OUsx6/4EQ+ja8pX93GCZO4BB7qdbs9ubVHqaUiLE19j5HszfHc8KG
TfVsMImxAoGBAPGACj1fi1ic7WxqgZ7uGy2WlA1b90q/m0oWyynvNBT658E4vggK
pJ0P2QVEkqixdCvZ7HnhH4lMMj9FXSmZKzFtnnTkK46ZVB6yB69YKVlCENOC4VBT
ltCHITlbKo9JCdmcQePmc4eI70DQ4Hu1OpZIvWsSlSXKWJ3nPB6AOP/1AoGBAOGk
Z0+qp+JdK9RDHYJ08ny7SCb6hO1iwOj9yewDNIZ0SjI9iK717kWBdFLUQyDinTux
lf9+fMNPuTQx418Vrc9zybwH5WjQ3LVkQB/HhcQKLx3qwYI05kwNqKx3UXM7IdzU
DS1lyeJqj+yRZF/qmn2N/wDbQkR/JMvbffZp/0Z9AoGAWw/0zH+ig59ox2DBz7Po
+5+z4/WxobXuUFmX8hAIi4CwsuM6hL2+pJq9ModQ5dtD+uUJjkudIKBisgjtwCnJ
Z7H19g2zfunCFnD7BNsxfD61KYxIeYmLbMYHeSEvQyg/VpbdIZpcJdc0oDQi8YK8
vL5g7mbrZnyOPbxbpwSieaECgYAWWvB1XyYE8lAuVnvs+eMwYmmymu1ii38rVkGU
JXklvQ3AzoHlO65gqoO41RjVgD4ttNl0l7aKrzJdLnglaoNu4zzgaTPcX50OR6Fm
xKDHHG8wmpqTaORMMqo8dBHYxcoEE+o+TjBjQ0WBHaKBMkAeIlxaXF7DZIljvRpM
uJG3DQKBgC2t/+48UBnWPCzLgl1L6LC/m90bHz3qUMx/kQzhpcavEoYh4rUVkOFO
NjxB+UkCkGaJdyqOGcav2m/9hM6ItdcInO+FKmHGoFOCAhmuLNWDJ9xuEKqaV6R0
udOG6qz3s77VAvAop3akyf6ROOk1zAYIjTYGCaaJs5eCEL6pj3wt
-----END RSA PRIVATE KEY-----
`
)

func TestRsa(t *testing.T) {
	msg := []byte("goutil")

	result, err := rsaEncrypt.Encrypt(msg)
	require.NoError(t, err)
	t.Logf("%s", base64.StdEncoding.EncodeToString(result))

	result, err = rsaEncrypt.Decrypt(result)
	require.NoError(t, err)

	require.Equal(t, msg, result)
}

func BenchmarkRsa_Encrypt(b *testing.B) {
	msg := []byte("goutil")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, _ := rsaEncrypt.Encrypt(msg)
		_ = result
	}
}

func BenchmarkRsa_Decrypt(b *testing.B) {
	data, err := base64.StdEncoding.DecodeString(`ZMNBpfCxhQfsH1P/TtkckmtcafjJzgx23XvX1vBBTPz7kUbG38L0MksiaDW/JlmroFPZXBUa2s/QINmFa5z11LvypA39YE3eRMRXkf9JkPDDQhrYo51LcE8jMUxMgxLxzudCI1ZKbQ5jP9BNzPgWpcOA
PJ+JZd16xOhosldrmEQiBNnW+pE3b0B3KjOzrJYpWa6pACHKWnpzh9xUv8ftgWZKMdy0Pkce1qEXWFjldBeU1430fjB5aLWwokQkkSMJFFTQ2xzD1TGJSiLWit0vY4lkWOYwhjG9rxttQCfq7JKEJcYWPpECtaniBovwTYM5yEfJ
BWADj684TWebyJFwrQ==`)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, _ := rsaEncrypt.Decrypt(data)
		_ = result
	}
}
