// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/8/9

package cmdutil

import (
	"sync"
)

type WorkerGroup struct {
	// Max 最大并发度，可选，默认为 1
	Max int

	once    sync.Once
	limiter chan struct{}
	wait    sync.WaitGroup
}

func (wg *WorkerGroup) init() {
	wg.once.Do(func() {
		wg.limiter = make(chan struct{}, max(wg.Max, 1))
	})
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
