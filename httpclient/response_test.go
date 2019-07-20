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
	"testing"

	"github.com/gogo/protobuf/proto"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	require.Equal(t, statusCode, resp.Raw().StatusCode)

	statusCode = http.StatusNotFound
	resp, err = req.Get(s.URL, nil, nil)
	require.NoError(t, err)
	require.Equal(t, statusCode, resp.Raw().StatusCode)
}

func TestResponse_DecodeJSON(t *testing.T) {
	jsonString := `{"code":200,"message":"","data":{"name":"golang"}}`
	handler := func(rw http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(rw, jsonString)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	require.NoError(t, err)
	var apiResponse struct {
		Code    int
		Message string
		Data    struct {
			Name string
		}
	}
	err = resp.DecodeJSON(&apiResponse)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, apiResponse.Code)
}

func TestResponse_DecodeProtoBuf(t *testing.T) {
	handler := func(rw http.ResponseWriter, req *http.Request) {
		message := &Message{
			Name: "protobuf",
		}
		data, err := proto.Marshal(message)
		if err != nil {
			panic(err)
		}
		_, _ = rw.Write(data)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	require.NoError(t, err)

	result := &Message{}
	err = resp.DecodeProtoBuf(result)
	require.NoError(t, err)

	require.Equal(t, "protobuf", result.Name)
}

func TestResponse_String(t *testing.T) {
	content := "string"
	handler := func(rw http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(rw, content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	require.NoError(t, err)
	data, err := resp.String()
	require.NoError(t, err)
	require.Equal(t, content, data)
}

func TestResponse_Bytes(t *testing.T) {
	content := []byte("bytes")
	handler := func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write(content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	require.NoError(t, err)
	data, err := resp.Bytes()
	require.NoError(t, err)
	require.Equal(t, content, data)
}

func TestResponse_Discard(t *testing.T) {
	content := []byte("discard")
	handler := func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write(content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	require.NoError(t, err)
	_, err = resp.Discard()
	require.NoError(t, err)
}

func TestResponse_WriteFile(t *testing.T) {
	content := []byte("write file")
	handler := func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write(content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	require.NoError(t, err)
	tmpFile, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	tmpFilename := tmpFile.Name()
	_ = tmpFile.Close()
	_, err = resp.WriteFile(tmpFilename)
	require.NoError(t, err)
	data, err := ioutil.ReadFile(tmpFilename)
	require.NoError(t, err)
	require.Equal(t, content, data)
}

type Writer struct {
	data []byte
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.data = p

	return len(p), nil
}

func (w *Writer) Bytes() []byte {
	return w.data
}

func TestResponse_WriteTo(t *testing.T) {
	content := []byte("write file")
	handler := func(rw http.ResponseWriter, req *http.Request) {
		_, _ = rw.Write(content)
	}
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	req := NewRequest()
	resp, err := req.Get(s.URL, nil, nil)
	require.NoError(t, err)

	w := &Writer{}
	_, err = resp.WriteTo(w)
	require.NoError(t, err)
	require.Equal(t, content, w.Bytes())
}
