// Package httpclient http 客户端
package httpclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// 默认超时时间
	defaultTimeout = 20 * time.Second
	// 默认UserAgent
	defaultUserAgent = "golang/http-client"
)

// RequestOption 请求可选项
type RequestOption struct {
	// Timeout 超时
	Timeout time.Duration
	// Header 请求Header
	Header http.Header
	// Body 请求Body
	Body io.Reader
	// 自动执行重定向
	AutoRedirect bool
}

// ResponseWrapper 响应包装
type ResponseWrapper struct {
	// StatusCode 响应码
	StatusCode int
	// Body 响应Body
	Body []byte
	// Header 响应Header
	Header http.Header
}

// Get 执行Get请求
func Get(url string, opt *RequestOption) (*ResponseWrapper, error) {
	if opt == nil {
		opt = &RequestOption{}
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return createRequestError(err)
	}

	return request(req, opt)
}

// PostParams POST请求按key value传递参数
func PostParams(url string, opt *RequestOption) (*ResponseWrapper, error) {
	if opt == nil {
		opt = &RequestOption{}
	}
	req, err := http.NewRequest("POST", url, opt.Body)
	if err != nil {
		return createRequestError(err)
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	return request(req, opt)
}

// PostJSON POST请求，发送JSON格式数据
func PostJSON(url string, opt *RequestOption) (*ResponseWrapper, error) {
	if opt == nil {
		opt = &RequestOption{}
	}
	req, err := http.NewRequest("POST", url, opt.Body)
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-type", "application/json")

	return request(req, opt)
}

// 执行请求
func request(req *http.Request, opt *RequestOption) (*ResponseWrapper, error) {
	client := &http.Client{}
	setClientOption(client, opt)
	setRequestHeader(req, opt)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("执行HTTP请求错误#%s", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取HTTP响应Body失败#%s", err.Error())
	}
	wrapper := &ResponseWrapper{}
	wrapper.StatusCode = resp.StatusCode
	wrapper.Body = body
	wrapper.Header = resp.Header

	return wrapper, nil
}

// 设置客户端选项
func setClientOption(c *http.Client, opt *RequestOption) {
	if opt.Timeout > 0 {
		c.Timeout = opt.Timeout
	} else {
		c.Timeout = defaultTimeout
	}
	if !opt.AutoRedirect {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}

// 设置请求Header
func setRequestHeader(req *http.Request, opt *RequestOption) {
	req.Header.Set("User-Agent", defaultUserAgent)
	if len(opt.Header) == 0 {
		return
	}
	for key := range opt.Header {
		req.Header.Set(key, opt.Header.Get(key))
	}
}

func createRequestError(err error) (*ResponseWrapper, error) {
	return nil, fmt.Errorf("创建HTTP请求错误#%s", err.Error())
}
