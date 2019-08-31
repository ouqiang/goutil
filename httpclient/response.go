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
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gogo/protobuf/proto"
)

// Response http响应
type Response struct {
	rawResp *http.Response
}

// newResponse 创建response
func newResponse(resp *http.Response) *Response {
	return &Response{resp}
}

// IsStatusOK 响应码是否为200
func (resp *Response) IsStatusOK() bool {
	return resp.rawResp.StatusCode == http.StatusOK
}

// DecodeJSON  json decode
func (resp *Response) DecodeJSON(v interface{}) error {
	err := json.NewDecoder(resp.rawResp.Body).Decode(v)
	_ = resp.rawResp.Body.Close()

	return err
}

// DecodeProtoBuf protoBuf decode
func (resp *Response) DecodeProtoBuf(v interface{}) error {
	data, err := resp.Bytes()
	if err != nil {
		return err
	}

	err = proto.Unmarshal(data, v.(proto.Message))

	return err
}

// String 读取http.Body, 返回string
func (resp *Response) String() (string, error) {
	b, err := resp.Bytes()
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// Bytes 读取http.Body, 返回bytes
func (resp *Response) Bytes() ([]byte, error) {
	b, err := ioutil.ReadAll(resp.rawResp.Body)
	_ = resp.rawResp.Body.Close()

	return b, err
}

// Discard 丢弃http.body
func (resp *Response) Discard() (int64, error) {
	n, err := io.Copy(ioutil.Discard, resp.rawResp.Body)
	_ = resp.rawResp.Body.Close()

	return n, err
}

// WriteFile 读取http.Body内容并写入文件中
func (resp *Response) WriteFile(filename string) (int64, error) {
	f, err := os.Create(filename)
	if err != nil {
		_, _ = resp.Discard()
		return 0, err
	}
	defer func() {
		_ = f.Close()
		_ = resp.rawResp.Body.Close()
	}()

	return io.Copy(f, resp.rawResp.Body)
}

// WriteTo 读取http.Body并写入w中
func (resp *Response) WriteTo(w io.Writer) (int64, error) {
	n, err := io.Copy(w, resp.rawResp.Body)
	_ = resp.rawResp.Body.Close()

	return n, err
}

// Header 获取header
func (resp *Response) Header() http.Header {
	return resp.rawResp.Header
}

// Raw 获取原始的http response
func (resp *Response) Raw() *http.Response {
	return resp.rawResp
}
