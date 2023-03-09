package gows

import (
	"errors"
	"io"
)

var frame_reader_idle FrameReader

type MessageReader struct {
	opts *AcceptOptions
	//若非压缩格式，调用者直接从frame中读取数据
	frame *FrameReader
	//若是压缩格式，调用者从flate读取数据
	flate io.ReadCloser
}

// 创建MessageReader时，若当前环境支持压缩，则创建压缩解码器
//
// 若是压缩数据，首先调用MessageReader.read()
// MessageReader.read() 调用mr.flate_flag.read()
// mr.flate_flag.read()调用mr.frame来读取数据
//
// 若是非压缩数据，则首先调用MessageReader.read()
// MessageReader.read() 直接调用mr.frame来读取数据
func newMessageReader(ws *WebSocket) *MessageReader {
	mr := &MessageReader{}
	mr.opts = ws.opts
	mr.frame = ws.msgReader
	if mr.opts.compressAlloter != nil {
		mr.flate, _ = mr.opts.compressAlloter.NewReader(ws.msgReader)
	}
	return mr
}

// 恢复压缩解码器的状态
func (mr *MessageReader) reset(ws *WebSocket) {
	mr.opts = ws.opts
	mr.frame = ws.msgReader
	if mr.flate != nil && mr.opts.compressAlloter != nil {
		mr.opts.compressAlloter.ResetReader(mr.flate, ws.msgReader)
	}
}

// 关闭时将压缩读取器指向一个空的FrameReader
func (mr *MessageReader) close() error {
	mr.frame = &frame_reader_idle
	if mr.flate != nil && mr.opts.compressAlloter != nil {
		mr.opts.compressAlloter.ResetReader(mr.flate, &frame_reader_idle)
	}
	return nil
}

func (mr *MessageReader) Close() error {
	return CloseReader(mr)
}

func (mr *MessageReader) readMessageHeader() (FrameHeader, error) {
	return mr.frame.readMessageHeader()
}

// 判断消息头，看看当前消息是否采用了压缩格式，若是压缩格式，则使用flate来读取数据
// 否则使用mr.frame来读取数据
func (mr *MessageReader) getMatchedReader() io.Reader {
	var reader io.Reader
	if mr.frame.header.flate == true {
		reader = mr.flate
	} else {
		mr.frame.extra.Reset("")
		reader = mr.frame
	}
	return reader
}

func (mr *MessageReader) Read(p []byte) (int, error) {
	reader := mr.getMatchedReader()
	if reader == nil {
		return 0, errors.New("can not allocate flate_flag reader")
	}
	return reader.Read(p)
}

func (mr *MessageReader) ReadAll() (*ByteBuffer, error) {
	reader := mr.getMatchedReader()
	if reader == nil {
		return nil, errors.New("can not allocate flate_flag reader")
	}
	mb := GetByteBuffer(mr.opts)
	err := mb.ReadAll(reader)
	return mb, err
}
