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
	"sync/atomic"
	"time"

	"golang.org/x/mod/semver"

	"github.com/fsgo/cmdutil"
)

// SDK 查找当前机器的 Go SDK 情况
type SDK struct {
	ExtDirs  []string // 除了 ~/sdk/ 其他的 go sdk 根目录，可选
	inPathGo string   // "go" 二进制程序的文件路径
	list     []string
	listEnv  []*goEnv
	once     sync.Once
}

// List 返回当前机器所有安装的，可用的 'go' 的地址
// 版本由高到低排序
// 若没有，会返回空
func (gs *SDK) List() []string {
	gs.doOnce()
	return gs.list
}

// Find 查找指定号的 go 命令的地址,若查找不到会返回空字符串
//
// version: go 版本号，如 1.21
//
// 返回如 /home/work/sdk/go1.21.12/bin/go
func (gs *SDK) Find(version string) string {
	gs.doOnce()
	if !strings.HasPrefix(version, "go") {
		version = "go" + version
	}
	for _, e := range gs.listEnv {
		if strings.HasPrefix(e.version, version) {
			return e.binPath
		}
	}
	return ""
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

	scanSDKDir := func(dir string) {
		ms, _ := filepath.Glob(filepath.Join(dir, "go1.*"))
		for i := 0; i < len(ms); i++ {
			gb := filepath.Join(ms[i], "bin", "go")
			if e := gs.findGo(gb); e != nil {
				all[e.binPath] = e
			}
		}
	}
	home, err := os.UserHomeDir()
	if err == nil {
		scanSDKDir(filepath.Join(home, "sdk"))
	}

	for _, dir := range gs.ExtDirs {
		scanSDKDir(dir)
	}

	listEnv := make([]*goEnv, 0, len(all))
	for _, e := range all {
		listEnv = append(listEnv, e)
	}
	sort.SliceStable(listEnv, func(i, j int) bool {
		return listEnv[i].greater(listEnv[j])
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

// LatestOrDefault 返回最新版本 "go" 二进制文件的路径，或者是 $PATH 里的 go 版本
// 若没有，也会返回 "go"
func (gs *SDK) LatestOrDefault() string {
	l := gs.Latest()
	if len(l) != 0 {
		return l
	}
	return gs.Default()
}

// Default 返回 $PATH 里的 "go" 二进制文件的路径，若不存在，会返回空
func (gs *SDK) Default() string {
	gs.doOnce()
	return gs.inPathGo
}

// Latest 返回最新版本的 "go" 的路径，若不存在，会返回空
func (gs *SDK) Latest() string {
	gs.doOnce()
	if len(gs.listEnv) >= 0 {
		return gs.listEnv[0].binPath
	}
	return ""
}

type goEnv struct {
	binPath string // 完整的 go 命令的路径
	version string // 版本号，如 1.21，1.22.1
}

func (gs *goEnv) greater(b *goEnv) bool {
	av := "v" + gs.version[2:]
	bv := "v" + b.version[2:]
	return semver.Compare(av, bv) >= 0
}

var defaultSDK atomic.Pointer[SDK]

func init() {
	Update()
}

// Update 更新默认的环境信息
func Update() {
	defaultSDK.Store(&SDK{})
}

// DefaultOrLatest 查找 $PATH 里的 "go" 二进制文件的路径 或者是最高版本的 go
// 若没有，也会返回 "go"
func DefaultOrLatest() string {
	return defaultSDK.Load().DefaultOrLatest()
}

// LatestOrDefault 返回最新版本，或者是 $PATH 里的 go 版本,
// 若没有，也会返回 "go"
func LatestOrDefault() string {
	return defaultSDK.Load().LatestOrDefault()
}

// Latest 返回最新版本的 Go 的路径，若不存在，会返回空
func Latest() string {
	return defaultSDK.Load().Latest()
}

// Default 返回 $PATH 里的 "go" 二进制文件的路径，若不存在，会返回空
func Default() string {
	return defaultSDK.Load().Default()
}

// List 返回当前机器所有安装的，可用的 'go' 的地址
// 版本由高到低排序
// 若没有，会返回空
func List() []string {
	return defaultSDK.Load().List()
}

// GoCmdEnv 根据 goBin 路径返回设置了 GOROOT 的环境变量
func GoCmdEnv(goBin string, env []string) []string {
	if len(env) == 0 {
		env = os.Environ()
	}
	ab, err := filepath.Abs(goBin)
	if err != nil {
		return env
	}

	goBinDir := filepath.Dir(ab)
	oe := &cmdutil.OSEnv{}
	oe.WithEnviron(env)
	_ = oe.Insert("PATH", goBinDir)

	goRoot := filepath.Dir(goBinDir)
	name := filepath.Join(goRoot, "api", "go1.txt")
	info, err := os.Stat(name)
	if err != nil || info.IsDir() {
		return env
	}
	_ = oe.Set("GOROOT", goRoot)
	return os.Environ()
}
