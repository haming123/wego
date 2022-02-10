package filters

import (
	"errors"
	"io"
	"net/http"
	"github.com/haming123/wego"
)

type maxBytesReader struct {
	c  	 	*wego.WebContext
	r   	io.ReadCloser
	n   	int64
	err 	error
	write 	bool
}

func (l *maxBytesReader) Read(p []byte) (n int, err error) {
	if l.err != nil {
		return 0, l.err
	}
	if len(p) == 0 {
		return 0, nil
	}
	if int64(len(p)) > l.n+1 {
		p = p[:l.n+1]
	}
	n, err = l.r.Read(p)

	if int64(n) <= l.n {
		l.n -= int64(n)
		l.err = err
		return n, err
	}

	n = int(l.n)
	l.n = 0

	l.err = errors.New("http: request body too large")
	if l.write == false {
		l.write = true
		l.c.SetHeader("connection", "close")
		l.c.AbortWithError(http.StatusRequestEntityTooLarge, l.err)
	}
	return n, l.err
}

func (l *maxBytesReader) Close() error {
	return l.r.Close()
}

func RequestSizeLimiter(limit int64) wego.HandlerFunc {
	return func(c *wego.WebContext) {
		c.Input.Body = &maxBytesReader{c: c, r: c.Input.Body, n: limit}
		c.Next()
	}
}
