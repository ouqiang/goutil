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

package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout = 20 * time.Second
)

var (
	defaultClient = &http.Client{Timeout: defaultTimeout}
	// 如果设置了Accept-Encoding, 不会自动解压
	defaultHeader = map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7",
		"Cache-Control":   "no-cache",
		"Pragma":          "no-cache",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.170 Safari/537.36",
	}
)

// Request http请求
// NOTICE: 不支持并发调用
type Request struct {
	client              *http.Client
	debug               bool
	timeout             time.Duration
	proxyURL            string
	retryTimes          int
	enableDefaultHeader bool
	shouldRetryFunc     func(*http.Response, error) bool
}

// NewRequest 创建request
func NewRequest() *Request {
	return &Request{}
}

// SetClient 自定义http client
func (req *Request) SetClient(client *http.Client) *Request {
	req.client = client
	return req
}

// SetTimeout 设置超时
func (req *Request) SetTimeout(timeout time.Duration) *Request {
	req.timeout = timeout

	return req
}

// SetRetryTimes 设置重试次数
func (req *Request) SetRetryTimes(retryTimes int) *Request {
	req.retryTimes = retryTimes

	return req
}

// SetShouldRetryFunc 设置是否应该重试方法
func (req *Request) SetShouldRetryFunc(f func(*http.Response, error) bool) *Request {
	req.shouldRetryFunc = f

	return req
}

// EnableDebug 开启调试模式
func (req *Request) EnableDebug() *Request {
	req.debug = true

	return req
}

// EnableDefaultHeader 自动设置默认header
func (req *Request) EnableDefaultHeader() *Request {
	req.enableDefaultHeader = true

	return req
}

// SetProxy 设置代理
func (req *Request) SetProxy(proxyURL string) *Request {
	req.proxyURL = proxyURL

	return req
}

// Get get请求
func (req *Request) Get(url string, data url.Values, header http.Header) (*Response, error) {
	url = req.makeURLWithParams(url, data)

	return req.do(http.MethodGet, url, nil, header)
}

// Post 普通post请求
func (req *Request) Post(url string, data interface{}, header http.Header) (*Response, error) {
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req.do(http.MethodPost, url, data, header)
}

// PostJSON 发送json body
func (req *Request) PostJSON(url string, data interface{}, header http.Header) (*Response, error) {
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/json")

	return req.do(http.MethodPost, url, data, header)
}

func (req *Request) do(method string, url string, data interface{}, header http.Header) (*Response, error) {
	targetReq, err := req.build(method, url, data, header)
	if err != nil {
		return nil, err
	}
	req.setClientIfNeed()
	req.beforeRequest(targetReq)
	execTimes := 1
	if req.retryTimes > 0 {
		execTimes += req.retryTimes
	}
	if req.shouldRetryFunc == nil {
		req.shouldRetryFunc = req.shouldRetry
	}
	var resp *http.Response
	for i := 0; i < execTimes; i++ {
		resp, err = req.client.Do(targetReq)
		req.afterResponse(resp, err)
		if req.retryTimes > 0 && !req.shouldRetryFunc(resp, err) {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	return newResponse(resp), err
}

func (req *Request) setClientIfNeed() {
	if req.client != nil {
		return
	}
	if req.timeout == 0 && req.proxyURL == "" {
		req.client = defaultClient
		return
	}
	req.client = &http.Client{}
	if req.timeout > 0 {
		req.client.Timeout = req.timeout
	} else {
		req.client.Timeout = defaultTimeout
	}

	if req.proxyURL != "" {
		req.client.Transport = req.transport()
	}
}

func (req *Request) transport() http.RoundTripper {
	trans := &http.Transport{
		Proxy: func(request *http.Request) (*url.URL, error) {
			if req.proxyURL != "" {
				return url.Parse(req.proxyURL)
			}

			return nil, nil
		},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return trans
}

// 构造http.Request
func (req *Request) build(method string, url string, data interface{}, header http.Header) (*http.Request, error) {
	body := req.makeBody(data)
	targetReq, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if header == nil {
		header = make(http.Header)
	}
	targetReq.Header = header
	host := header.Get("Host")
	if host != "" {
		targetReq.Host = host
	}
	if req.enableDefaultHeader {
		for k, v := range defaultHeader {
			targetReq.Header.Add(k, v)
		}
	}

	return targetReq, nil
}

func (req *Request) beforeRequest(r *http.Request) {
	req.dumpRequestIfNeed(r)
}

func (req *Request) afterResponse(resp *http.Response, err error) {
	req.dumpResponseIfNeed(resp, err)
}

// request调试输出
func (req *Request) dumpRequestIfNeed(r *http.Request) {
	if !req.debug {
		return
	}
	reqDump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[Request]\n\n%s\n", reqDump)
}

// response调试输出
func (req *Request) dumpResponseIfNeed(resp *http.Response, err error) {
	if !req.debug {
		return
	}
	if err != nil {
		fmt.Printf("[Response]\n\n%s\n", err)
		return
	}
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[Response]\n\n %s\n", respDump)
}

// 是否要重试
func (req *Request) shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	if resp == nil {
		return true
	}
	if resp.StatusCode != http.StatusOK {
		return true
	}

	return false
}

// 参数追加到url末尾
func (req *Request) makeURLWithParams(url string, data url.Values) string {
	if url == "" {
		return url
	}
	if data == nil {
		return url
	}
	params := data.Encode()
	if strings.Contains(url, "?") {
		if url[len(url)-1] != '?' {
			url += "&"
		}
	} else {
		url += "?"
	}
	url += params

	return url
}

// 生成请求Body
func (req *Request) makeBody(data interface{}) io.Reader {
	if data == nil {
		return nil
	}
	var body io.Reader
	switch v := data.(type) {
	case string:
		body = strings.NewReader(v)
	case []byte:
		body = bytes.NewBuffer(v)
	case url.Values:
		body = strings.NewReader(v.Encode())
	case io.Reader:
		body = v
	default:
		panic("data is not support type")
	}

	return body
}
