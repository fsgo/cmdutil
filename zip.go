// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/16

package cmdutils

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
)

type Zip struct {
	// StripComponents Unpack 的时候，忽略掉前 N 层目录
	StripComponents uint

	// IgnoreFailed 是否忽略异常
	// 不会忽略 UnpackNextBefore 返回的 error
	IgnoreFailed bool

	// MinSize 最小文件大小，>0 时有效
	MinSize int64

	// MaxSize 最小文件大小，>0 时有效
	MaxSize int64

	// UnpackNextBefore 在 Unpack 时，解析到下一个 Header 后，实际 unpack 前的回调
	UnpackNextBefore func(f *zip.File) (skip bool, err error)

	// UnpackNextAfter 在 Unpack 时，解析到下一个 Header 后，实际 unpack 后的回调
	UnpackNextAfter func(f *zip.File, err error) error
}

func (tr *Zip) unpackTo(p string) string {
	if tr.StripComponents == 0 {
		return p
	}
	sc := int(tr.StripComponents)
	ps := strings.Split(filepath.Clean(p), string(filepath.Separator))
	if len(ps) < sc {
		return ""
	}

	return filepath.Join(ps[sc:]...)
}

// Unpack 解压缩文件到指定目录
func (tr *Zip) Unpack(archiveFile string, targetDir string) error {
	zr, err := zip.OpenReader(archiveFile)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, f := range zr.File {
		if tr.UnpackNextBefore != nil {
			if skip, err4 := tr.UnpackNextBefore(f); skip {
				continue
			} else if err4 != nil {
				return err4
			}
		}

		err3 := tr.unpackOne(f, targetDir)

		if tr.UnpackNextAfter != nil {
			err3 = tr.UnpackNextAfter(f, err3)
		}

		if err3 != nil && !tr.IgnoreFailed {
			return err3
		}
	}
	return nil
}

func (tr *Zip) checkMinMaxIgnore(f *zip.File) bool {
	if !f.FileInfo().Mode().IsRegular() {
		return false
	}
	size := f.FileInfo().Size()
	if tr.MinSize > 0 && size < tr.MinSize {
		return true
	}
	if tr.MaxSize > 0 && size > tr.MaxSize {
		return true
	}
	return false
}

func (tr *Zip) mkdir(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil && os.IsExist(err) {
		return nil
	}
	return err
}

func (tr *Zip) unpackOne(f *zip.File, targetDir string) error {
	to := tr.unpackTo(f.Name)
	// 若是文件名为空，则此文件忽略掉
	if len(to) == 0 {
		return nil
	}

	outPath := filepath.Join(targetDir, to)

	if f.FileInfo().IsDir() {
		return tr.mkdir(outPath)
	}

	rc, err2 := f.Open()
	if err2 != nil {
		return err2
	}
	defer rc.Close()

	if err3 := tr.mkdir(filepath.Dir(outPath)); err3 != nil {
		return err3
	}

	out, err4 := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err4 != nil {
		return err4
	}
	defer out.Close()
	return copyFile(rc, out, -1)
}
