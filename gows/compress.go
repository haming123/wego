package gows

import "io"

type CompressMode int

const (
	CompressDisabled CompressMode = iota
	CompressContextTakeover
	CompressNoContextTakeover
)

type CompressAlloter interface {
	GetReponseExtensions(args []string) string
	NewWriter(mw *FrameWriter) (io.WriteCloser, error)
	ResetWriter(fw io.WriteCloser, mw *FrameWriter) error
	FlushWriter(fw io.WriteCloser) error
	NewReader(mr *FrameReader) (io.ReadCloser, error)
	ResetReader(fr io.ReadCloser, mr *FrameReader) error
}
