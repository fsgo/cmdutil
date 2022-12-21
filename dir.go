// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/12/21

package cmdutil

import (
	"os"
)

func MustChdir(to string) *Chdir {
	c, err := NewChdir(to)
	if err == nil {
		return c
	}
	panic(err)
}

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

type Chdir struct {
	last string
}

func (c *Chdir) GoBack() error {
	return os.Chdir(c.last)
}

func (c *Chdir) MustGoBack() {
	err := c.GoBack()
	if err == nil {
		return
	}
	panic(err)
}

type DirPushd struct {
	list []*Chdir
}

func (dp *DirPushd) MustPushd(dir string) {
	err := dp.Pushd(dir)
	if err != nil {
		panic(err)
	}
}

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

func (dp *DirPushd) Popd() error {
	if len(dp.list) == 0 {
		return nil
	}
	tail := dp.list[len(dp.list)-1]
	dp.list = dp.list[0 : len(dp.list)-1]
	return tail.GoBack()
}
