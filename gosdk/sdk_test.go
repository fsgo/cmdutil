// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/7

package gosdk

import (
	"strings"
	"testing"

	"github.com/fsgo/fst"
)

func TestDefaultOrLatest(t *testing.T) {
	got := DefaultOrLatest(t.Context())
	fst.NotEmpty(t, got)
}

func TestLatest(t *testing.T) {
	got := Latest(t.Context())
	fst.NotEmpty(t, got)
}

func TestDefault(t *testing.T) {
	got := Default(t.Context())
	fst.NotEmpty(t, got)
}

func TestLatestOrDefault(t *testing.T) {
	got := LatestOrDefault(t.Context())
	fst.NotEmpty(t, got)
}

func TestList(t *testing.T) {
	got := List(t.Context())
	fst.NotEmpty(t, got)
}

func TestGoCmdEnv(t *testing.T) {
	l := LatestOrDefault(t.Context())
	var env []string
	got := GoCmdEnv(l, env)
	fst.NotEmpty(t, got)
	str := strings.Join(got, ";")
	fst.Contains(t, str, "GOROOT")
	fst.Contains(t, str, "PATH")
}
