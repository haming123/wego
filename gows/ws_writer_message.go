package gows

import (
	"io"
)

type MessageWriter struct {
	ws     *WebSocket
	opts   *AcceptOptions
	opcode int

	//非压缩的数据写到frame_writer中
	frame_writer FrameWriter
	//压缩结果写到frame_flate中
	frame_flate FrameWriter
	//用于压缩的Writer
	flate_writer io.WriteCloser
	//最终数据是否采用了压缩
	flate_flag bool
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
	w.flate_flag = false
	w.flate_writer = nil

	//初始化非压缩的FrameWriter
	w.frame_writer.init(ws, opcode)
	//若AcceptOptions支持压缩，并且握手时协商采用压缩，则创建用于压缩的Writer：w.flate_writer
	if ws.opts.compress_alloter != nil && ws.useFlateWrite == true {
		w.frame_flate.init(ws, opcode)
		w.flate_writer, _ = ws.opts.compress_alloter.NewWriter(&w.frame_flate)
	}
}

func (w *MessageWriter) reset(ws *WebSocket, opcode int) {
	w.ws = ws
	w.opts = ws.opts
	w.opcode = opcode

	//重置w.frame_writer
	w.frame_writer.Reset(ws, opcode)
	//重置w.frame_flate以及 w.flate_writer
	if w.flate_writer != nil && ws.opts.compress_alloter != nil {
		w.frame_flate.Reset(ws, opcode)
		ws.opts.compress_alloter.ResetWriter(w.flate_writer, &w.frame_flate)
	}
}

func (w *MessageWriter) close() error {
	//若开启压缩，则先调用w.flate_writer.Flush()，将最后的数据写到w.frame_flate
	//然后调用w.frame_flate.Close()将最后的数据写到网络接口中
	//若没有使用压缩，则数据在w.frame_writer中，只需要调用w.frame_writer.Close()即可
	if w.flate_flag == true && w.flate_writer != nil {
		err := w.opts.compress_alloter.FlushWriter(w.flate_writer)
		if err != nil {
			return err
		}
		err = w.frame_flate.Close()
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
	w.flate_flag = false
	return nil
}

func (w *MessageWriter) Close() error {
	return CloseWriter(w)
}

func (w *MessageWriter) Write(p []byte) (int, error) {
	//没有开启压缩发送，则使用frame_writer
	if w.flate_writer == nil {
		return w.frame_writer.Write(p)
	}

	//若数据大小<指定长度，则使用frame_writer
	if w.flate_flag == false {
		buff_len_total := w.frame_writer.pos + len(p)
		data_len_total := buff_len_total - maxFrameHeaderSize
		if data_len_total < w.opts.minCompressSize && buff_len_total < len(w.frame_writer.buff) {
			return w.frame_writer.Write(p)
		}
	}

	//若达到压缩条件，则将w.frame_writer的数据写到flate_writer
	//然后清空frame_writer。
	if w.flate_flag == false {
		w.flate_flag = true
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
