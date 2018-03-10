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

package goutil

// Semaphore 信号量模拟, 用于控制goroutine并发数量
type Semaphore struct {
	queue chan struct{}
}

// NewSemaphore 初始化Semaphore
func NewSemaphore(size int) *Semaphore {
	sem := &Semaphore{
		queue: make(chan struct{}, size),
	}

	return sem
}

// Add 向sem发送数据
func (sem *Semaphore) Add() {
	sem.queue <- struct{}{}
}

// Done 从sem读取数据
func (sem *Semaphore) Done() {
	<-sem.queue
}
