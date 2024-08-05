// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/8/5

package gosdk

import (
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"

	"github.com/fsgo/cmdutil"
)

// RunGo 执行 go 命令
//
// root: Go SDK 的目录，如 ~/sdk/go1.22.5/
func RunGo(root string) {
	goBin := filepath.Join(root, "bin", "go"+cmdutil.Exe())
	cmd := exec.Command(goBin, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	oe := &cmdutil.OSEnv{}
	oe.MustSet("GOROOT", root)
	oe.MustInsert("PATH", filepath.Join(root, "bin"))

	cmd.Env = oe.Environ()

	handleSignals()

	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func handleSignals() {
	signal.Notify(make(chan os.Signal, 1), signalsToIgnore...)
}
