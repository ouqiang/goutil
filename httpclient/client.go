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

// Package httpclient http客户端
package httpclient

import (
	"net/http"
	"net/url"
	"time"
)

var (
	defaultClient = NewRequest(
		WithConnectTimeout(2*time.Second),
		WithTimeout(10*time.Second),
		WithEnableDefaultHeader(),
		WithMaxIdleConnsPerHost(5),
	)
)

// Get get请求
func Get(url string, data url.Values, header http.Header) (*Response, error) {
	return defaultClient.Get(url, data, header)
}

// Post 普通post请求
func Post(url string, data interface{}, header http.Header) (*Response, error) {
	return defaultClient.Post(url, data, header)
}

// PostJSON 发送json body
func PostJSON(url string, data interface{}, header http.Header) (*Response, error) {
	return defaultClient.PostJSON(url, data, header)
}
