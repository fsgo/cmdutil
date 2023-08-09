// Copyright(C) 2023 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2023/8/9

package cmdutil

import "regexp"

var colorReg = regexp.MustCompile(`\x1b\[\d+m`)

func CleanColor(str string) string {
	return colorReg.ReplaceAllString(str, "")
}
