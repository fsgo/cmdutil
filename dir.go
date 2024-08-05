// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/12/21

package cmdutil

import (
	"os"
)

// MustChdir 进入指定目录，若失败会 panic
func MustChdir(to string) *Chdir {
	c, err := NewChdir(to)
	if err == nil {
		return c
	}
	panic(err)
}

// NewChdir 进入指定目录
func NewChdir(to string) (*Chdir, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	err = os.Chdir(to)
	if err != nil {
		return nil, err
	}
	return &Chdir{
		last: wd,
	}, nil
}

// Chdir 支持回退到之前目录的当前目录切换功能
type Chdir struct {
	last string
}

// GoBack 回到初始目录
func (c *Chdir) GoBack() error {
	return os.Chdir(c.last)
}

// MustGoBack 回到初始目录，若失败会 panic
func (c *Chdir) MustGoBack() {
	err := c.GoBack()
	if err == nil {
		return
	}
	panic(err)
}

// DirPushd 类似 pushd、popd 命令的目录切换功能
type DirPushd struct {
	list []*Chdir
}

func (dp *DirPushd) MustPushd(dir string) {
	err := dp.Pushd(dir)
	if err != nil {
		panic(err)
	}
}

// Pushd 进入一个目录，并将当前目录入栈
func (dp *DirPushd) Pushd(dir string) error {
	c, err := NewChdir(dir)
	if err != nil {
		return err
	}
	dp.list = append(dp.list, c)
	return nil
}

func (dp *DirPushd) MustPopd() {
	err := dp.Popd()
	if err != nil {
		panic(err)
	}
}

// Popd 返回前一个目录，若目录栈为空则直接返回 nil
func (dp *DirPushd) Popd() error {
	if len(dp.list) == 0 {
		return nil
	}
	tail := dp.list[len(dp.list)-1]
	dp.list = dp.list[0 : len(dp.list)-1]
	return tail.GoBack()
}
