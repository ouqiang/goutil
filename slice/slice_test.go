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

	"github.com/stretchr/testify/require"
)

func TestContainsString(t *testing.T) {
	a := []string{"foo", "test", "example"}
	require.True(t, ContainsString(a, "example"))
	require.False(t, ContainsString(a, "none"))
}

func TestContainsInt(t *testing.T) {
	a := []int{100, 1098, 2890}
	require.True(t, ContainsInt(a, 2890))
	require.False(t, ContainsInt(a, 9000))
}

func TestRemoveEmptyString(t *testing.T) {
	a := []string{"foo", "test", "", "example", ""}
	expected := []string{"foo", "test", "example"}
	require.Equal(t, expected, RemoveEmptyString(a))

	a = []string{"foo", "test", "example"}
	require.Equal(t, a, RemoveEmptyString(a))
}
