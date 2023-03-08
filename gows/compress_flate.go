package gows

import (
	"compress/flate"
	"io"
)

// https://www.rfc-editor.org/rfc/rfc7692
//
// 这些字节："\x00\x00\xff\xff"在发送时被要求删除。
// 因此在进行消息读取时需要将它们添加回，
// 否则flate.Reader会继续尝试读取更多字节。
//
// 添加结束标志："\x01\x00\x00\xff\xff"来防止flate.reader产生：unexpected EOF错误
const deflateMessageTail = "\x00\x00\xff\xff\x01\x00\x00\xff\xff"

// Four extension parameters are defined for "permessage-deflate" to help endpoints manage per-connection resource usage.
//
//	"server_no_context_takeover"
//	"client_no_context_takeover"
//	"server_max_window_bits"
//	"client_max_window_bits"
//
// "use context takeover" :
//
//	The term "use context takeover" means that the same LZ77 sliding window
//	used by the endpoint to build frames of the previous sent message
//	is reused to build frames of the next message to be sent.
//
// "server_no_context_takeover" ：
//
//	Extension Parameter：If the peer server doesn't use context takeover,
//	the client doesn't need to reserve memory to retain the LZ77 sliding window between messages.
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

func (this *FlateAlloter) WebsocketExtension(params []string) string {
	return "permessage-deflate; server_no_context_takeover; client_no_context_takeover"
}

func (this *FlateAlloter) NewWriter(mw *FrameWriter) (io.WriteCloser, error) {
	//logPrint("new use_flate.Writer !!!")
	mw.flate = true
	//用于删除这四个字节"\x00\x00\xff\xff"
	mw.SetTrimlength(4)
	return flate.NewWriter(mw, this.level)
}

func (this *FlateAlloter) ResetWriter(fw io.WriteCloser, mw *FrameWriter) error {
	//logPrint("reset use_flate.Writer...")
	mw.flate = true
	//用于删除这四个字节"\x00\x00\xff\xff"
	mw.SetTrimlength(4)
	fw.(*flate.Writer).Reset(mw)
	return nil
}

func (this *FlateAlloter) FlushWriter(fw io.WriteCloser) error {
	return fw.(*flate.Writer).Flush()
}

func (this *FlateAlloter) NewReader(mr *FrameReader) (io.ReadCloser, error) {
	//logPrint("new use_flate.Reader !!!")
	//添加结束标志，防止flate.reader产生：unexpected EOF错误
	mr.extra.Reset(deflateMessageTail)
	return flate.NewReader(mr), nil
}

func (this *FlateAlloter) ResetReader(fr io.ReadCloser, mr *FrameReader) error {
	//logPrint("reset use_flate.Reade...")
	//添加结束标志，防止flate.reader产生：unexpected EOF错误
	mr.extra.Reset(deflateMessageTail)
	return fr.(flate.Resetter).Reset(mr, nil)
}
