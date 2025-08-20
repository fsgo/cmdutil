// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/7

package gosdk

import (
	"context"
	"debug/buildinfo"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/mod/semver"

	"github.com/fsgo/cmdutil"
)

// SDK 查找当前机器的 Go SDK 情况
type SDK struct {
	ExtDirs    []string // 除了 ~/sdk/ 其他的 go sdk 根目录，可选
	inPathGo   string   // "go" 二进制程序的文件路径
	list       []string
	listEnv    []*goEnv
	once       sync.Once
	goEnvCache sync.Map
}

// List 返回当前机器所有安装的，可用的 'go' 的地址
// 版本由高到低排序
// 若没有，会返回空
func (gs *SDK) List(ctx context.Context) []string {
	gs.doOnce(ctx)
	return gs.list
}

// Find 查找指定号的 go 命令的地址,若查找不到会返回空字符串
//
// version: go 版本号，如 1.21
//
// 返回如 /home/work/sdk/go1.21.12/bin/go
func (gs *SDK) Find(ctx context.Context, version string) string {
	gs.doOnce(ctx)
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

func (gs *SDK) doOnce(ctx context.Context) {
	gs.once.Do(func() {
		gs.findList(ctx)
	})
}

func (gs *SDK) findList(ctx context.Context) {
	all := make(map[string]*goEnv)
	if e := gs.findGo(ctx, "go"); e != nil {
		gs.inPathGo = e.binPath
		all[e.binPath] = e
	}

	scanSDKDir := func(dir string) {
		ms, _ := filepath.Glob(filepath.Join(dir, "go1.*"))
		for i := 0; i < len(ms); i++ {
			gb := filepath.Join(ms[i], "bin", "go")
			if e := gs.findGo(ctx, gb); e != nil {
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

func (gs *SDK) findGo(ctx context.Context, binPath string) *goEnv {
	getLogger().Println("SDK.findGo, binPath=", binPath)
	binPath = filepath.Clean(binPath)

	value, ok := gs.goEnvCache.Load(binPath)
	if ok {
		if result, ok1 := value.(*goEnv); ok1 {
			return result
		}
		return nil
	}
	bp, err := exec.LookPath(binPath)
	if err != nil {
		return nil
	}

	bi, err := buildinfo.ReadFile(bp)
	if err != nil {
		return nil
	}
	bf, _ := json.Marshal(bi)
	getLogger().Printf("SDK.findGo, bp=%q buildInfo=%#v, err=%v\n", bp, bf, err)
	if bi.Path != "cmd/go" {
		return nil
	}

	result := &goEnv{
		binPath: bp,
		version: bi.GoVersion,
	}
	gs.goEnvCache.Store(binPath, result)
	return result
}

// DefaultOrLatest 查找 $PATH 里的 go 或者是最高版本的 go
// 若没有，也会返回 "go"
func (gs *SDK) DefaultOrLatest(ctx context.Context) string {
	gs.doOnce(ctx)
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
func (gs *SDK) LatestOrDefault(ctx context.Context) string {
	l := gs.Latest(ctx)
	if len(l) != 0 {
		return l
	}
	return gs.Default(ctx)
}

// Default 返回 $PATH 里的 "go" 二进制文件的路径，若不存在，会返回空
func (gs *SDK) Default(ctx context.Context) string {
	gs.doOnce(ctx)
	return gs.inPathGo
}

// Latest 返回最新版本的 "go" 的路径，若不存在，会返回空
func (gs *SDK) Latest(ctx context.Context) string {
	gs.doOnce(ctx)
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
func DefaultOrLatest(ctx context.Context) string {
	return defaultSDK.Load().DefaultOrLatest(ctx)
}

// LatestOrDefault 返回最新版本，或者是 $PATH 里的 go 版本,
// 若没有，也会返回 "go"
func LatestOrDefault(ctx context.Context) string {
	return defaultSDK.Load().LatestOrDefault(ctx)
}

// Latest 返回最新版本的 Go 的路径，若不存在，会返回空
func Latest(ctx context.Context) string {
	return defaultSDK.Load().Latest(ctx)
}

// Default 返回 $PATH 里的 "go" 二进制文件的路径，若不存在，会返回空
func Default(ctx context.Context) string {
	return defaultSDK.Load().Default(ctx)
}

// List 返回当前机器所有安装的，可用的 'go' 的地址
// 版本由高到低排序
// 若没有，会返回空
func List(ctx context.Context) []string {
	return defaultSDK.Load().List(ctx)
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
