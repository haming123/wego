package gows

import (
	"io"
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

func messageReadLoop(ws *WebSocket, handler MessageHandler) {
	defer func() {
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
