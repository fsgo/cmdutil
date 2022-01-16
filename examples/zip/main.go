// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/16

package main

import (
	"archive/zip"
	"fmt"
	"log"

	"github.com/fsgo/cmdutils"
)

func main() {
	t := &cmdutils.Zip{
		StripComponents: 0,
		IgnoreFailed:    false,
		MinSize:         1,
		UnpackNextBefore: func(f *zip.File) (skip bool, err error) {
			log.Println(f.Name)
			return false, nil
		},
	}
	err := t.Unpack("a.zip", "./tmp/")
	fmt.Println("err:", err)
}
