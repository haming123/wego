package gows

import (
	"compress/flate"
	"io"
)

// 这些字节："\x00\x00\xff\xff"在发送时被要求删除。
// 因此在进行消息读取时需要将它们添加回，
// 否则flate.Reader会继续尝试读取更多字节。
//
// 添加结束标志："\x01\x00\x00\xff\xff"来防止flate.reader产生：unexpected EOF错误
const deflateMessageTail = "\x00\x00\xff\xff\x01\x00\x00\xff\xff"

var flate_default FlateAlloter

func init() {
	flate_default.level = flate.BestSpeed
}

type FlateAlloter struct {
	level int
}

func NewFlateAlloter(level int) *FlateAlloter {
	ent := &FlateAlloter{}
	ent.level = level
	return ent
}

func (this *FlateAlloter) GetReponseExtensions(params []string) string {
	return "permessage-deflate; server_no_context_takeover; client_no_context_takeover"
}

func (this *FlateAlloter) NewWriter(mw *FrameWriter) (io.WriteCloser, error) {
	//logPrint("new flate.Writer !!!")
	mw.flate = true
	//用于删除这四个字节"\x00\x00\xff\xff"
	mw.SetTrimlength(4)
	return flate.NewWriter(mw, this.level)
}

func (this *FlateAlloter) ResetWriter(fw io.WriteCloser, mw *FrameWriter) error {
	//logPrint("reset flate.Writer...")
	mw.flate = true
	//用于删除这四个字节"\x00\x00\xff\xff"
	mw.SetTrimlength(4)
	fw.(*flate.Writer).Reset(mw)
	return nil
}

func (this *FlateAlloter) NewReader(mr *FrameReader) (io.ReadCloser, error) {
	//logPrint("new flate.Reader !!!")
	//添加结束标志，防止flate.reader产生：unexpected EOF错误
	mr.extra.Reset(deflateMessageTail)
	return flate.NewReader(mr), nil
}

func (this *FlateAlloter) ResetReader(fr io.ReadCloser, mr *FrameReader) error {
	//logPrint("reset flate.Reade...")
	//添加结束标志，防止flate.reader产生：unexpected EOF错误
	mr.extra.Reset(deflateMessageTail)
	return fr.(flate.Resetter).Reset(mr, nil)
}
