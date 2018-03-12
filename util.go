// Copyright 2018 ouqiang authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package goutil

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"runtime"
	"time"
)

// MD5 生成MD5摘要
func MD5(s string) string {
	m := md5.New()
	m.Write([]byte(s))

	return hex.EncodeToString(m.Sum(nil))
}

// RandNumber 生成min - max之间的随机数
func RandNumber(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return min + r.Intn(max-min)
}

// PanicToError Panic转换为error
func PanicToError(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic: %s", e)
		}
	}()
	f()
	return
}

// PrintAppVersion 打印应用版本
func PrintAppVersion(appVersion, GitCommit, BuildDate string) {
	content := `
   Version: {{.Version}}
Go Version: {{.GoVersion}}
Git Commit: {{.GitCommit}}
     Built: {{.BuildDate}}
   OS/ARCH: {{.GOOS}}/{{.GOARCH}}
`
	tpl, err := template.New("version").Parse(content)
	if err != nil {
		panic("parse version template error: " + err.Error())
	}
	err = tpl.Execute(os.Stdout, map[string]string{
		"Version":   appVersion,
		"GoVersion": runtime.Version(),
		"GitCommit": GitCommit,
		"BuildDate": BuildDate,
		"GOOS":      runtime.GOOS,
		"GOARCH":    runtime.GOARCH,
	})
	if err != nil {
		panic("replace version template error: " + err.Error())
	}
}
