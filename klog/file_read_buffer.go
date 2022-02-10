package klog

import (
	"io"
)

const minReadBufferSize = 81920

type Reader struct {
	buf		[]byte
	r		int
	w		int
}

func (b *Reader) Reset() {
	b.r = 0
	b.w = 0
}

func NewReaderSize(size int) *Reader {
	if size < minReadBufferSize {
		size = minReadBufferSize
	}

	rb := new(Reader)
	rb.buf = make([]byte, size)
	rb.r = 0
	rb.w = 0

	return rb
}

func NewReader() *Reader {
	return NewReaderSize(minReadBufferSize)
}

func (b *Reader) Buffered() int {
	return b.w - b.r
}

func (b *Reader) ReadFromFile(rd io.Reader) {
	if b.r > 0 {
		copy(b.buf, b.buf[b.r:b.w])
		b.w -= b.r
		b.r = 0
	}

	if b.w >= len(b.buf) {
		return
	}

	n, _ := rd.Read(b.buf[b.w:])
	b.w += n
}

func (b *Reader) Read(p []byte) (int, error) {
	p = b.buf[b.r:b.w]
	return b.Buffered(), nil
}

func (b *Reader) GetBytes() []byte  {
	return b.buf[b.r:b.w]
}
