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
	got := DefaultOrLatest()
	fst.NotEmpty(t, got)
}

func TestLatest(t *testing.T) {
	got := Latest()
	fst.NotEmpty(t, got)
}

func TestDefault(t *testing.T) {
	got := Default()
	fst.NotEmpty(t, got)
}

func TestLatestOrDefault(t *testing.T) {
	got := LatestOrDefault()
	fst.NotEmpty(t, got)
}

func TestList(t *testing.T) {
	got := List()
	fst.NotEmpty(t, got)
}

func TestGoCmdEnv(t *testing.T) {
	l := LatestOrDefault()
	var env []string
	got := GoCmdEnv(l, env)
	fst.NotEmpty(t, got)
	str := strings.Join(got, ";")
	fst.StringContains(t, str, "GOROOT")
	fst.StringContains(t, str, "PATH")
}
