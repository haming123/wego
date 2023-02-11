package gows

func getMessageReader(ws *WebSocket) *MessageReader {
	fw, ok := ws.opts.messageReaderPool.Get().(*MessageReader)
	if !ok {
		//logPrint("new MessageReader !!!")
		return newMessageReader(ws)
	}
	//logPrint("get MessageReader from poll")
	fw.reset(ws)
	return fw
}

func putMessageReader(opts *AcceptOptions, w *MessageReader) {
	opts.messageReaderPool.Put(w)
	//logPrint("free MessageReader")
}

func (ws *WebSocket) NextReader() (FrameHeader, *MessageReader, error) {
	r := getMessageReader(ws)
	h, err := r.readMessageHeader()
	return h, r, err
}

func CloseReader(r *MessageReader) error {
	err := r.close()
	if err != nil {
		return err
	}
	putMessageReader(r.opts, r)
	return nil
}
