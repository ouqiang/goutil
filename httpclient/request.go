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
	"context"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
)

const (
	defaultTimeout             = 20 * time.Second
	defaultConnectTimeout      = 5 * time.Second
	defaultMaxIdleConnsPerHost = 2
)

var (
	// 如果设置了Accept-Encoding, 不会自动解压
	defaultHeader = map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7",
		"Cache-Control":   "no-cache",
		"Pragma":          "no-cache",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.170 Safari/537.36",
	}
)

type RequestInterceptor func(req *http.Request)
type ResponseInterceptor func(req *http.Request, resp *http.Response, err error)

type options struct {
	client              *http.Client
	debug               bool
	cookieJar           http.CookieJar
	timeout             time.Duration
	connectTimeout      time.Duration
	maxIdleConnsPerHost int
	proxyURL            string
	retryTimes          int
	enableDefaultHeader bool
	disableKeepAlive    bool
	dnsResolver         DNSResolverFunc
	unixSocketPath      string
	shouldRetryFunc     func(*http.Request, *http.Response, error) bool
	requestInterceptor  RequestInterceptor
	responseInterceptor ResponseInterceptor
	clientTrace         *httptrace.ClientTrace
}

// DNSResolver DNS解析
type DNSResolverFunc func(host string) (ip string, err error)

type DialContext func(ctx context.Context, network, addr string) (net.Conn, error)

type Option func(*options)

// WithClient 自定义http client
func WithClient(client *http.Client) Option {
	return func(opt *options) {
		opt.client = client
	}
}

func WithClientTrace(t *httptrace.ClientTrace) Option {
	return func(opt *options) {
		opt.clientTrace = t
	}
}

func WithRequestInterceptor(i RequestInterceptor) Option {
	return func(opt *options) {
		opt.requestInterceptor = i
	}
}

func WithResponseInterceptor(i ResponseInterceptor) Option {
	return func(opt *options) {
		opt.responseInterceptor = i
	}
}

func WithCookieJar(cookieJar http.CookieJar) Option {
	return func(opt *options) {
		opt.cookieJar = cookieJar
	}
}

func WithDNSResolver(dnsResolver DNSResolverFunc) Option {
	return func(opt *options) {
		opt.dnsResolver = dnsResolver
	}
}

// WithShouldRetryFunc 自定义是否需要重试
func WithShouldRetryFunc(f func(*http.Request, *http.Response, error) bool) Option {
	return func(opt *options) {
		opt.shouldRetryFunc = f
	}
}

// WithEnableDefaultHeader 设置默认header
func WithEnableDefaultHeader() Option {
	return func(opt *options) {
		opt.enableDefaultHeader = true
	}
}

// WithRetryTime 设置重试次数
func WithRetryTime(retryTimes int) Option {
	return func(opt *options) {
		opt.retryTimes = retryTimes
	}
}

// WithProxyURL 设置代理
func WithProxyURL(proxyURL string) Option {
	return func(opt *options) {
		opt.proxyURL = proxyURL
	}
}

func WithUnixSocketPath(path string) Option {
	return func(opt *options) {
		opt.unixSocketPath = opt.unixSocketPath
	}
}

// WithTimeout 设置超时
func WithTimeout(timeout time.Duration) Option {
	return func(opt *options) {
		opt.timeout = timeout
	}
}

// WithConnectTimeout 设置连接超时
func WithConnectTimeout(timeout time.Duration) Option {
	return func(opt *options) {
		opt.connectTimeout = timeout
	}
}

// WithMaxIdleConnsPerHost 设置每个host最大空闲连接数
func WithMaxIdleConnsPerHost(n int) Option {
	return func(opt *options) {
		opt.maxIdleConnsPerHost = n
	}
}

// WithDisableKeepAlive 连接重用
func WithDisableKeepAlive() Option {
	return func(opt *options) {
		opt.disableKeepAlive = true
	}
}

// WithDebug 开启调试模式
func WithDebug() Option {
	return func(opt *options) {
		opt.debug = true
	}
}

// Request http请求
type Request struct {
	opts options
}

// NewRequest 创建request
func NewRequest(opt ...Option) *Request {
	req := &Request{}
	req.opts = options{}
	for _, o := range opt {
		o(&req.opts)
	}
	req.init()

	return req
}

