// Package environment 环境
package environment

const (
	DevMode  = "dev"   // 开发环境
	TestMode = "test"  // 测试环境
	PreMode   = "pre"  // 预发布环境
	ProdMode = "prod"  // 生产环境
)

var mode string

func Init(m string) {
	mode = m
}

func Mode() string {
	return mode
}

func IsDev() bool {
	return mode == DevMode
}

func IsTest() bool {
	return mode == TestMode
}

func IsPre() bool {
	return mode == PreMode
}

func IsProd() bool {
	return mode == ProdMode
}
