// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/16

package main

import (
	"archive/tar"
	"fmt"
	"log"

	"github.com/fsgo/cmdutil"
)

func main() {
	t := &cmdutil.Tar{
		StripComponents: 0,
		IgnoreFailed:    false,
		MinSize:         1,
		UnpackNextBefore: func(h *tar.Header) (skip bool, err error) {
			log.Println(h.Name, h.Size)
			return false, nil
		},
	}
	err := t.Unpack("a.tar.gz", "./tmp/")
	fmt.Println("err:", err)
}