func (req *Request) init() {
	if req.opts.connectTimeout <= 0 {
		req.opts.connectTimeout = defaultConnectTimeout
	}
	if req.opts.maxIdleConnsPerHost <= 0 {
		req.opts.maxIdleConnsPerHost = defaultMaxIdleConnsPerHost
	}
	trans := &http.Transport{
		Proxy: func(request *http.Request) (*url.URL, error) {
			if req.opts.proxyURL != "" {
				return url.Parse(req.opts.proxyURL)
			}

			return nil, nil
		},
		DialContext: (&net.Dialer{
			Timeout:   req.opts.connectTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   req.opts.maxIdleConnsPerHost,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if req.opts.disableKeepAlive {
		trans.DisableKeepAlives = true
	}

	if req.opts.client == nil {
		req.opts.client = &http.Client{}
	}
	if req.opts.dnsResolver != nil {
		trans.DialContext = req.dialContext()
	}
	if req.opts.unixSocketPath != "" {
		trans.DialContext = req.dialContextForUnixDomainSocket
	}
	if req.opts.client.Transport == nil {
		req.opts.client.Transport = trans
	}
	if req.opts.timeout == 0 {
		req.opts.client.Timeout = defaultTimeout
	}
	if req.opts.shouldRetryFunc == nil {
		req.opts.shouldRetryFunc = req.shouldRetry
	}
	if req.opts.cookieJar != nil {
		req.opts.client.Jar = req.opts.cookieJar
	}
}

// Get get请求
func (req *Request) Get(url string, params url.Values, header http.Header) (*Response, error) {
	url = req.makeURLWithParams(url, params)

	return req.do(http.MethodGet, url, nil, header)
}

// Post 普通post请求
func (req *Request) Post(url string, data interface{}, header http.Header) (*Response, error) {

	return req.do(http.MethodPost, url, data, header)
}

// Put Put请求
func (req *Request) Put(url string, data interface{}, header http.Header) (*Response, error) {
	return req.do(http.MethodPut, url, data, header)
}

// Delete Delete请求
func (req *Request) Delete(url string, data interface{}, header http.Header) (*Response, error) {
	return req.do(http.MethodDelete, url, data, header)
}

// PostJSON 发送json body
func (req *Request) PostJSON(url string, data interface{}, header http.Header) (*Response, error) {
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/json")
	var body interface{}
	switch data.(type) {
	case string, []byte, io.Reader:
		body = data
	default:
		var err error
		body, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	return req.do(http.MethodPost, url, body, header)
}

// PostProtoBuf 发送protoBuf body
func (req *Request) PostProtoBuf(url string, v proto.Message, header http.Header) (*Response, error) {
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	body, err := proto.Marshal(v)
	if err != nil {
		return nil, err
	}

	return req.do(http.MethodPost, url, body, header)
}

// UploadFile 上传文件
func (req *Request) UploadFile(url string, reader io.Reader, filename string, header http.Header, params map[string]string) (*Response, error) {
	pipeReader, pipeWriter := io.Pipe()
	mr := multipart.NewWriter(pipeWriter)
	go func() {
		defer func() {
			_ = mr.Close()
			_ = pipeWriter.Close()
		}()

		fileFieldName := "file"
		if len(params) > 0 && params["_file_field_name"] != "" {
			fileFieldName = params["_file_field_name"]
		}

		var part io.Writer
		part, err := mr.CreateFormFile(fileFieldName, filename)
		if err != nil {
			return
		}
		_, _ = io.Copy(part, reader)
		for k, v := range params {
			_ = mr.WriteField(k, v)
		}
	}()
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", mr.FormDataContentType())

	resp, respErr := req.Post(url, pipeReader, header)

	return resp, respErr
}

func (req *Request) do(method string, url string, data interface{}, header http.Header) (*Response, error) {
	targetReq, err := req.build(method, url, data, header)
	if err != nil {
		return nil, err
	}
	req.beforeRequest(targetReq)
	execTimes := 1
	retryInterval := 300 * time.Millisecond
	if req.opts.retryTimes > 0 {
		execTimes += req.opts.retryTimes
	}
	var resp *http.Response
	for i := 0; i < execTimes; {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		resp, err = req.opts.client.Do(targetReq)
		req.afterResponse(targetReq, resp, err)
		if req.opts.retryTimes > 0 && !req.opts.shouldRetryFunc(targetReq, resp, err) {
			break
		}
		i++
		retryInterval *= 2
		if i < execTimes {
			time.Sleep(retryInterval)
		}
	}
	if err != nil {
		return nil, err
	}

	return newResponse(resp), err
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
	if method != http.MethodGet && header.Get("Content-Type") == "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	targetReq.Header = header
	host := header.Get("Host")
	if host != "" {
		targetReq.Host = host
	}
	if req.opts.enableDefaultHeader {
		for k, v := range defaultHeader {
			targetReq.Header.Add(k, v)
		}
	}
	if req.opts.clientTrace != nil {
		targetReq = targetReq.WithContext(httptrace.WithClientTrace(targetReq.Context(), req.opts.clientTrace))

	}

	return targetReq, nil
}

func (req *Request) beforeRequest(r *http.Request) {
	if req.opts.requestInterceptor != nil {
		req.opts.requestInterceptor(r)
	}
	req.dumpRequestIfNeed(r)
}

func (req *Request) afterResponse(r *http.Request, resp *http.Response, err error) {
	if req.opts.responseInterceptor != nil {
		req.opts.responseInterceptor(r, resp, err)
	}
	req.dumpResponseIfNeed(resp, err)
}

// request调试输出
func (req *Request) dumpRequestIfNeed(r *http.Request) {
	if !req.opts.debug {
		return
	}
	reqDump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		panic(err)
	}
	log.Printf("[Request]\n\n%s\n", reqDump)
}

// response调试输出
func (req *Request) dumpResponseIfNeed(resp *http.Response, err error) {
	if !req.opts.debug {
		return
	}
	if err != nil {
		log.Printf("[Response]\n\n%s\n", err)
		return
	}
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Printf("[Response]\n\n %s\n", respDump)
}

// 是否要重试
func (req *Request) shouldRetry(request *http.Request, resp *http.Response, err error) bool {
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

func (req *Request) dialContext() DialContext {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout:   req.opts.connectTimeout,
			KeepAlive: 30 * time.Second,
		}

		u, err := url.Parse(addr)
		if err == nil && u.Scheme == "unix" {
			return dialer.DialContext(ctx, u.Scheme, u.Path)
		}

		separator := strings.LastIndex(addr, ":")
		ip, err := req.opts.dnsResolver(addr[:separator])
		if err != nil {
			return nil, err
		}

		addr = ip + addr[separator:]

		return dialer.DialContext(ctx, network, addr)
	}
}

func (req *Request) dialContextForUnixDomainSocket(ctx context.Context, network, address string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   req.opts.connectTimeout,
		KeepAlive: 30 * time.Second,
	}

	return dialer.DialContext(ctx, "unix", req.opts.unixSocketPath)
}
