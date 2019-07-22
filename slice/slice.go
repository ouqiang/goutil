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

// ContainsString 是否包含字符串元素
func ContainsString(a []string, s string) bool {
	for _, item := range a {
		if item == s {
			return true
		}
	}

	return false
}

// ContainsInt 是否包含整型元素
func ContainsInt(a []int, i int) bool {
	for _, item := range a {
		if item == i {
			return true
		}
	}

	return false
}

// RemoveEmptyString 移除空字符串
func RemoveEmptyString(a []string) []string {
	length := len(a)
	if length == 0 {
		return a
	}
	out := make([]string, 0, length)
	for _, item := range a {
		if item != "" {
			out = append(out, item)
		}
	}
	if len(out) == length {
		return a
	}

	return out
}

// Remove 移除元素
func Remove(a []string, s string) []string {
        length := len(a)
        if length == 0 {
                return a
        }
        out := make([]string, 0, length)
        for _, item := range a {
                if item != s {
                        out = append(out, item)
                }
        }

        return out
}

