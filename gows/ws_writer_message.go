package gows

import (
	"io"
)

type MessageWriter struct {
	ws     *WebSocket
	opts   *AcceptOptions
	opcode int

	//非压缩的数据写到mframe中
	mframe FrameWriter
	//压缩结果写到fframe中
	fframe FrameWriter
	//用于压缩的Writer
	fwrite io.WriteCloser
	//最终数据是否采用了压缩
	flate bool
}

func newMessageWriter(ws *WebSocket, opcode int) *MessageWriter {
	var w MessageWriter
	w.init(ws, opcode)
	return &w
}

func (w *MessageWriter) init(ws *WebSocket, opcode int) {
	w.ws = ws
	w.opts = ws.opts
	w.opcode = opcode
	w.flate = false
	w.fwrite = nil

	//初始化非压缩的FrameWriter
	w.mframe.init(ws, opcode)
	//若websocket支持压缩，并且握手时协商采用压缩，则创建用于压缩的Writer：w.fwrite
	if ws.opts.compress_alloter != nil && ws.flateWrite == true {
		w.fframe.init(ws, opcode)
		flate_writer, err := ws.opts.compress_alloter.NewWriter(&w.fframe)
		if err == nil {
			w.fwrite = flate_writer
		}
	}
}

func (w *MessageWriter) reset(ws *WebSocket, opcode int) {
	w.ws = ws
	w.opts = ws.opts
	w.opcode = opcode

	//重置w.mframe
	w.mframe.Reset(ws, opcode)

	//重置w.fframe，将压缩结果的输出指向w.fframe
	if w.fwrite != nil {
		w.fframe.Reset(ws, opcode)
		ws.opts.compress_alloter.ResetWriter(w.fwrite, &w.fframe)
	}
}

func (w *MessageWriter) close() error {
	//若开启压缩，则优先调用w.fwrite.Close()， 将最后的数据写到w.fframe
	//然后调用w.fframe将最后的数据写到网络接口中
	//若没有使用压缩，则胡数据在w.mframe中，只需要调用w.mframe.Close()即可
	if w.flate == true && w.fwrite != nil {
		err := w.fwrite.Close()
		if err != nil {
			return err
		}
		err = w.fframe.Close()
		if err != nil {
			return err
		}
	} else {
		err := w.mframe.Close()
		if err != nil {
			return err
		}
	}

	w.ws = nil
	w.opts = nil
	w.flate = false
	return nil
}

func (w *MessageWriter) Close() error {
	return CloseWriter(w)
}

func (w *MessageWriter) Write(p []byte) (int, error) {
	//没有开启压缩发送，则使用FrameWriter
	if w.fwrite == nil {
		return w.mframe.Write(p)
	}

	//若数据大小<指定长度，则使用FrameWriter
	if w.flate == false {
		buff_len_total := w.mframe.pos + len(p)
		data_len_total := buff_len_total - maxFrameHeaderSize
		if data_len_total < w.opts.minCompressSize && buff_len_total < len(w.mframe.buff) {
			return w.mframe.Write(p)
		}
	}

	//若达到压缩条件，则首先创建flateWriter，
	//然后将w.frame的数据写到flateWriter。
	if w.flate == false {
		w.flate = true
		if w.mframe.GetPayloadLength() > 0 {
			err := WriteAllTo(w.mframe.GetPayload(), w.fwrite)
			if err != nil {
				return 0, err
			}
		}
		w.mframe.opcode = Frame_Null
		w.mframe.pos = maxFrameHeaderSize
	}

	return w.fwrite.Write(p)
}

func (w *MessageWriter) WriteAll(data []byte) error {
	for len(data) > 0 {
		nn, err := w.Write(data)
		if err != nil {
			return err
		}
		data = data[nn:]
	}
	return nil
}

func (w *MessageWriter) WriteString(str string) error {
	return w.WriteAll(StringToBytes(str))
}

func (w *MessageWriter) WriteControlFrame(data []byte) error {
	return w.mframe.WriteControlFrame(data)
}
