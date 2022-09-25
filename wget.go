// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/16

package cmdutils

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Wget struct {
	LogWriter io.Writer

	Proxy func(*http.Request) (*url.URL, error)

	// Timeout 整体超时
	Timeout time.Duration

	// ConnectTimeout 连接超时
	ConnectTimeout time.Duration

	InsecureSkipVerify bool
}

func (w *Wget) getProxy() func(*http.Request) (*url.URL, error) {
	if w.Proxy != nil {
		return w.Proxy
	}
	return http.ProxyFromEnvironment
}

func (w *Wget) logit(msgs ...interface{}) {
	if w.LogWriter == nil {
		return
	}
	var b strings.Builder
	for _, m := range msgs {
		b.WriteString(fmt.Sprint(m))
		b.WriteString(" ")
	}
	b.WriteString("\n")
	fmt.Fprint(w.LogWriter, b.String())
}

func (w *Wget) getClient() *http.Client {
	tr := &http.Transport{
		DisableCompression: true,
		DisableKeepAlives:  true,
		Proxy:              w.getProxy(),
		DialContext:        w.dialContext,
	}
	if w.InsecureSkipVerify {
		tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true, // 不校验 server 的证书的有效性
		}
	}
	return &http.Client{
		Transport: tr,
		Timeout:   w.Timeout,
	}
}

func (w *Wget) dialContext(ctx context.Context, network, addr string) (c net.Conn, err error) {
	if w.ConnectTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, w.ConnectTimeout)
		defer cancel()
	}
	start := time.Now()
	w.logit("connect start", network, addr)
	defer func() {
		cost := time.Since(start)
		if err == nil {
			w.logit("connect success", network, addr, "cost=", cost.String())
		} else {
			w.logit("connect failed", network, addr, "cost=", cost.String(), err)
		}
	}()
	return (&net.Dialer{}).DialContext(ctx, network, addr)
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

	bw := bufio.NewWriter(dstFile)
	defer bw.Flush()

	if err = w.DownloadToWriter(src, bw); err != nil {
		return err
	}

	return dstFile.Close()
}

func (w *Wget) DownloadToWriter(src string, dst io.Writer) error {
	client := w.getClient()

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
	if w.LogWriter != nil {
		pw = &progressWriter{
			w:     dst,
			po:    w.LogWriter,
			total: res.ContentLength,
		}
		ww = pw
	} else {
		ww = dst
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
	return nil
}

type progressWriter struct {
	last     time.Time
	w        io.Writer
	po       io.Writer
	n        int64
	total    int64
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
