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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"errors"

	"github.com/ouqiang/goutil"
)

func TestRequest_Get(t *testing.T) {
	content := "get method"
	handler := func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("name", req.URL.Query().Get("name"))
		rw.Header().Set("X-Forwarded-For", req.Header.Get("X-Forwarded-For"))
		_, _ = io.WriteString(rw, content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	params := url.Values{}
	params.Set("name", "golang")

	header := make(http.Header)
	header.Set("X-Forwarded-For", "8.8.8.8")
	resp, err := req.Get(s.URL, params, header)
	require.NoError(t, err)
	data, err := resp.String()
	require.NoError(t, err)
	require.Equal(t, content, data)
	require.Equal(t, "golang", resp.Header().Get("name"))
	require.Equal(t, "8.8.8.8", resp.Header().Get("X-Forwarded-For"))
}

func TestRequest_makeBody(t *testing.T) {
	req := NewRequest()
	r := req.makeBody(nil)
	require.Nil(t, r)
	s := "name=golang"
	r = req.makeBody(s)
	out, err := ioutil.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, s, string(out))

	b := []byte(s)
	r = req.makeBody(b)
	out, err = ioutil.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, s, string(out))

	v := url.Values{}
	v.Add("name", "golang")
	r = req.makeBody(v)
	out, err = ioutil.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, s, string(out))

	r = req.makeBody(strings.NewReader(s))
	out, err = ioutil.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, s, string(out))
	err = goutil.PanicToError(func() {
		r = req.makeBody(1)
	})
	require.NotNil(t, err)
}

func TestRequest_makeURLWithParams(t *testing.T) {
	req := NewRequest()
	baseURL := "https://golang.org"
	u := req.makeURLWithParams(baseURL, nil)
	require.Equal(t, baseURL, u)

	data := url.Values{}
	data.Set("name", "golang")
	u = req.makeURLWithParams(baseURL, data)
	expected := baseURL + "?name=golang"
	require.Equal(t, expected, u)

	data = url.Values{}
	data.Set("name", "golang")
	baseURL += "?"
	u = req.makeURLWithParams(baseURL, data)
	require.Equal(t, expected, u)
}

func TestRequest_shouldRetry(t *testing.T) {
	resp := new(http.Response)
	err := errors.New("should retry")

	req := NewRequest()
	require.True(t, req.shouldRetry(nil, resp, err))

	err = nil
	resp.StatusCode = http.StatusNotFound
	require.True(t, req.shouldRetry(nil, resp, err))

	err = nil
	resp.StatusCode = http.StatusOK
	require.False(t, req.shouldRetry(nil, resp, err))

}

func TestRequest_Post(t *testing.T) {
	handler := func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("name", req.FormValue("name"))
		rw.Header().Set("From", req.Header.Get("From"))
		rw.Header().Set("X-Forwared-For", req.Header.Get("X-Forwared-For"))
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	data := url.Values{}
	data.Set("name", "golang")

	header := make(http.Header)
	header.Set("From", "Post")
	header.Set("X-Forwared-For", "8.8.8.8")
	resp, err := req.Post(s.URL, data, header)
	require.NoError(t, err)
	require.Equal(t, "golang", resp.Header().Get("name"))
	require.Equal(t, "Post", resp.Header().Get("From"))
	require.Equal(t, "8.8.8.8", resp.Header().Get("X-Forwared-For"))
}

func TestRequest_PostJSON(t *testing.T) {
	jsonString := `{"code":200,"message":"","data":{"name":"golang"}}`
	handler := func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Content-Type") != "application/json" {
			panic("invalid content-type")
		}
		_, _ = io.Copy(rw, req.Body)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.PostJSON(s.URL, jsonString, nil)
	require.NoError(t, err)
	body, err := resp.String()
	require.NoError(t, err)
	require.Equal(t, jsonString, body)
}

func TestRequest_PostProtoBuf(t *testing.T) {
	message := &Message{
		Name: "protobuf",
	}
	handler := func(rw http.ResponseWriter, req *http.Request) {
		_, _ = io.Copy(rw, req.Body)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.PostProtoBuf(s.URL, message, nil)
	require.NoError(t, err)

	result := &Message{}
	err = resp.DecodeProtoBuf(result)
	require.NoError(t, err)
	require.Equal(t, message.Name, result.Name)
}

func TestRequest_UploadFile(t *testing.T) {
	fileContent := "test file content"

	params := map[string]string{
		"name": "httpclient",
	}
	handler := func(rw http.ResponseWriter, req *http.Request) {
		file, _, err := req.FormFile("file")
		if err != nil {
			panic(fmt.Errorf("读取文件错误: %s", err))
		}
		name := req.FormValue("name")
		if name != "httpclient" {
			panic("invalid post params")
		}

		defer func() {
			_ = file.Close()
		}()
		_, _ = io.Copy(rw, file)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.UploadFile(s.URL, strings.NewReader(fileContent), "upload.txt", nil, params)
	require.NoError(t, err)
	responseContent, err := resp.String()
	require.NoError(t, err)
	require.Equal(t, fileContent, responseContent)
}

func TestRequest_SetRetryTimes(t *testing.T) {
	retryTimes := 0
	handler := func(rw http.ResponseWriter, req *http.Request) {
		if retryTimes == 2 {
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
		retryTimes--
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()
	req := NewRequest(WithRetryTime(retryTimes))
	resp, err := req.Get(s.URL, nil, nil)
	require.NoError(t, err)
	require.False(t, resp.IsStatusOK())
	require.Equal(t, -1, retryTimes)

	retryTimes = 3
	req = NewRequest(WithRetryTime(retryTimes))
	resp, err = req.Get(s.URL, nil, nil)
	require.NoError(t, err)
	require.True(t, resp.IsStatusOK())
	require.Equal(t, 1, retryTimes)

	retryTimes = 1
	req = NewRequest(WithRetryTime(retryTimes))
	resp, err = req.Get(s.URL, nil, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.Raw().StatusCode)
	require.Equal(t, -1, retryTimes)
}
