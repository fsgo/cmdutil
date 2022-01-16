// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/16

package cmdutils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Tar tape archive 工具，目前已经具备压缩文件
type Tar struct {
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
	UnpackNextBefore func(h *tar.Header) (skip bool, err error)

	// UnpackNextAfter 在 Unpack 时，解析到下一个 Header 后，实际 unpack 后的回调
	UnpackNextAfter func(h *tar.Header, err error) error

	// UnCompress Unpack 时的解压缩方法，可选
	// 默认为按照文件后缀自动选择：
	// 1.后缀为 .gz 和 .tgz 时选择 gzip
	UnCompress func(rd io.Reader) (io.Reader, error)
}

func (tr *Tar) validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}

func (tr *Tar) unpackTo(p string) (string, error) {
	if !tr.validRelPath(p) {
		return "", fmt.Errorf("tar file contained invalid name %q", p)
	}
	if tr.StripComponents == 0 {
		return p, nil
	}
	sc := int(tr.StripComponents)
	ps := strings.Split(filepath.Clean(p), string(filepath.Separator))
	if len(ps) <= sc {
		return "", fmt.Errorf("invalid StripComponents=%d for filepath %q", sc, p)
	}

	return filepath.Join(ps[sc:]...), nil
}

func (tr *Tar) unCompress(rd io.Reader, name string) (io.Reader, error) {
	if tr.UnCompress != nil {
		return tr.UnCompress(rd)
	}

	if strings.HasSuffix(name, ".gz") || strings.HasSuffix(name, ".tgz") {
		return gzip.NewReader(rd)
	}

	return rd, nil
}

// Unpack 解压缩文件到指定目录
func (tr *Tar) Unpack(archiveFile string, targetDir string) error {
	tf, err0 := os.Open(archiveFile)
	if err0 != nil {
		return err0
	}
	defer tf.Close()

	zr, err1 := tr.unCompress(tf, filepath.Base(archiveFile))
	if err1 != nil {
		return err1
	}

	if rc, ok := zr.(io.Closer); ok {
		defer rc.Close()
	}

	madeDir := map[string]bool{}

	trd := tar.NewReader(zr)
	for {
		th, err2 := trd.Next()
		if err2 == io.EOF {
			break
		}
		if err2 != nil {
			return err2
		}

		if tr.UnpackNextBefore != nil {
			if skip, err4 := tr.UnpackNextBefore(th); skip {
				return nil
			} else if err4 != nil {
				return err4
			}
		}

		err3 := tr.unpackOne(trd, th, targetDir, madeDir)

		if tr.UnpackNextAfter != nil {
			err3 = tr.UnpackNextAfter(th, err3)
		}

		if err3 != nil && !tr.IgnoreFailed {
			return err3
		}
	}
	return nil
}

func (tr *Tar) checkMinMaxIgnore(th *tar.Header) bool {
	if !th.FileInfo().Mode().IsRegular() {
		return false
	}
	if tr.MinSize > 0 && th.Size < tr.MinSize {
		return true
	}
	if tr.MaxSize > 0 && th.Size > tr.MaxSize {
		return true
	}
	return false
}

func (tr *Tar) unpackOne(trd *tar.Reader, th *tar.Header, targetDir string, madeDir map[string]bool) (ret error) {
	to, err3 := tr.unpackTo(th.Name)
	if err3 != nil {
		return err3
	}

	if tr.checkMinMaxIgnore(th) {
		return nil
	}

	abs := filepath.Join(targetDir, to)
	fi := th.FileInfo()

	mode := fi.Mode()
	switch {
	case mode.IsRegular():
		// Make the directory. This is redundant because it should
		// already be made by a directory entry in the tar
		// beforehand. Thus, don't check for errors; the next
		// write will fail with the same error.
		dir := filepath.Dir(abs)
		if !madeDir[dir] {
			if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
				return err
			}
			madeDir[dir] = true
		}
		wf, err := os.OpenFile(abs, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode.Perm())
		if err != nil {
			return err
		}
		n, err := io.Copy(wf, trd)
		if closeErr := wf.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
		if err != nil {
			return fmt.Errorf("error writing to %s: %v", abs, err)
		}
		if n != th.Size {
			return fmt.Errorf("only wrote %d bytes to %s; expected %d", n, abs, th.Size)
		}
		if !th.ModTime.IsZero() {
			if err := os.Chtimes(abs, th.ModTime, th.ModTime); err != nil {
				return err
			}
		}
	case mode.IsDir():
		if err := os.MkdirAll(abs, 0755); err != nil {
			return err
		}
		madeDir[abs] = true
	default:
		return fmt.Errorf("tar file entry %s contained unsupported file type %v", th.Name, mode)
	}
	return nil
}
