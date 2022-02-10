package log

import (
	"io"
	"os"
	"time"
)

type TermWriter struct {
	out   io.Writer
}

func NewTermWriter() *TermWriter {
	consoleWriter := &TermWriter{
		out:  os.Stdout,
	}
	return consoleWriter
}

func (cw *TermWriter) Write(tm time.Time, data []byte) {
	cw.out.Write([]byte(data))
}

func (cw *TermWriter) Flush() error {
	return nil
}

func (cw *TermWriter) Close() error {
	return nil
}



