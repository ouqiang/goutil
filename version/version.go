// Package version 版本
package version

import "github.com/ouqiang/goutil"

var (
	// 应用版本
	appVersion string
	// 应用构建日期
	buildDate string
	// 最后提交的git commit id
	gitCommit string
)

// Init 初始化版本信息
func Init(version, date, commit string) {
	appVersion = version
	buildDate = date
	gitCommit = commit

}

// Get 获取应用版本
func Get() string {
	return appVersion
}

// Format 格式化输出
func Format() string {
	v, err := goutil.FormatAppVersion(appVersion, gitCommit, buildDate)
	if err != nil {
		panic(err)
	}

	return v
}
