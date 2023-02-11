package gows

import (
	"strings"
	"sync"
)

const minCompressSize = 1024
const frameWriteBuffSize = 4096
const frameReadBuffSize = 4096
const messageBufferSize = 4096

type AcceptOptions struct {
	SubProtocols []string

	compress_alloter CompressAlloter
	minCompressSize  int

	frameWriteBuffSize int
	frameReadBuffSize  int
	messageBufferSize  int

	messageWriterPool sync.Pool
	messageReaderPool sync.Pool
	messageBufferPool sync.Pool
}

func NewAcceptOptions() *AcceptOptions {
	var opts AcceptOptions
	opts.init()
	return &opts
}

func (opts *AcceptOptions) init() {
	opts.compress_alloter = nil
	opts.minCompressSize = minCompressSize
	opts.frameWriteBuffSize = frameWriteBuffSize
	opts.frameReadBuffSize = frameReadBuffSize
	opts.messageBufferSize = messageBufferSize
}

func (this *AcceptOptions) selectSubProtocol(cps []string) string {
	for _, sp := range this.SubProtocols {
		for _, cp := range cps {
			if strings.EqualFold(sp, cp) {
				return cp
			}
		}
	}
	return ""
}

func (this *AcceptOptions) UseFlate(val ...CompressAlloter) {
	this.compress_alloter = &flate_default
	if len(val) == 1 {
		this.compress_alloter = val[0]
	}
}

func (this *AcceptOptions) SetMinCompressSize(size int) {
	this.minCompressSize = size
}

func (this *AcceptOptions) SetFrameReadBuffSize(size int) {
	if size >= 512 {
		this.frameReadBuffSize = size
	}
}

func (this *AcceptOptions) SetMessageBufferSize(size int) {
	//if size >= messageBufferSize {
	this.messageBufferSize = size
	//}
}
