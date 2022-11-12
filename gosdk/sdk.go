// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/7

package gosdk

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/mod/semver"
)

// SDK 查找当前机器的 Go SDK 情况
type SDK struct {
	inPathGo string
	list     []string
	listEnv  []*goEnv
	once     sync.Once
}

// List 返回当前机器所有安装的，可用的 'go'的地址
// 安装版本由高到低排序
// 若没有，会返回空
func (gs *SDK) List() []string {
	gs.doOnce()
	return gs.list
}

func (gs *SDK) doOnce() {
	gs.once.Do(gs.findList)
}

func (gs *SDK) findList() {
	all := make(map[string]*goEnv)
	if e := gs.findGo("go"); e != nil {
		gs.inPathGo = e.binPath
		all[e.binPath] = e
	}
	home, err := os.UserHomeDir()
	if err == nil {
		ms, _ := filepath.Glob(filepath.Join(home, "sdk", "go1.*"))
		for i := 0; i < len(ms); i++ {
			gb := filepath.Join(ms[i], "bin", "go")
			if e := gs.findGo(gb); e != nil {
				all[e.binPath] = e
			}
		}
	}
	listEnv := make([]*goEnv, 0, len(all))
	for _, e := range all {
		listEnv = append(listEnv, e)
	}
	sort.SliceStable(listEnv, func(i, j int) bool {
		return listEnv[i].Greater(listEnv[j])
	})
	gs.listEnv = listEnv

	list := make([]string, 0, len(listEnv))
	for _, e := range listEnv {
		list = append(list, e.binPath)
	}
	gs.list = list
}

func (gs *SDK) findGo(binPath string) *goEnv {
	bp, err := exec.LookPath(binPath)
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bp, "version")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	str := string(out)
	if !strings.HasPrefix(str, "go version go") {
		return nil
	}
	fields := strings.Fields(str)
	version := fields[2]
	return &goEnv{
		binPath: bp,
		version: version,
	}
}

// DefaultOrLatest 查找 $PATH 里的 go 或者是最高版本的 go
// 若没有，也会返回 "go"
func (gs *SDK) DefaultOrLatest() string {
	gs.doOnce()
	if len(gs.inPathGo) != 0 {
		return gs.inPathGo
	}
	if len(gs.listEnv) >= 0 {
		return gs.listEnv[0].binPath
	}
	return "go"
}

// Default 返回 $PATH 里的 Go 的路径
func (gs *SDK) Default() string {
	gs.doOnce()
	return gs.inPathGo
}

// Latest 返回最新版本的 Go 的路径，若不存在，会返回空
func (gs *SDK) Latest() string {
	gs.doOnce()
	if len(gs.listEnv) >= 0 {
		return gs.listEnv[0].binPath
	}
	return ""
}

type goEnv struct {
	binPath string
	version string
}

func (gs *goEnv) Greater(b *goEnv) bool {
	av := "v" + gs.version[2:]
	bv := "v" + b.version[2:]
	return semver.Compare(av, bv) >= 0
}

var defaultSDK = &SDK{}

// DefaultOrLatest 查找 $PATH 里的 go 或者是最高版本的 go
// 若没有，也会返回 "go"
func DefaultOrLatest() string {
	return defaultSDK.DefaultOrLatest()
}

// Latest 返回最新版本的 Go 的路径，若不存在，会返回空
func Latest() string {
	return defaultSDK.Latest()
}

// Default 返回 $PATH 里的 Go 的路径
func Default() string {
	return defaultSDK.Default()
}