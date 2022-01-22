// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/16

package main

import (
	"log"
	"os"

	"github.com/fsgo/cmdutils"
)

func main() {
	w := &cmdutils.Wget{
		LogWriter: os.Stderr,
	}
	u := "https://go.dev/dl/go1.17.6.darwin-amd64.tar.gz"
	u = "http://127.0.0.1:8088/chunk?repeat=10"
	err := w.Download(u, "tmp/a.tar.gz")
	log.Println("err:", err)
}
