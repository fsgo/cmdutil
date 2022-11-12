// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/12

package cmdutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOSEnv1(t *testing.T) {
	t.Setenv("PATH", "hello1:world1")
	de := &OSEnv{}
	key1 := "PATH"
	path := os.Getenv(key1)

	key2 := "PATH2022"
	path2022 := os.Getenv(key2)
	require.Empty(t, path2022)

	require.NoError(t, de.Insert(key1, "abc"))
	require.Equal(t, "abc:"+path, de.Get(key1))

	require.NoError(t, de.Insert(key1, "hello"))
	require.Equal(t, "hello:abc:"+path, de.Get(key1))

	require.NoError(t, de.Append(key1, "world"))
	require.Equal(t, "hello:abc:"+path+":world", de.Get(key1))

	require.NoError(t, de.Set(key1, "world"))
	require.Equal(t, "world", de.Get(key1))

	require.NoError(t, de.Delete(key1))
	require.Equal(t, "", de.Get(key1))

	require.NoError(t, de.Delete(key2))
	require.NoError(t, de.Append(key2, "world"))
	require.Equal(t, "world", de.Get(key2))

	require.NoError(t, de.Delete(key2))
	require.Equal(t, "", de.Get(key2))

	require.NoError(t, de.Insert(key2, "world"))
	require.Equal(t, "world", de.Get(key2))

	require.NoError(t, de.Delete(key2))
	require.Equal(t, "", de.Get(key2))

	require.NoError(t, de.Set(key2, "world"))
	require.Equal(t, "world", de.Get(key2))

	require.Error(t, de.Insert("", "abc"))
	require.Error(t, de.Append("", "abc"))
	require.Error(t, de.Set("", "abc"))
	require.Error(t, de.Delete(""))
	require.Empty(t, de.Get(""))
}

func TestOSEnv2(t *testing.T) {
	de := &OSEnv{}
	de.WithEnviron([]string{"P6", "PATH=abc", "P2", "P3", "P4", "P5"})
	require.Equal(t, "abc", de.Get("PATH"))
	require.Equal(t, "", de.Get("P2"))

	require.NoError(t, de.Delete("P5"))

	require.NoError(t, de.Set("P4", "v4"))
	require.Equal(t, "v4", de.Get("P4"))

	require.NoError(t, de.Insert("P2", "v2"))
	require.Equal(t, "v2", de.Get("P2"))

	require.NoError(t, de.Append("P3", "v3"))
	require.Equal(t, "v3", de.Get("P3"))

	all := []string{"P6", "PATH=abc", "P2=v2", "P3=v3", "P4=v4"}
	require.Equal(t, all, de.Environ())
}

func TestOSEnv3(t *testing.T) {
	old := caseInsensitiveEnv
	defer func() {
		caseInsensitiveEnv = old
	}()
	caseInsensitiveEnv = true
	de := &OSEnv{}
	de.WithEnviron([]string{"P6", "PATH=abc", "P2", "P3", "P4", "P5"})
	require.Equal(t, "abc", de.Get("path"))
	require.Equal(t, "", de.Get("p2"))

	require.NoError(t, de.Delete("p5"))

	require.NoError(t, de.Set("p4", "v4"))
	require.Equal(t, "v4", de.Get("P4"))

	require.NoError(t, de.Insert("p2", "v2"))
	require.Equal(t, "v2", de.Get("P2"))

	require.NoError(t, de.Append("p3", "v3"))
	require.Equal(t, "v3", de.Get("P3"))

	all := []string{"P6", "PATH=abc", "p2=v2", "p3=v3", "p4=v4"}
	require.Equal(t, all, de.Environ())
}

func TestOSEnv_unique(t *testing.T) {
	oe := &OSEnv{}
	require.Equal(t, "abc", oe.unique("abc"))
	require.Equal(t, "abc:def", oe.unique("abc:def:abc"))
	require.Equal(t, "abc:def", oe.unique("abc:def:abc ::"))
}
