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

// 关闭帧可能包含数据部分（应用数据帧），该部分表明了关闭的原因，例如端点关闭、端点接收帧过大或端点收到的帧不符合预期。
// 如果有数据部分，则数据的前两个字节必须是一个无符号整数（网络字节序），该无符号整数表示了一个状态码，具体定义哪些关闭码将在后面的文章中介绍。
// 在无符号整数后面，可能还有一个UTF-8编码的数据，表示关闭原因，关闭原因由开发者自行定义（可选），并无规范。
// 关闭原因并不一定是对人可读的，但会对调试或传递相关信息起到一定的作用。由于数据不能保证可读，所以客户端不应将其显示给用户（会在关闭事件onclose中）。
//
// 应用程序在发送了一个关闭帧后，禁止再发送任何数据（此时处于CLOSING状态）。
// 如果端点（客户端或服务器）收到了一个关闭帧，并且之前没有发送过关闭帧，则端点必须发送一个关闭帧作为响应。当端点可以发送关闭响应时应尽快发送关闭响应。
// 一个端点可以延迟发送响应直到它的当前消息发送完毕（例如，已经发送了大多数的消息片段，则端点可能会在发送关闭响应帧前先将剩下的消息帧发送出去）。
// 但不能保证对方在已经发送了关闭帧后还能够继续处理这些数据。
// 在双方都以发送并接收了关闭帧后，端点需要断掉WebSocket连接并且必须关闭底层的TCP连接。服务器必须立即切断底层TCP连接，
// 客户端最好等待服务器断开连接，但也可以在发送并接收了关闭帧后任何时候断开连接，例如在一段时间内服务器仍没有断开TCP连接。
// 如果服务器和客户端同时发送了关闭帧，两端都会接收关闭帧，并且都需要断开TCP连接。
func (ws *WebSocket) WriteClose(data []byte) error {
	writer := ws.NextWriter(Frame_Close)
	err := writer.WriteControlFrame(data)
	if err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}

// Ping和Pong是websocket里的心跳，用来保证客户端是在线的，
func (ws *WebSocket) WritePing(data []byte) error {
	writer := ws.NextWriter(Frame_Ping)
	err := writer.WriteControlFrame(data)
	if err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}

// 当接收到 0x9 Ping 操作码的控制帧以后，应当立即发送一个包含 pong 操作码的帧响应。
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
