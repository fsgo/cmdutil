// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/8/5

package main

import (
	"fmt"

	"github.com/fsgo/cmdutil/gosdk"
)

func main() {
	s := &gosdk.SDK{}
	fmt.Println("Default()=", s.Default())
	fmt.Println("Latest()=", s.Latest())
	fmt.Println("DefaultOrLatest()=", s.DefaultOrLatest())
	fmt.Println("Find(1.21)=", s.Find("1.21"))
}
