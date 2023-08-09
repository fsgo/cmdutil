// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/8/9

package cmdutil

import (
	"sync"
	"sync/atomic"
)

type WorkerGroup struct {
	// Max 最大并发度，默认为 1
	Max int

	once    atomic.Bool
	limiter chan struct{}
	wait    sync.WaitGroup
}

func (wg *WorkerGroup) init() {
	if wg.once.CompareAndSwap(false, true) {
		var m = wg.Max
		if m <= 0 {
			m = 1
		}
		wg.limiter = make(chan struct{}, m)
	}
}

func (wg *WorkerGroup) Run(fn func()) {
	wg.init()
	wg.limiter <- struct{}{}
	wg.wait.Add(1)
	go func() {
		defer func() {
			<-wg.limiter
			wg.wait.Done()
		}()
		fn()
	}()
}

func (wg *WorkerGroup) Wait() {
	wg.wait.Wait()
}
