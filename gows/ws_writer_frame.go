package gows

import (
	"encoding/binary"
	"errors"
	"math"
	"time"
)

const maxFrameHeaderSize = 14 // header(2)+ max-length(8) + mask(4)
const maxControlFrameSize = 125

type FrameWriter struct {
	ws     *WebSocket
	opts   *AcceptOptions
	buff   []byte
	opcode int
	flate  bool
	pos    int
	trim   int
}

func NewFrameWriter(ws *WebSocket, opcode int) *FrameWriter {
	var w FrameWriter
	w.init(ws, opcode)
	return &w
}

func (w *FrameWriter) init(ws *WebSocket, opcode int) {
	w.ws = ws
	w.opts = ws.opts
	w.opcode = opcode
	w.buff = make([]byte, ws.opts.frameWriteBuffSize)
	w.pos = maxFrameHeaderSize
	w.flate = false
	w.trim = 0
}

func (w *FrameWriter) Reset(ws *WebSocket, opcode int) error {
	w.ws = ws
	w.opcode = opcode
	return nil
}

// init后，opcode != (Frame_Continue, Frame_Null)
// 写数据后，opcode == (Frame_Continue, Frame_Null)
// colse后，opcode == Frame_Null
//
// colse前检查是否已经发送了结束帧（opcode == Frame_Null），
// 若已经发送了结束帧，则退出
func (w *FrameWriter) Close() error {
	if w.opcode == Frame_Null {
		return nil
	}

	err := w.writeMessageFrame(true)
	if err != nil {
		return err
	}

	w.ws = nil
	w.opcode = Frame_Null
	return nil
}

func (w *FrameWriter) available() int {
	return len(w.buff) - w.pos
}

func (w *FrameWriter) GetPayloadLength() int {
	return w.pos - maxFrameHeaderSize
}

func (w *FrameWriter) GetPayload() []byte {
	return w.buff[maxFrameHeaderSize:w.pos]
}

func (w *FrameWriter) SetTrimlength(trim int) {
	w.trim = trim
}

func (w *FrameWriter) Write(p []byte) (int, error) {
	//close后就不能调用Write
	if w.opcode == Frame_Null {
		return 0, errors.New("invalid opcode")
	}

	plen := len(p)
	for len(p) > 0 {
		nn := copy(w.buff[w.pos:], p)
		w.pos += nn
		p = p[nn:]

		if w.pos >= len(w.buff) {
			err := w.writeMessageFrame(false)
			if err != nil {
				return 0, err
			}
		}
	}

	return plen, nil
}

func (w *FrameWriter) WriteAll(data []byte) error {
	for len(data) > 0 {
		nn, err := w.Write(data)
		if err != nil {
			return err
		}
		data = data[nn:]
	}
	return nil
}

func (w *FrameWriter) WriteString(str string) error {
	return w.WriteAll(StringToBytes(str))
}

func (w *FrameWriter) writeMessageFrame(final bool) error {
	//close后就不能调用writeMessageFrame
	if w.opcode == Frame_Null {
		return errors.New("invalid opcode")
	}

	//若是压缩格式，则总是预留w.trim_len个字节
	data_end := w.pos
	if w.trim > 0 {
		data_end -= w.trim
		if data_end < maxFrameHeaderSize {
			data_end = maxFrameHeaderSize
		}
	}

	b0 := byte(w.opcode)
	if final {
		b0 |= 1 << 7
	}
	if w.flate {
		b0 |= 1 << 6
	}

	//data_beg为帧的起始位置， 真的payload的开始位置在：maxFrameHeaderSize
	//不同的payload的长度会造成data_beg位置的不同，以下代码用于计算data_beg的位置
	b1 := byte(0)
	data_beg := 4 //server side has no mask(4 byte)
	payload_len := data_end - maxFrameHeaderSize
	switch {
	case payload_len >= math.MaxUint16:
		data_beg += 0
		w.buff[data_beg] = b0
		w.buff[data_beg+1] = b1 | 127
		binary.BigEndian.PutUint64(w.buff[data_beg+2:], uint64(payload_len))
	case payload_len > 125:
		data_beg += 6
		w.buff[data_beg] = b0
		w.buff[data_beg+1] = b1 | 126
		binary.BigEndian.PutUint16(w.buff[data_beg+2:], uint16(payload_len))
	default:
		data_beg += 8
		w.buff[data_beg] = b0
		w.buff[data_beg+1] = b1 | byte(payload_len)
	}

	//将buff中的数据写入：ws_cnn.conn
	w.ws.mux.Lock()
	buff_temp := w.buff[data_beg:data_end]
	writeWait := w.ws.writeTimeOut
	if writeWait > 0 {
		w.ws.cnn.SetWriteDeadline(time.Now().Add(writeWait))
	}
	for len(buff_temp) > 0 {
		nn, err := w.ws.cnn.Write(buff_temp)
		if err != nil {
			w.ws.mux.Unlock()
			return err
		}
		buff_temp = buff_temp[nn:]
	}
	w.ws.mux.Unlock()
	logPrintf4ws(w.ws, "send message frame opcode=%d fin=%v len=%d use_flate=%v\n", w.opcode, final, payload_len, w.flate)

	//若是压缩格式，将没有写入ws_cnn.conn的数据添加到w.buff中
	//若是结束帧，则不用执行此操作
	buff_len := maxFrameHeaderSize
	if final == false && w.pos-data_end > 0 {
		nn := copy(w.buff[buff_len:], w.buff[data_end:])
		buff_len += nn
	}

	//恢复MessageWriter的状态
	// RSV1 bit must be set only for the first frame.
	if final == true {
		w.pos = maxFrameHeaderSize
		w.opcode = Frame_Null
		w.flate = false
	} else {
		w.pos = buff_len
		w.opcode = Frame_Continue
		w.flate = false
	}
	return nil
}

func (w *FrameWriter) WriteControlFrame(data []byte) error {
	if len(data) > maxControlFrameSize {
		return errors.New("websocket: invalid control frame")
	}

	dlen := len(data)
	data_beg := maxFrameHeaderSize - 2
	buff := w.buff[data_beg:]
	buff[0] = byte(w.opcode) | 1<<7
	buff[1] = byte(dlen)
	copy(buff[2:], data)

	w.ws.mux.Lock()
	buff_temp := buff[0 : dlen+2]
	writeWait := w.ws.writeTimeOut
	if writeWait > 0 {
		w.ws.cnn.SetWriteDeadline(time.Now().Add(writeWait))
	}
	for len(buff_temp) > 0 {
		nn, err := w.ws.cnn.Write(buff_temp)
		if err != nil {
			w.ws.mux.Unlock()
			return err
		}
		buff_temp = buff_temp[nn:]
	}
	w.ws.mux.Unlock()
	logPrintf4ws(w.ws, "send control frame opcode=%d len=%d\n", w.opcode, dlen)

	w.pos = maxFrameHeaderSize
	w.opcode = Frame_Null
	w.flate = false
	return nil
}
