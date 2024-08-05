// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/8/5

//go:build plan9 || windows

package gosdk

import (
	"os"
)

var signalsToIgnore = []os.Signal{os.Interrupt}
