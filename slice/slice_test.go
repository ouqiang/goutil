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

package slice

import (
	"testing"

	"reflect"

	"fmt"

	"github.com/ouqiang/goutil"
)

func TestCheckSlice(t *testing.T) {
	items := [...]interface{}{1, "golang", [...]int{1, 2, 3}, func() {}}
	var err error
	for _, item := range items {
		err = goutil.PanicToError(func() {
			checkSlice(item)
		})
		if err == nil {
			t.Fatal("got err is nil, want err is not nil")
		}
	}
	err = goutil.PanicToError(func() {
		var arr []int
		checkSlice(arr)
	})
	if err != nil {
		t.Fatalf("got err [%s], want err is nil", err)
	}
}

func TestCheckFunc(t *testing.T) {
	items := [...]interface{}{1, "golang", []int{1, 2, 3}}
	var err error
	for _, item := range items {
		err = goutil.PanicToError(func() {
			checkFunc(item)
		})
		if err == nil {
			t.Fatal("got err is nil, want err is not nil")
		}
	}
	err = goutil.PanicToError(func() {
		f := func() {}
		checkFunc(f)
	})
	if err != nil {
		t.Fatalf("got err [%s], want err is nil", err)
	}
}

func TestTypeError_Error(t *testing.T) {
	e := &TypeError{"arr", "is not slice"}
	expected := "arr is not slice"
	if e.Error() != expected {
		t.Errorf("got %s, want %s", e.Error(), expected)
	}
}

func TestEvery(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	f := func(currentValue interface{}, index int, l interface{}) bool {
		v := currentValue.(int)

		if arr[index] != v {
			t.Fatalf("got %d, want %d", arr[index], v)
		}
		if !reflect.DeepEqual(arr, l) {
			t.Fatalf("got not equal, want equal")
		}

		return v > 0
	}

	if !Every(arr, f) {
		t.Error("got false, want true")
	}

	arr = []int{0, 1, 2, 3, 4, 5}
	if Every(arr, f) {
		t.Error("got true, want false")
	}
}

func TestFilter(t *testing.T) {
	arr := []int{0, 1, 2, 3, 4, 5}
	f := func(currentValue interface{}, index int, arr interface{}) bool {
		v := currentValue.(int)
		return v >= 3
	}

	out := Filter(arr, f)
	expected := []int{3, 4, 5}
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("got %+v, want %+v", out, expected)
	}
}

func TestMap(t *testing.T) {
	arr := []int{0, 1, 2, 3, 4, 5}
	f := func(currentValue interface{}, index int, arr interface{}) interface{} {
		v := currentValue.(int)
		return v + 1
	}
	out := Map(arr, f)
	expected := []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("got %+v, want %+v", out, expected)
	}
}

func ExampleForEach() {
	arr := []int{0, 1, 2}
	f := func(currentValue interface{}, index int, arr interface{}) bool {
		fmt.Println(currentValue)
		return true
	}
	ForEach(arr, f)

	f = func(currentValue interface{}, index int, arr interface{}) bool {
		fmt.Println(currentValue)

		return index < 1
	}
	ForEach(arr, f)
	// OUTPUT: 0
	// 1
	// 2
	// 0
	// 1
}

func TestInclude(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	if !Include(arr, 1) {
		t.Errorf("got false, want true")
	}
	if Include(arr, 10) {
		t.Errorf("got true, want false")
	}
}

func TestIndex(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	if Index(arr, 1) == ItemNotFound {
		t.Errorf("got false, want true")
	}
	if Index(arr, 10) != ItemNotFound {
		t.Errorf("got true, want false")
	}
}
