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

// Package slice 数组常用操作
package slice

import (
	"fmt"
	"reflect"
)

const (
	// ItemNotFound 元素未找到
	ItemNotFound = -1
)

// TypeError 类型错误
type TypeError struct {
	typeName string
	msg      string
}

func (te *TypeError) Error() string {
	return fmt.Sprintf("%s %s", te.typeName, te.msg)
}

// ItemFunc 对元素进行处理
type ItemFunc func(currentValue interface{}, index int, arr interface{})

// ItemFuncFilter 判断元素是否满足要求
type ItemFuncFilter func(currentValue interface{}, index int, arr interface{}) bool

// ItemFuncMap 对元素进行处理, 并返回新值
type ItemFuncMap func(currentValue interface{}, index int, arr interface{}) interface{}

// Every 判断数组元素是否都满足f
func Every(arr interface{}, f ItemFuncFilter) bool {
	checkFunc(f)
	elem, length := reflectSlice(arr)
	if length == 0 {
		return false
	}
	for i := 0; i < length; i++ {
		if !f(elem.Index(i).Interface(), i, arr) {
			return false
		}
	}

	return true
}

// Filter 对数组元素过滤, 返回过滤后的数组
func Filter(arr interface{}, f ItemFuncFilter) interface{} {
	checkFunc(f)
	elem, length := reflectSlice(arr)
	if length == 0 {
		return arr
	}
	out := reflect.MakeSlice(reflect.TypeOf(arr), 0, length)
	for i := 0; i < length; i++ {
		if f(elem.Index(i).Interface(), i, arr) {
			out = reflect.Append(out, elem.Index(i))
		}
	}

	return out.Slice(0, out.Len()).Interface()
}

// Map 对数组每个元素执行f, 并返回新值, 不会修改原数组
func Map(arr interface{}, f ItemFuncMap) interface{} {
	checkFunc(f)
	elem, length := reflectSlice(arr)
	if length == 0 {
		return arr
	}
	out := reflect.MakeSlice(reflect.TypeOf(arr), 0, length)
	for i := 0; i < length; i++ {
		v := f(elem.Index(i).Interface(), i, arr)
		out = reflect.Append(out, reflect.ValueOf(v))
	}

	return out.Interface()
}

// ForEach 对数组每个元素进行处理
func ForEach(arr interface{}, f ItemFuncFilter) {
	checkFunc(f)
	elem, length := reflectSlice(arr)
	for i := 0; i < length; i++ {
		if !f(elem.Index(i).Interface(), i, arr) {
			return
		}
	}
}

// Include arr中是否包含value
func Include(arr interface{}, value interface{}) bool {
	index := Index(arr, value)

	return index != ItemNotFound
}

// Index arr中存在value, 返回索引, 不存在返回-1
func Index(arr interface{}, value interface{}) int {
	elem, length := reflectSlice(arr)
	var found bool
	for i := 0; i < length; i++ {
		switch elem.Kind() {
		case reflect.String:
			valueString, ok := value.(string)
			if !ok {
				panic("required value is string type")
			}
			found = elem.Index(i).String() == valueString
		default:
			found = reflect.DeepEqual(elem.Index(i).Interface(), value)
		}
		if found {
			return i
		}
	}

	return ItemNotFound
}

func reflectSlice(arr interface{}) (reflect.Value, int) {
	checkSlice(arr)
	refValue := reflect.ValueOf(arr)
	length := refValue.Len()

	return refValue, length
}

func checkSlice(arr interface{}) {
	refType := reflect.TypeOf(arr)
	if refType.Kind() != reflect.Slice {
		panic(&TypeError{refType.Name(), "is not slice"})
	}
}

func checkFunc(f interface{}) {
	refType := reflect.TypeOf(f)
	if refType.Kind() != reflect.Func {
		panic(&TypeError{refType.Name(), "is not func"})
	}
}
