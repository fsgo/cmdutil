// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/12

package cmdutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xanygo/anygo/xt"
)

func TestOSEnv1(t *testing.T) {
	const (
		key1 = "PATH"
		key2 = "PATH2022"
		sp   = string(filepath.ListSeparator)
	)
	t.Setenv(key1, "hello1"+sp+"world1")
	t.Setenv(key2, "")

	path := os.Getenv(key1)
	xt.Equal(t, "hello1"+sp+"world1", path)

	xt.Empty(t, os.Getenv(key2))

	de := &OSEnv{}
	xt.Equal(t, []string{"hello1", "world1"}, de.GetValues(key1))
	xt.Equal(t, path, de.Get(key1))

	xt.NoError(t, de.Insert(key1, "abc"))
	xt.Equal(t, "abc"+sp+path, de.Get(key1))

	xt.NoError(t, de.Insert(key1, "hello"))
	xt.Equal(t, "hello"+sp+"abc"+sp+path, de.Get(key1))

	xt.NoError(t, de.Append(key1, "world"))
	xt.Equal(t, "hello"+sp+"abc"+sp+path+sp+"world", de.Get(key1))

	xt.NoError(t, de.Set(key1, "world"))
	xt.Equal(t, "world", de.Get(key1))

	xt.NoError(t, de.Delete(key1))
	xt.Equal(t, "", de.Get(key1))

	xt.NoError(t, de.Delete(key2))
	xt.NoError(t, de.Append(key2, "world"))
	xt.Equal(t, "world", de.Get(key2))

	xt.NoError(t, de.Delete(key2))
	xt.Equal(t, "", de.Get(key2))

	xt.NoError(t, de.Insert(key2, "world"))
	xt.Equal(t, "world", de.Get(key2))

	xt.NoError(t, de.Delete(key2))
	xt.Equal(t, "", de.Get(key2))

	xt.NoError(t, de.Set(key2, "world"))
	xt.Equal(t, "world", de.Get(key2))

	xt.Error(t, de.Insert("", "abc"))
	xt.Error(t, de.Append("", "abc"))
	xt.Error(t, de.Set("", "abc"))
	xt.Error(t, de.Delete(""))
	xt.Empty(t, de.Get(""))

	xt.NoError(t, de.Set(key1, path))
	xt.NoError(t, de.DeleteValue(key1, "hello1"))
	xt.Equal(t, "world1", de.Get(key1))

	xt.NoError(t, de.DeleteValue(key1, "not-found"))
	xt.Equal(t, "world1", de.Get(key1))
}

func TestOSEnv2(t *testing.T) {
	de := &OSEnv{}
	de.WithEnviron([]string{"P6", "PATH=abc", "P2", "P3", "P4", "P5"})
	xt.Equal(t, "abc", de.Get("PATH"))
	xt.Equal(t, "", de.Get("P2"))

	xt.NoError(t, de.Delete("P5"))

	xt.NoError(t, de.Set("P4", "v4"))
	xt.Equal(t, "v4", de.Get("P4"))

	xt.NoError(t, de.Insert("P2", "v2"))
	xt.Equal(t, "v2", de.Get("P2"))

	xt.NoError(t, de.Append("P3", "v3"))
	xt.Equal(t, "v3", de.Get("P3"))

	all := []string{"P6", "PATH=abc", "P2=v2", "P3=v3", "P4=v4"}
	xt.Equal(t, all, de.Environ())
}

func TestOSEnv3(t *testing.T) {
	old := caseInsensitiveEnv
	defer func() {
		caseInsensitiveEnv = old
	}()
	caseInsensitiveEnv = true
	de := &OSEnv{}
	de.WithEnviron([]string{"P6", "PATH=abc", "P2", "P3", "P4", "P5"})
	xt.Equal(t, "abc", de.Get("path"))
	xt.Equal(t, "", de.Get("p2"))

	xt.NoError(t, de.Delete("p5"))

	xt.NoError(t, de.Set("p4", "v4"))
	xt.Equal(t, "v4", de.Get("P4"))

	xt.NoError(t, de.Insert("p2", "v2"))
	xt.Equal(t, "v2", de.Get("P2"))

	xt.NoError(t, de.Append("p3", "v3"))
	xt.Equal(t, "v3", de.Get("P3"))

	all := []string{"P6", "PATH=abc", "p2=v2", "p3=v3", "p4=v4"}
	xt.Equal(t, all, de.Environ())
}

func TestOSEnv_unique(t *testing.T) {
	const sp = string(filepath.ListSeparator)
	oe := &OSEnv{}
	xt.Equal(t, "abc", oe.unique("abc"))
	xt.Equal(t, "abc"+sp+"def", oe.unique("abc"+sp+"def"+sp+"abc"))
	xt.Equal(t, "abc"+sp+"def", oe.unique("abc"+sp+"def"+sp+"abc "+sp+sp))
}
