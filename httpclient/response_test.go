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
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestResponse_IsStatusOK(t *testing.T) {
	statusCode := http.StatusOK
	handler := func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(statusCode)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsStatusOK() {
		t.Errorf("got status code %d, want %d", resp.Raw().StatusCode, statusCode)
	}

	statusCode = http.StatusNotFound
	resp, err = req.Get(s.URL, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.IsStatusOK() {
		t.Errorf("got status code %d, want %d", resp.Raw().StatusCode, statusCode)
	}
}

func TestResponse_DecodeJSON(t *testing.T) {
	jsonString := `{"code":200,"message":"","data":{"name":"golang"}}`
	handler := func(rw http.ResponseWriter, req *http.Request) {
		io.WriteString(rw, jsonString)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	var apiResponse struct {
		Code    int
		Message string
		Data    struct {
			Name string
		}
	}
	err = resp.DecodeJSON(&apiResponse)
	if err != nil {
		t.Fatal(err)
	}
	if apiResponse.Code != 200 {
		t.Errorf("got code %d, want %d", apiResponse.Code, 200)
	}
}

func TestResponse_String(t *testing.T) {
	content := "string"
	handler := func(rw http.ResponseWriter, req *http.Request) {
		io.WriteString(rw, content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
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
}

func TestResponse_Bytes(t *testing.T) {
	content := []byte("bytes")
	handler := func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	data, err := resp.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, content) {
		t.Errorf("got %s, want %s", data, content)
	}
}

func TestResponse_Discard(t *testing.T) {
	content := []byte("discard")
	handler := func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, err = resp.Discard()
	if err != nil {
		t.Fatal(err)
	}
}

func TestResponse_WriteFile(t *testing.T) {
	content := []byte("write file")
	handler := func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	tmpFilename := tmpFile.Name()
	tmpFile.Close()
	_, err = resp.WriteFile(tmpFilename)
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadFile(tmpFilename)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(content, data) {
		t.Errorf("got %s, want %s", data, content)
	}
}

func ExampleResponse_WriteTo() {
	content := []byte("write file")
	handler := func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	if err != nil {
		panic(err)
	}
	_, err = resp.WriteTo(os.Stdout)
	if err != nil {
		panic(err)
	}
	// OUTPUT: write file
}
