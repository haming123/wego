package klog

import "time"

const min_log_buf_size = 4096
const max_log_buf_size = 1024*1024
type LogBuffer struct {
	buf []byte
	n   int
}

func NewLogBufferSize(size int) *LogBuffer {
	if size <= min_log_buf_size {
		size = min_log_buf_size
	}
	return &LogBuffer{
		buf: make([]byte, size),
	}
}

func NewLogBuffer() *LogBuffer {
	return NewLogBufferSize(min_log_buf_size)
}

func (b *LogBuffer) Size() int {
	return len(b.buf)
}

func (b *LogBuffer) Reset() {
	if len(b.buf) > max_log_buf_size {
		b.buf = make([]byte, min_log_buf_size)
	}
	b.n = 0
}

func (b *LogBuffer) GetBytes() []byte {
	return b.buf[0:b.n]
}

func (b *LogBuffer) Available() int {
	return len(b.buf) - b.n
}

func (b *LogBuffer) Buffered() int {
	return b.n
}

func (b *LogBuffer)Extend(size int)  {
	if len(b.buf) > size {
		return
	}
	buf := make([]byte, size)
	copy(buf, b.buf[0:b.n])
	b.buf = buf
}

func (b *LogBuffer)Write(data []byte) {
	d_len := len(data)
	if d_len > b.Available() {
		size := b.n + d_len
		b.Extend(size)
	}
	copy(b.buf[b.n:], data)
	b.n += d_len
}

func (b *LogBuffer)WriteByte(c byte) {
	d_len := 1
	if d_len > b.Available() {
		size := b.n + d_len
		b.Extend(size)
	}
	b.buf[b.n] = c
	b.n++
}

func (b *LogBuffer)WriteString(s string) {
	d_len := len(s)
	if d_len < 0 {
		return
	}
	if d_len > b.Available() {
		size := b.n + d_len
		b.Extend(size)
	}
	copy(b.buf[b.n:], s)
	b.n += d_len
}

func itoa_width(val int, wid int) []byte {
	var b [20]byte
	bp := len(b) - 1
	for val >= 10 || wid > 1 {
		wid--
		q := val / 10
		b[bp] = byte('0' + val - q*10)
		bp--
		val = q
	}
	// i < 10
	b[bp] = byte('0' + val)
	return b[bp:]
}

func (b *LogBuffer)WriteTimeString(t time.Time) {
	//date
	year, month, day := t.Date()
	b.Write(itoa_width(year, 4))
	b.WriteByte('/')
	b.Write(itoa_width(int(month), 2))
	b.WriteByte('/')
	b.Write(itoa_width(day, 2))
	b.WriteByte(' ')
	//time
	hour, min, sec := t.Clock()
	b.Write(itoa_width(hour, 2))
	b.WriteByte(':')
	b.Write(itoa_width(min, 2))
	b.WriteByte(':')
	b.Write(itoa_width(sec, 2))
}


