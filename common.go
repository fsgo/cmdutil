// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/16

package cmdutil

import (
	"os"
	"runtime"
)

func mkdir(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil && os.IsExist(err) {
		return nil
	}
	return err
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func Exe() string {
	if isWindows() {
		return ".exe"
	}
	return ""
}
