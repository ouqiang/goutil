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

// Package goutil
package goutil

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMD5(t *testing.T) {
	str := "goutil"
	result := MD5(str)
	expected := "3fa6e676e7d5c558e9a49b599cbc975f"
	if result != expected {
		t.Errorf("got %s, want %s", result, expected)
	}
}

func BenchmarkMD5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MD5("goutil")
	}
}

func ExampleMD5() {
	str := "goutil"
	fmt.Println(MD5(str))
	// OUTPUT: 3fa6e676e7d5c558e9a49b599cbc975f
}

func TestRandNumber(t *testing.T) {
	f := func(min, max int) {
		for i := min; i <= max; i++ {
			num := RandNumber(i, max)
			if min <= num && num <= max {
				continue
			}
			t.Errorf("got %d, want range %d - %d", num, min, max)
		}
	}
	f(0, 1000)
	f(0, 0)
	f(-1000, 0)
	f(-1000, 1000)
}

func BenchmarkRandNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandNumber(1, math.MaxInt32)
	}
}

func ExampleRandNumber() {
	num := RandNumber(1, 1000)
	fmt.Println(num)
}

func TestPanicToError(t *testing.T) {
	err := PanicToError(func() {
		var m map[string]string
		m["key"] = "value"
	})
	if err == nil {
		t.Errorf("got err is nil, want err is not nil")
	}
}

func ExamplePanicToError() {
	err := PanicToError(func() {
		var m map[string]string
		m["key"] = "value"
	})
	fmt.Println(err)
}

func TestDownloadFile(t *testing.T) {
	handler := func(rw http.ResponseWriter, req *http.Request) {
		filename := filepath.Join("testdata", "download.txt")
		err := DownloadFile(filename, rw)
		if err != nil {
			t.Fatal(err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	expected := "golang"
	if string(body) != expected {
		t.Fatalf("got file content [%s], want content [%s]", body, expected)
	}
	fields := strings.Split(resp.Header.Get("Content-Disposition"), "=")
	if len(fields) != 2 {
		t.Fatalf("unexpected download filename")
	}
	if strings.TrimSpace(fields[1]) != `"download.txt"` {
		t.Fatalf("unexpected download filename: %s", fields[1])
	}
}

func TestWorkDir(t *testing.T) {
	wd, err := WorkDir()
	if err != nil {
		t.Fatal(err)
	}
	if wd != filepath.Dir(os.Args[0]) {
		t.Fatalf("got working dir [%s], want working dir [%s]", wd, os.Args[0])
	}
}
