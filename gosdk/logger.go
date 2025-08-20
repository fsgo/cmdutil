//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-08-20

package gosdk

import (
	"io"
	"log"
	"sync/atomic"
)

var logger atomic.Value

func init() {
	lg := log.New(io.Discard, "", log.LstdFlags)
	logger.Store(lg)
}

func SetLogger(l *log.Logger) {
	if l == nil {
		l = log.New(io.Discard, "", log.LstdFlags)
	}
	logger.Store(l)
}

func getLogger() *log.Logger {
	return logger.Load().(*log.Logger)
}
