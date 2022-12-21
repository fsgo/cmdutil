// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/12/21

package cmdutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func checkDir(t *testing.T, want string) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.Contains(t, filepath.Base(wd), want)
}

func TestDirPushd_MustPushd(t *testing.T) {
	dp := &DirPushd{}
	dp.MustPushd("gosdk")
	checkDir(t, "gosdk")

	dp.MustPushd("../_example")
	checkDir(t, "_example")

	dp.MustPopd()
	checkDir(t, "gosdk")
	dp.MustPopd()
	checkDir(t, "cmdutil")

	// 额外多调用一次也不会有问题
	dp.MustPopd()
}

func TestMustChdir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := MustChdir("gosdk")
		checkDir(t, "gosdk")
		c.MustGoBack()
	})

	t.Run("dir not exists", func(t *testing.T) {
		defer func() {
			re := recover()
			require.NotNil(t, re)
		}()
		MustChdir("not_found")
	})
}
