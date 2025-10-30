// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/7

package gosdk

import (
	"strings"
	"testing"

	"github.com/xanygo/anygo/xt"
)

func TestDefaultOrLatest(t *testing.T) {
	got := DefaultOrLatest(t.Context())
	xt.NotEmpty(t, got)
}

func TestLatest(t *testing.T) {
	got := Latest(t.Context())
	xt.NotEmpty(t, got)
}

func TestDefault(t *testing.T) {
	got := Default(t.Context())
	xt.NotEmpty(t, got)
}

func TestLatestOrDefault(t *testing.T) {
	got := LatestOrDefault(t.Context())
	xt.NotEmpty(t, got)
}

func TestList(t *testing.T) {
	got := List(t.Context())
	xt.NotEmpty(t, got)
}

func TestGoCmdEnv(t *testing.T) {
	l := LatestOrDefault(t.Context())
	var env []string
	got := GoCmdEnv(l, env)
	xt.NotEmpty(t, got)
	str := strings.Join(got, ";")
	xt.Contains(t, str, "GOROOT")
	xt.Contains(t, str, "PATH")
}
