package gows

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

const g_writeTimeOut = 10 * time.Second
const g_readTimeOut = 10 * time.Second

var wsConnPool sync.Pool

type WebSocket struct {
	cnn  net.Conn
	opts *AcceptOptions

	mux        sync.Mutex
	ch_writer  chan struct{}
	closed     bool
	wroteClose bool

	useFlateWrite bool
	writeTimeOut  time.Duration
	readTimeOut   time.Duration

	msgReader *FrameReader
	handler   SocketHandler
}

func newWebSocket(cnn net.Conn, opts *AcceptOptions, br *bufio.Reader) *WebSocket {
	ch_writer := make(chan struct{}, 1)
	ws := &WebSocket{
		cnn:          cnn,
		opts:         opts,
		ch_writer:    ch_writer,
		writeTimeOut: g_writeTimeOut,
		readTimeOut:  g_readTimeOut,
	}
	ws.msgReader = newFrameReader(ws, br)
	return ws
}

func NewWebSocket(cnn net.Conn, opts *AcceptOptions, br *bufio.Reader) *WebSocket {
	return newWebSocket(cnn, opts, br)
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

func (ws *WebSocket) SetReadTimeOut(readTimeOut time.Duration) {
	ws.readTimeOut = readTimeOut
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

// Ping和Pong是websocket里的心跳，用来保证客户端是在线的，
// 目前浏览器中没有相关api发送ping给服务器，只能由服务器发ping给浏览器，浏览器返回pong消息。
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

var errWroteClose = errors.New("has wrote close frame")

// 连接任一端想关闭websocket，就发一个close frame给对端。
// 对端收到该frame后，若之前没有发过close frame，则必须回复一个close frame。
func (ws *WebSocket) WriteClose(data []byte) error {
	ws.mux.Lock()
	wroteClose := ws.wroteClose
	ws.wroteClose = true
	ws.mux.Unlock()
	if wroteClose {
		return errWroteClose
	}

	writer := ws.NextWriter(Frame_Close)
	err := writer.WriteControlFrame(data)
	if err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}

func (ws *WebSocket) WriteCloseText(code CloseCode, text string) error {
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

func (ws *WebSocket) WriteCloseError(code CloseCode, err error) error {
	return ws.WriteCloseText(code, err.Error())
}

func (ws *WebSocket) WriteCloseProtocolError(err error) error {
	return ws.WriteCloseError(CloseProtocolError, err)
}

func (ws *WebSocket) Serve(handler MessageHandler) {
	ws.handler = handler
	go messageReadLoop(ws, handler)
}

func (ws *WebSocket) ServeChunk(handler ChuckReadHandler) {
	ws.handler = handler
	go chunkReadLoop(ws, handler)
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
	if err == io.EOF {
		err = nil
	}
	reader.Close()

	return head.opcode, p.GetBytes(), err
}
*/

func (ws *WebSocket) CloseWebsocket(code CloseCode, info string) error {
	//发送close控制帧
	err := ws.WriteCloseText(code, info)
	if err != nil && err != errWroteClose {
		ws.Close()
		return err
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

	if err == nil {
		err = errors.New("read time out")
	} else if err == errCloseFrame {
		err = nil
	}
	if err != nil {
		ws.Close()
		return err
	}

	err = ws.Close()
	return err
}
