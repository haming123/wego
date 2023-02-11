package gows

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type FrameHeader struct {
	opcode  int
	payload int64
	maskKey [4]byte
	masked  bool
	flate   bool
	isfin   bool
}

type FrameReader struct {
	ws      *WebSocket
	header  FrameHeader
	buffr   *bufio.Reader
	extra   strings.Reader
	maskPos int
}

func newFrameReader(ws *WebSocket, br *bufio.Reader) *FrameReader {
	mr := &FrameReader{}
	mr.Init(ws, br)
	return mr
}

func (mr *FrameReader) Init(ws *WebSocket, br *bufio.Reader) {
	mr.ws = ws
	if br != nil && br.Size() == ws.opts.frameReadBuffSize {
		mr.buffr = br
	} else {
		mr.buffr = bufio.NewReaderSize(ws.cnn, ws.opts.frameReadBuffSize)
	}
	mr.maskPos = 0
	mr.header.opcode = Frame_Null
	mr.header.isfin = true
}

func (mr *FrameReader) Close() error {
	mr.maskPos = 0
	mr.header.opcode = Frame_Null
	mr.header.isfin = true
	return nil
}

func (mr *FrameReader) BufferSize() int {
	return mr.buffr.Size()
}

// 使用peekRead函数，不用进行字节数组的拷贝，有助于提升性能
// 前提是n < len(buffr)
func (mr *FrameReader) framePeekRead(n int) ([]byte, error) {
	p, err := mr.buffr.Peek(n)
	mr.buffr.Discard(len(p))
	return p, err
}

// 读取帧的header信息
func (mr *FrameReader) readFrameHeader() (FrameHeader, error) {
	var fr FrameHeader

	//读取数据1：读取第一个字节
	//用于获取FIN标志、操作码(Opcode)、压缩标志
	byte1, err := mr.buffr.ReadByte()
	if err != nil {
		return fr, err
	}

	//获取FIN标志、操作码(Opcode)、压缩标志
	fr.isfin = byte1&(1<<7) != 0
	fr.flate = byte1&(1<<6) != 0
	fr.opcode = int(byte1 & 0xf)
	if isControlFrame(fr.opcode) == false && isMessageFrame(fr.opcode) == false {
		err = errors.New("unknown opcode " + strconv.Itoa(fr.opcode))
		mr.ws.WiteCloseProtocolError(err)
		return fr, err
	} else if isControlFrame(fr.opcode) && fr.isfin == false {
		err = errors.New("control frame not final")
		mr.ws.WiteCloseProtocolError(err)
		return fr, err
	}

	//读取数据2：读取第二个字节
	//用于获取掩码标志、数据长度
	byte2, err := mr.buffr.ReadByte()
	if err != nil {
		return fr, err
	}

	//读取数据3：获取数据长度(Payload len)
	//如果数据长度Payload len在0-125之间，那么Payload len用7位表示足以，表示的数也就是净荷长度
	//如果数据长度Payload len等于126，接下来2字节表示的16位无符号整数才是这一帧的长度
	//如果数据长度Payload len等于127，接下来8字节表示的64位无符号整数才是这一帧的长度
	fr.payload = int64(byte2 & 0x7f)
	switch fr.payload {
	case 126:
		p, err := mr.framePeekRead(2)
		if err != nil {
			return fr, err
		}
		fr.payload = int64(binary.BigEndian.Uint16(p))
	case 127:
		p, err := mr.framePeekRead(8)
		if err != nil {
			return fr, err
		}
		fr.payload = int64(binary.BigEndian.Uint64(p))
	}
	//payload长度检查
	if fr.payload < 0 {
		err = errors.New(fmt.Sprintf("received negative payload length: %v", fr.payload))
		mr.ws.WiteCloseProtocolError(err)
		return fr, err
	} else if isControlFrame(fr.opcode) && fr.payload > maxControlFrameSize {
		err = errors.New("control frame length > 125")
		mr.ws.WiteCloseProtocolError(err)
		return fr, err
	}

	//读取数据4：获取掩码(Mask)标志，并读取掩码（4个字节）
	fr.masked = byte2&(1<<7) != 0
	if fr.masked {
		mr.maskPos = 0
		p, err := mr.framePeekRead(4)
		if err != nil {
			return fr, err
		}
		copy(fr.maskKey[:], p)
	}

	logPrintf4ws(mr.ws, "read frame: opcode=%d fin=%v len=%d flate=%v", fr.opcode, fr.isfin, fr.payload, fr.flate)
	return fr, nil
}

