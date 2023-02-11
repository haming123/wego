package gows

func getMessageWriter(ws *WebSocket, opcode int) *MessageWriter {
	fw, ok := ws.opts.messageWriterPool.Get().(*MessageWriter)
	if !ok {
		//logPrint("new MessageWriter !!!")
		return newMessageWriter(ws, opcode)
	}
	fw.reset(ws, opcode)
	//logPrint("get MessageWriter from poll")
	return fw
}

func putMessageWriter(ws *WebSocket, w *MessageWriter) {
	ws.opts.messageWriterPool.Put(w)
	//logPrint("free MessageWriter")
}

// 设置BeginMessage并发限制
// Frame_Binary、Frame_Text消息不允许并发发送
func lockNextWriter(ws *WebSocket, opcode int) {
	if opcode == Frame_Binary || opcode == Frame_Text {
		ws.ch_writer <- struct{}{}
	}
}

// 解除BeginMessage并发限制
// Frame_Binary、Frame_Text消息不允许并发发送
func unlockNextWriter(ws *WebSocket, opcode int) {
	if opcode == Frame_Binary || opcode == Frame_Text {
		<-ws.ch_writer
	}
}

// 获取一个MessageWriter， 若存在已经创建的MessageWriter，则该调用被阻塞
func (ws *WebSocket) NextWriter(opcode int) *MessageWriter {
	lockNextWriter(ws, opcode)
	return getMessageWriter(ws, opcode)
}

// 关闭MessageWriter，MessageWriter关闭后会解除对NextWriter调用的阻塞
func CloseWriter(writer *MessageWriter) error {
	ws := writer.ws
	opcode := writer.opcode
	defer func() {
		putMessageWriter(ws, writer)
		unlockNextWriter(ws, opcode)
	}()
	return writer.close()
}
