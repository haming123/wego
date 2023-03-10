package gows

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

type SocketHandler interface {
	OnClose(ws *WebSocket)
}

type MessageHandler interface {
	SocketHandler
	OnMessage(ws *WebSocket, opcode int, buff *ByteBuffer) error
}

type StreamReadHandler interface {
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

func messageReadLoop(ws *WebSocket, handler MessageHandler) {
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			logPrint4ws(ws, trace(message))
		}
		ws.Close()
	}()

	for {
		head, reader, err := ws.NextReader()
		if err != nil {
			reader.Close()
			logPrint4ws(ws, err)
			return
		}

		p, err := reader.ReadAll()
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			reader.Close()
			logPrint4ws(ws, err)
			return
		}

		err = handler.OnMessage(ws, head.opcode, p)
		if err != nil {
			reader.Close()
			logPrint4ws(ws, err)
			return
		}

		reader.Close()
	}
}

func streamReadLoop(ws *WebSocket, handler StreamReadHandler) {
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			logPrint4ws(ws, trace(message))
		}
		ws.Close()
	}()

	for {
		head, reader, err := ws.NextReader()
		if err != nil {
			reader.Close()
			logPrint4ws(ws, err)
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
				logPrint4ws(ws, err)
				return
			}

			handler.OnData(ws, head.opcode, fin, mb)
			if err != nil {
				reader.Close()
				logPrint4ws(ws, err)
				return
			}
			if fin == true {
				break
			}
		}

		reader.Close()
	}
}
