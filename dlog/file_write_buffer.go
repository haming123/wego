package log

import (
	"io"
	"os"
)

const def_file_buff_size = 8192
type BufferWriter struct {
	name string
	err error
	buf []byte
	n   int
	file  *os.File
}

func NewBufferWriter(name string, flag int, perm os.FileMode) (*BufferWriter, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	bf := &BufferWriter {
		buf: 	make([]byte, def_file_buff_size),
		file:  	file,
		name:	name,
	}
	return bf, nil
}

func (b *BufferWriter) Name() string {
	return b.name
}

func (b *BufferWriter) Size() int {
	return len(b.buf)
}

func (b *BufferWriter) Reset() {
	b.err = nil
	b.n = 0
}

func (b *BufferWriter) Flush() error {
	if b.err != nil {
		return b.err
	}
	if b.n == 0 {
		return nil
	}
	n, err := b.file.Write(b.buf[0:b.n])
	if n < b.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.n {
			copy(b.buf[0:b.n-n], b.buf[n:b.n])
		}
		b.n -= n
		b.err = err
		return err
	}
	b.n = 0
	return nil
}

func (b *BufferWriter)Close() error {
	err := b.Flush()
	if err != nil {
		return err
	}
	return b.file.Close()
}

func (b *BufferWriter) Available() int {
	return len(b.buf) - b.n
}

func (b *BufferWriter) Buffered() int {
	return b.n
}

func (b *BufferWriter) Write(data []byte) (nn int, err error) {
	d_len := len(data)

	if d_len > b.Available() {
		b.Flush()
	}

	if d_len > b.Available() {
		return b.file.Write(data)
	}

	copy(b.buf[b.n:], data)
	b.n += d_len
	return d_len, nil
}
