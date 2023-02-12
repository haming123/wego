package gows

import (
	"bufio"
	"encoding/json"
	"net"
	"sync"
	"time"
)

const writeTimeOut = 10 * time.Second

var wsConnPool sync.Pool

type WebSocket struct {
	cnn  net.Conn
	opts *AcceptOptions

	mux       sync.Mutex
	ch_writer chan struct{}
	closed    bool

	flateWrite   bool
	writeTimeOut time.Duration

	msgReader *FrameReader
	handler   SocketHandler
}

func newWebSocket(cnn net.Conn, opts *AcceptOptions, br *bufio.Reader) *WebSocket {
	ch_writer := make(chan struct{}, 1)
	ws := &WebSocket{
		cnn:          cnn,
		opts:         opts,
		ch_writer:    ch_writer,
		writeTimeOut: writeTimeOut,
	}
	ws.msgReader = newFrameReader(ws, br)
	return ws
}

func (ws *WebSocket) Close() error {
	ws.mux.Lock()
	defer ws.mux.Unlock()

	if ws.closed == true {
		return nil
	}
	ws.closed = true

	err := ws.cnn.Close()
	if err != nil {
		return err
	}

	if ws.handler != nil {
		ws.handler.OnClose(ws)
	}

	logPrint4ws(ws, "websocket is closed!!!")
	return nil
}

func (ws *WebSocket) LocalAddr() net.Addr {
	return ws.cnn.LocalAddr()
}

func (ws *WebSocket) RemoteAddr() net.Addr {
	return ws.cnn.RemoteAddr()
}

func (ws *WebSocket) SetWriteTimeOut(writeTimeOut time.Duration) {
	ws.writeTimeOut = writeTimeOut
}

func (ws *WebSocket) WriteMessage(opcode int, data []byte) error {
	writer := ws.NextWriter(opcode)
	err := writer.WriteAll(data)
	if err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}

func (ws *WebSocket) WriteText(data []byte) error {
	return ws.WriteMessage(Frame_Text, data)
}

func (ws *WebSocket) WriteBinary(data []byte) error {
	return ws.WriteMessage(Frame_Binary, data)
}

func (ws *WebSocket) WriteJSON(v interface{}) error {
	writer := ws.NextWriter(Frame_Text)
	err := json.NewEncoder(writer).Encode(v)
	if err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}

func (ws *WebSocket) WriteClose(data []byte) error {
	writer := ws.NextWriter(Frame_Close)
	err := writer.WriteControlFrame(data)
	if err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}

func (ws *WebSocket) WritePing(data []byte) error {
	writer := ws.NextWriter(Frame_Ping)
	err := writer.WriteControlFrame(data)
	if err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}

func (ws *WebSocket) WritePong(data []byte) error {
	writer := ws.NextWriter(Frame_Pong)
	err := writer.WriteControlFrame(data)
	if err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}

func (ws *WebSocket) WiteCloseText(code CloseCode, text string) error {
	data, err := MarshalCloseInfo(code, text)
	if err != nil {
		return err
	}
	err = ws.WriteClose(data)
	if err != nil {
		return err
	}
	return nil
}

func (ws *WebSocket) WiteCloseError(code CloseCode, err error) error {
	return ws.WiteCloseText(code, err.Error())
}

func (ws *WebSocket) WiteCloseProtocolError(err error) error {
	return ws.WiteCloseError(CloseProtocolError, err)
}

func (ws *WebSocket) Serve(handler SocketHandler) {
	ws.handler = handler
	switch handler.(type) {
	case MessageHandler:
		messageReadLoop(ws, handler.(MessageHandler))
	case StreamReadHandler:
		streamReadLoop(ws, handler.(StreamReadHandler))
	default:
		panic("incorrect handler type")
	}
}

func (ws *WebSocket) ServeMessage(handler MessageHandler) {
	ws.handler = handler
	messageReadLoop(ws, handler)
}

func (ws *WebSocket) ServeStream(handler StreamReadHandler) {
	ws.handler = handler
	streamReadLoop(ws, handler)
}

/*
func (ws *WebSocket) ReadMessage() (int, []byte, error) {
	if ws.handler != nil {
		return 0, nil, errors.New("WebSocket.handler != nil")
	}

	head, reader, err := ws.NextReader()
	if err != nil {
		return head.opcode, nil, err
	}

	p, err := reader.ReadAll()
	reader.Close()

	return head.opcode, p.GetBytes(), err
}*/
