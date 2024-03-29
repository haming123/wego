package gows

import "io"

type CompressMode int

const (
	CompressDisabled CompressMode = iota
	CompressContextTakeover
	CompressNoContextTakeover
)

type CompressAlloter interface {
	WebsocketExtension(args []string) string
	NewWriter(mw *FrameWriter) (io.WriteCloser, error)
	FlushWriter(fw io.WriteCloser) error
	ResetWriter(fw io.WriteCloser, mw *FrameWriter) error
	NewReader(mr *FrameReader) (io.ReadCloser, error)
	ResetReader(fr io.ReadCloser, mr *FrameReader) error
}
