package gows

import (
	"io"
	"unsafe"
)

func GetByteBuffer(opts *AcceptOptions) *ByteBuffer {
	b, ok := opts.messageBufferPool.Get().(*ByteBuffer)
	if !ok {
		//logPrint("new ByteBuffer !!!")
		return newByteBuffer(opts)
	}
	b.Reset()
	//logPrint("get ByteBuffer from poll")
	return b
}

func PutByteBuffer(b *ByteBuffer) {
	b.opts.messageBufferPool.Put(b)
	//logPrint("free ByteBuffer")
}

type ByteBuffer struct {
	opts *AcceptOptions
	buf  []byte
	pos  int
}

func newByteBuffer(opts *AcceptOptions) *ByteBuffer {
	return &ByteBuffer{
		opts: opts,
		buf:  make([]byte, opts.messageBufferSize),
		pos:  0,
	}
}

func (b *ByteBuffer) Close() {
	size := len(b.buf)
	if b.opts != nil && size == b.opts.messageBufferSize {
		PutByteBuffer(b)
	}
}

func (b *ByteBuffer) Reset() {
	b.pos = 0
}

func (b *ByteBuffer) Size() int {
	return b.pos
}

func (b *ByteBuffer) ReadFull(reader io.Reader) error {
	size := len(b.buf)
	for b.pos < size {
		n, err := reader.Read(b.buf[b.pos:])
		b.pos += n
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *ByteBuffer) ReadAll(reader io.Reader) error {
	for {
		size := len(b.buf)
		if b.pos >= size {
			btmp := b.buf
			b.buf = make([]byte, size*2)
			copy(b.buf, btmp)
		}

		n, err := reader.Read(b.buf[b.pos:])
		b.pos += n
		if err != nil {
			return err
		}
	}
}

func (b *ByteBuffer) GetBytes() []byte {
	return b.buf[0:b.pos]
}

func (b *ByteBuffer) CloneBytes() []byte {
	data := make([]byte, b.pos)
	copy(data, b.buf[0:b.pos])
	return data
}

func (b *ByteBuffer) GetString() string {
	if b.buf == nil {
		return ""
	}
	buff := b.buf[0:b.pos]
	return string(buff)
}

func (b *ByteBuffer) GetVolatileString() string {
	if b.buf == nil {
		return ""
	}
	buff := b.buf[0:b.pos]
	return *(*string)(unsafe.Pointer(&buff))
}
