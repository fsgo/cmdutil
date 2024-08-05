// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/4/11

package gosdk

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
)

// MustAutoDisableGoWork AutoDisableGoWork 的快捷调用,若失败会 panic
func MustAutoDisableGoWork() {
	err := AutoDisableGoWork()
	if err != nil {
		panic(err)
	}
}

// TryAutoDisableGoWork AutoDisableGoWork 的快捷调用，若失败，会打印日志
func TryAutoDisableGoWork() {
	err := AutoDisableGoWork()
	if err != nil {
		log.Println("AutoDisableGoWork failed:", err)
	}
}

// AutoDisableGoWork 自动禁用 go work 功能
//
// 若当前项目未在 go.work 文件中定义，则设置环境变量 GOWORK = off
func AutoDisableGoWork() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, DefaultOrLatest(), "env", "GOWORK")
	cmd.Stdin = os.Stdin
	out, err0 := cmd.Output()
	if err0 != nil {
		return err0
	}
	fp := strings.TrimSpace(string(out))
	if fp == "off" || fp == "" {
		return nil
	}
	bf, err1 := os.ReadFile(fp)
	if err1 != nil {
		return err1
	}
	wf, err2 := modfile.ParseWork(fp, bf, nil)
	if err2 != nil {
		return err2
	}

	wd, err3 := os.Getwd()
	if err3 != nil {
		return err3
	}
	dir := filepath.Dir(fp)

	var found bool
	for _, u := range wf.Use {
		ap := filepath.Join(dir, u.Path)
		if strings.HasPrefix(wd, ap) {
			found = true
			break
		}
	}
	if !found {
		return os.Setenv("GOWORK", "off")
	}
	return nil
}