func (mr *FrameReader) handleControlFrame(header *FrameHeader) error {
	//读取控制帧的payload
	data, err := mr.framePeekRead(int(header.payload))
	if err != nil {
		return err
	}
	maskBytes(header.maskKey, 0, data)

	if header.opcode == Frame_Ping {
		//logPrint("recieved ping")
		return mr.ws.WritePong(data)
	} else if header.opcode == Frame_Pong {
		//logPrint("recieved pong")
		return nil
	}

	ce, err := parseClosePayload(data)
	if err != nil {
		return err
	}

	logPrint4ws(mr.ws, "recieved close frame")
	logPrint4ws(mr.ws, ce.Code, ce.Text)
	err = mr.ws.WiteCloseText(ce.Code, ce.Text)
	if err == nil {
		err = errors.New("close frame be sent")
	}
	return err
}

func (mr *FrameReader) readMessageHeader() (FrameHeader, error) {
	//只有消息已经读取结束后才可以读取下一个消息的header
	if mr.header.payload > 0 || mr.header.isfin == false {
		return mr.header, errors.New("current frame not final")
	}

	for {
		//读取帧的header
		header, err := mr.readFrameHeader()
		if err != nil {
			return header, err
		}

		//若是控制帧，则处理控制帧，并重新读取下一帧
		if isControlFrame(header.opcode) {
			err := mr.handleControlFrame(&header)
			if err != nil {
				return header, err
			}
			continue
		}

		//消息的开始帧不能是Frame_Continue
		if header.opcode == Frame_Continue {
			return header, errors.New("begin from a continuation frame")
		}

		mr.header = header
		return header, nil
	}
}

func (mr *FrameReader) ReadMessagePayload(p []byte) (int, error) {
	for {
		//读取帧数据
		if mr.header.payload > 0 {
			nn, err := mr.buffr.Read(p)
			mr.maskPos = maskBytes(mr.header.maskKey, mr.maskPos, p[:nn])
			mr.header.payload -= int64(nn)
			if mr.header.payload < 1 && mr.header.isfin == true {
				//logPrint("received data message")
				return nn, io.EOF
			}
			return nn, err
		}

		//若没有数据可读，并且是结束帧，则返回io.EOF
		if mr.header.payload < 1 && mr.header.isfin == true {
			//logPrint("received data message")
			return 0, io.EOF
		}

		//当前帧的数据已经读取完成，读取下一个帧的header
		header, err := mr.readFrameHeader()
		if err != nil {
			return 0, err
		}

		//若是控制帧，则处理控制帧，并重新读取下一帧数据
		if isControlFrame(header.opcode) {
			err := mr.handleControlFrame(&header)
			if err != nil {
				return 0, err
			}
			continue
		}

		//下一个帧必须是continuation帧
		if header.opcode != Frame_Continue {
			return 0, errors.New("next frame must be a continuation frame")
		}

		mr.header = header
	}
}

// mr.extra是结束标志："\x01\x00\x00\xff\xff"来防止flate.reader产生：unexpected EOF错误
// 首先从网络接口读取消息数据，消息数据读取完成后从mr.extra读取压缩数据结束标志
func (mr *FrameReader) Read(p []byte) (int, error) {
	nn, err := mr.ReadMessagePayload(p)
	if err == io.EOF && mr.extra.Len() > 0 {
		n, err := mr.extra.Read(p[nn:])
		nn += n
		return nn, err
	}
	return nn, err
}
