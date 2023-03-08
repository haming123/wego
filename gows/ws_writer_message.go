package gows

import (
	"io"
)

type MessageWriter struct {
	ws     *WebSocket
	opts   *AcceptOptions
	opcode int

	//非压缩的数据写到mframe中
	frame_writer FrameWriter
	//压缩结果写到fframe中
	flate_frame FrameWriter
	//用于压缩的Writer
	flate_writer io.WriteCloser
	//最终数据是否采用了压缩
	use_flate bool
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
	w.use_flate = false
	w.flate_writer = nil

	//初始化非压缩的FrameWriter
	w.frame_writer.init(ws, opcode)
	//若websocket支持压缩，并且握手时协商采用压缩，则创建用于压缩的Writer：w.flate_writer
	if ws.opts.compress_alloter != nil && ws.useFlateWrite == true {
		w.flate_frame.init(ws, opcode)
		w.flate_writer, _ = ws.opts.compress_alloter.NewWriter(&w.flate_frame)
	}
}

func (w *MessageWriter) reset(ws *WebSocket, opcode int) {
	w.ws = ws
	w.opts = ws.opts
	w.opcode = opcode

	//重置w.frame_writer
	w.frame_writer.Reset(ws, opcode)
	//重置w.flate_frame，将压缩结果的输出指向w.flate_frame
	if w.flate_writer != nil {
		w.flate_frame.Reset(ws, opcode)
		ws.opts.compress_alloter.ResetWriter(w.flate_writer, &w.flate_frame)
	}
}

func (w *MessageWriter) close() error {
	//若开启压缩，则优先调用w.flate_writer.Close()， 将最后的数据写到w.flate_frame
	//然后调用w.flate_frame将最后的数据写到网络接口中
	//若没有使用压缩，则胡数据在w.mframe中，只需要调用w.frame_writer.Close()即可
	if w.use_flate == true && w.flate_writer != nil {
		err := w.opts.compress_alloter.FlushWriter(w.flate_writer)
		if err != nil {
			return err
		}
		err = w.flate_frame.Close()
		if err != nil {
			return err
		}
	} else {
		err := w.frame_writer.Close()
		if err != nil {
			return err
		}
	}

	w.ws = nil
	w.opts = nil
	w.use_flate = false
	return nil
}

func (w *MessageWriter) Close() error {
	return CloseWriter(w)
}

func (w *MessageWriter) Write(p []byte) (int, error) {
	//没有开启压缩发送，则使用FrameWriter
	if w.flate_writer == nil {
		return w.frame_writer.Write(p)
	}

	//若数据大小<指定长度，则使用FrameWriter
	if w.use_flate == false {
		buff_len_total := w.frame_writer.pos + len(p)
		data_len_total := buff_len_total - maxFrameHeaderSize
		if data_len_total < w.opts.minCompressSize && buff_len_total < len(w.frame_writer.buff) {
			return w.frame_writer.Write(p)
		}
	}

	//若达到压缩条件，则首先创建flateWriter，
	//然后将w.frame的数据写到flateWriter。
	if w.use_flate == false {
		w.use_flate = true
		if w.frame_writer.GetPayloadLength() > 0 {
			err := WriteAllTo(w.frame_writer.GetPayload(), w.flate_writer)
			if err != nil {
				return 0, err
			}
		}
		w.frame_writer.opcode = Frame_Null
		w.frame_writer.pos = maxFrameHeaderSize
	}

	return w.flate_writer.Write(p)
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
	return w.frame_writer.WriteControlFrame(data)
}
