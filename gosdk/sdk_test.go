// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/7

package gosdk

import (
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
