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

// Zip 解压缩 .zip 文件
type Zip struct {
	// UnpackNextBefore 在 Unpack 时，解析到下一个 Header 后，实际 unpack 前的回调
	UnpackNextBefore func(f *zip.File) (skip bool, err error)

	// UnpackNextAfter 在 Unpack 时，解析到下一个 Header 后，实际 unpack 后的回调
	UnpackNextAfter func(f *zip.File, err error) error

	// StripComponents Unpack 的时候，忽略掉前 N 层目录
	StripComponents uint

	// MinSize 最小文件大小，>0 时有效
	MinSize int64

	// MaxSize 最小文件大小，>0 时有效
	MaxSize int64

	// IgnoreFailed 是否忽略异常
	// 不会忽略 UnpackNextBefore 返回的 error
	IgnoreFailed bool
}

func (zp *Zip) unpackTo(p string) string {
	if zp.StripComponents == 0 {
		return p
	}
	sc := int(zp.StripComponents)
	ps := strings.Split(filepath.Clean(p), string(filepath.Separator))
	if len(ps) < sc {
		return ""
	}

	return filepath.Join(ps[sc:]...)
}

// Unpack 解压缩文件到指定目录
func (zp *Zip) Unpack(archiveFile string, targetDir string) error {
	zr, err := zip.OpenReader(archiveFile)
	if err != nil {
		return err
	}
	defer zr.Close()
	return zp.UnpackFromReader(&zr.Reader, targetDir)
}

// UnpackFromReader 解压 zip.Reader
func (zp *Zip) UnpackFromReader(zrd *zip.Reader, targetDir string) error {
	for _, f := range zrd.File {
		if !zp.checkMinMaxIgnore(f) {
			continue
		}
		if zp.UnpackNextBefore != nil {
			if skip, err4 := zp.UnpackNextBefore(f); skip {
				continue
			} else if err4 != nil {
				return err4
			}
		}

		err3 := zp.unpackOne(f, targetDir)

		if zp.UnpackNextAfter != nil {
			err3 = zp.UnpackNextAfter(f, err3)
		}

		if err3 != nil && !zp.IgnoreFailed {
			return err3
		}
	}
	return nil
}

func (zp *Zip) checkMinMaxIgnore(f *zip.File) bool {
	if !f.FileInfo().Mode().IsRegular() {
		return false
	}
	size := f.FileInfo().Size()
	if zp.MinSize > 0 && size < zp.MinSize {
		return true
	}
	if zp.MaxSize > 0 && size > zp.MaxSize {
		return true
	}
	return false
}

func (zp *Zip) unpackOne(f *zip.File, targetDir string) error {
	to := zp.unpackTo(f.Name)
	// 若是文件名为空，则此文件忽略掉
	if len(to) == 0 {
		return nil
	}

	outPath := filepath.Join(targetDir, to)

	if f.FileInfo().IsDir() {
		return mkdir(outPath)
	}

	rc, err2 := f.Open()
	if err2 != nil {
		return err2
	}
	defer rc.Close()

	if err3 := mkdir(filepath.Dir(outPath)); err3 != nil {
		return err3
	}

	out, err4 := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err4 != nil {
		return err4
	}
	defer out.Close()
	return copyFile(rc, out, -1)
}
