// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/7

package gosdk

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultOrLatest(t *testing.T) {
	got := DefaultOrLatest()
	require.NotEmpty(t, got)
}

func TestLatest(t *testing.T) {
	got := Latest()
	require.NotEmpty(t, got)
}

func TestDefault(t *testing.T) {
	got := Default()
	require.NotEmpty(t, got)
}

func TestLatestOrDefault(t *testing.T) {
	got := LatestOrDefault()
	require.NotEmpty(t, got)
}

func TestList(t *testing.T) {
	got := List()
	require.NotEmpty(t, got)
}

func TestGoCmdEnv(t *testing.T) {
	l := LatestOrDefault()
	var env []string
	got := GoCmdEnv(l, env)
	require.NotEmpty(t, got)
	str := strings.Join(got, ";")
	require.Contains(t, str, "GOROOT")
	require.Contains(t, str, "PATH")
}
