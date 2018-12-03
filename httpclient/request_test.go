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
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"errors"

	"github.com/ouqiang/goutil"
)

func TestRequest_Get(t *testing.T) {
	content := "get method"
	handler := func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("name", req.URL.Query().Get("name"))
		rw.Header().Set("X-Forwarded-For", req.Header.Get("X-Forwarded-For"))
		io.WriteString(rw, content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	params := url.Values{}
	params.Set("name", "golang")

	header := make(http.Header)
	header.Set("X-Forwarded-For", "8.8.8.8")
	resp, err := req.Get(s.URL, params, header)
	if err != nil {
		t.Fatal(err)
	}
	data, err := resp.String()
	if err != nil {
		t.Fatal(err)
	}
	if data != content {
		t.Errorf("got %s, want %s", data, content)
	}
	if resp.Header().Get("name") != "golang" {
		t.Errorf("got name %s, want %s", resp.Header().Get("name"), "golang")
	}
	if resp.Header().Get("X-Forwarded-For") != "8.8.8.8" {
		t.Errorf("got X-Forwarded-For %s, want %s", resp.Header().Get("X-Forwarded-For"), "8.8.8.8")
	}
}

func TestRequest_makeBody(t *testing.T) {
	req := NewRequest()
	var r io.Reader
	r = req.makeBody(nil)
	if r != nil {
		t.Fatal("got reader not nil, want nil")
	}
	s := "name=golang"
	r = req.makeBody(s)
	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if s != string(out) {
		t.Errorf("got %s, want %s", out, s)
	}
	b := []byte(s)
	r = req.makeBody(b)
	out, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if s != string(out) {
		t.Errorf("got %s, want %s", out, s)
	}
	v := url.Values{}
	v.Add("name", "golang")
	r = req.makeBody(v)
	out, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if s != string(out) {
		t.Errorf("got %s, want %s", out, s)
	}
	r = req.makeBody(strings.NewReader(s))
	out, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if s != string(out) {
		t.Errorf("got %s, want %s", out, s)
	}
	err = goutil.PanicToError(func() {
		r = req.makeBody(1)
	})
	if err == nil {
		t.Errorf("got err is nil, want err is not nil")
	}
}

func TestRequest_makeURLWithParams(t *testing.T) {
	req := NewRequest()
	baseURL := "https://golang.org"
	u := req.makeURLWithParams(baseURL, nil)
	if u != baseURL {
		t.Errorf("got %s, want %s", u, baseURL)
	}
	data := url.Values{}
	data.Set("name", "golang")
	u = req.makeURLWithParams(baseURL, data)
	expected := baseURL + "?name=golang"
	if u != expected {
		t.Errorf("got %s, want %s", u, expected)
	}
	data = url.Values{}
	data.Set("name", "golang")
	baseURL += "?"
	u = req.makeURLWithParams(baseURL, data)
	if u != expected {
		t.Errorf("got %s, want %s", u, expected)
	}
}

func TestRequest_shouldRetry(t *testing.T) {
	resp := new(http.Response)
	err := errors.New("should retry")

	req := NewRequest()
	if !req.shouldRetry(resp, err) {
		t.Errorf("got no retry, should retry")
	}

	err = nil
	resp.StatusCode = http.StatusNotFound
	if !req.shouldRetry(resp, err) {
		t.Errorf("got no retry, should retry")
	}

	err = nil
	resp.StatusCode = http.StatusOK
	if req.shouldRetry(resp, err) {
		t.Errorf("got retry, should no retry")
	}

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
	if err != nil {
		t.Fatal(err)
	}
	if resp.Header().Get("name") != "golang" {
		t.Errorf("want name is golang")
	}
	if resp.Header().Get("From") != "Post" {
		t.Errorf("want From is post")
	}
	if resp.Header().Get("X-Forwared-For") != "8.8.8.8" {
		t.Errorf("want X-Forwared-For is 8.8.8.8")
	}
}

func TestRequest_PostJSON(t *testing.T) {
	jsonString := `{"code":200,"message":"","data":{"name":"golang"}}`
	handler := func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Content-Type") != "application/json" {
			panic("invalid content-type")
		}
		io.Copy(rw, req.Body)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.PostJSON(s.URL, jsonString, nil)
	if err != nil {
		t.Fatal(err)
	}
	body, err := resp.String()
	if err != nil {
		t.Fatal(err)
	}
	if body != jsonString {
		t.Errorf("got %s, want %s", body, jsonString)
	}
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
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsStatusOK() {
		t.Errorf("got status code %d, want %d", resp.Raw().StatusCode, http.StatusNotFound)
	}
	if retryTimes != -1 {
		t.Errorf("got retrytimes %d, want %d", retryTimes, -1)
	}

	retryTimes = 3
	req = NewRequest(WithRetryTime(retryTimes))
	resp, err = req.Get(s.URL, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsStatusOK() {
		t.Errorf("got status code %d, want %d", resp.Raw().StatusCode, http.StatusOK)
	}

	if retryTimes != 1 {
		t.Errorf("got retrytimes %d, want %d", retryTimes, 1)
	}

	retryTimes = 1
	req = NewRequest(WithRetryTime(retryTimes))
	resp, err = req.Get(s.URL, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Raw().StatusCode != http.StatusNotFound {
		t.Errorf("got status code %d, want %d", resp.Raw().StatusCode, http.StatusNotFound)
	}

	if retryTimes != -1 {
		t.Errorf("got retrytimes %d, want %d", retryTimes, -1)
	}
}
