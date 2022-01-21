// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/16

package cmdutils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Wget struct {
	PrintProgress io.Writer
	Timeout       time.Duration
	Proxy         func(*http.Request) (*url.URL, error)
}

func (w *Wget) getProxy() func(*http.Request) (*url.URL, error) {
	if w.Proxy != nil {
		return w.Proxy
	}
	return http.ProxyFromEnvironment
}

func (w *Wget) Download(src string, dst string) error {
	if len(dst) == 0 {
		return errors.New("empty output path")
	}

	if err := mkdir(filepath.Dir(dst)); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			dstFile.Close()
			os.Remove(dst)
		}
	}()

	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
			DisableKeepAlives:  true,
			Proxy:              w.getProxy(),
		},
		Timeout: w.Timeout,
	}

	res, err := client.Get(src)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %s", res.Status)
	}

	var pw *progressWriter
	var ww io.Writer
	if w.PrintProgress != nil {
		pw = &progressWriter{
			w:     dstFile,
			po:    w.PrintProgress,
			total: res.ContentLength,
		}
		ww = pw
	} else {
		ww = dstFile
	}
	n, err := io.Copy(ww, res.Body)
	if err != nil {
		return err
	}
	if res.ContentLength != -1 && res.ContentLength != n {
		return fmt.Errorf("copied %v bytes; expected %v", n, res.ContentLength)
	}

	if pw != nil {
		pw.finish()
	}
	return dstFile.Close()
}

type progressWriter struct {
	w        io.Writer
	n        int64
	total    int64
	last     time.Time
	po       io.Writer
	finished bool
}

func (p *progressWriter) finish() {
	p.finished = true
	p.update()
}

func (p *progressWriter) update() {
	end := " ..."
	if p.finished {
		end = ""
	}
	if p.total > 0 {
		fmt.Fprintf(p.po, "Downloaded %5.1f%% (%*d / %d bytes)%s\n",
			(100.0*float64(p.n))/float64(p.total),
			p.ndigits(p.total), p.n, p.total, end)
		return
	}
	fmt.Fprintf(p.po, "Downloaded %d bytes %s\n", p.n, end)
}

func (p *progressWriter) ndigits(i int64) int {
	var n int
	for ; i != 0; i /= 10 {
		n++
	}
	return n
}

func (p *progressWriter) Write(buf []byte) (n int, err error) {
	n, err = p.w.Write(buf)
	p.n += int64(n)
	if now := time.Now(); now.Unix() != p.last.Unix() {
		p.update()
		p.last = now
	}
	return
}
