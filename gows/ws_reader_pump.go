package gows

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"
)

type SocketHandler interface {
	OnClose(ws *WebSocket)
}

type MessageHandler interface {
	SocketHandler
	OnMessage(ws *WebSocket, opcode int, buff *ByteBuffer) error
}

type ChuckReadHandler interface {
	SocketHandler
	OnData(ws *WebSocket, opcode int, fin bool, buff *ByteBuffer) error
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func onReadError(ws *WebSocket, err error) {
	//接收到客户端的关闭请求则直接关闭链接
	if err == errCloseFrame {
		ws.Close()
		return
	}

	err_code := ws.err_code
	if err_code < 1 {
		err_code = CloseProtocolError
	}

	//发送close控制帧
	err_write := ws.writeCloseFrame(err_code, err.Error())
	if err_write != nil && err_write != errWroteClose {
		ws.Close()
		return
	}

	//设置Handshake的响应时间
	readWait := ws.readTimeOut
	if readWait < 1 {
		readWait = 5
	}
	tm_wait := time.Now().Add(readWait)

	//读取Handshake响应。
	//若收到Handshake响应，则回读取到一个error：errCloseFrame
	for {
		var reader *MessageReader
		ws.cnn.SetReadDeadline(tm_wait)
		_, reader, err = ws.NextReader()
		if err != nil {
			reader.Close()
			break
		}

		_, err = reader.ReadAll()
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			reader.Close()
			break
		}
		reader.Close()
	}

	ws.Close()
}

func messageReadLoop(ws *WebSocket, handler MessageHandler) {
	var err error
	defer func() {
		logPrint4ws(ws, err)
		if err_recover := recover(); err_recover != nil {
			message := fmt.Sprintf("%s", err_recover)
			logPrint4ws(ws, trace(message))
		}
		onReadError(ws, err)
	}()

	for {
		var head FrameHeader
		var reader *MessageReader
		head, reader, err = ws.NextReader()
		if err != nil {
			reader.Close()
			return
		}

		var p *ByteBuffer
		p, err = reader.ReadAll()
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			reader.Close()
			return
		}

		err = handler.OnMessage(ws, head.opcode, p)
		if err != nil {
			if ws.err_code < 1 {
				ws.err_code = CloseUnsupportedData
			}
			reader.Close()
			return
		}

		reader.Close()
	}
}

func chunkReadLoop(ws *WebSocket, handler ChuckReadHandler) {
	var err error
	defer func() {
		if err_recover := recover(); err_recover != nil {
			message := fmt.Sprintf("%s", err_recover)
			logPrint4ws(ws, trace(message))
		}
		onReadError(ws, err)
	}()

	for {
		var head FrameHeader
		var reader *MessageReader
		head, reader, err = ws.NextReader()
		if err != nil {
			reader.Close()
			return
		}

		for {
			fin := false
			mb := GetByteBuffer(ws.opts)
			err := mb.ReadFull(reader)
			if err == io.EOF {
				err = nil
				fin = true
			}
			if err != nil {
				reader.Close()
				return
			}

			handler.OnData(ws, head.opcode, fin, mb)
			if err != nil {
				if ws.err_code < 1 {
					ws.err_code = CloseUnsupportedData
				}
				reader.Close()
				return
			}
			if fin == true {
				break
			}
		}

		reader.Close()
	}
}
