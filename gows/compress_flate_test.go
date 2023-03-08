package gows

import (
	"bytes"
	"compress/flate"
	"io"
	"testing"
)

func TestFlateWriteAndRead(t *testing.T) {
	var buf bytes.Buffer
	flate_writer, _ := flate.NewWriter(&buf, flate.BestCompression)
	flate_writer.Write([]byte("12345"))
	flate_writer.Flush()
	buf.Write([]byte("\x01\x00\x00\xff\xff"))

	flate_reader := flate.NewReader(&buf)
	data, err := io.ReadAll(flate_reader)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestWebsocketFlateWriteAndRead(t *testing.T) {
	var cnn mocConn
	accept_options.UseFlate()
	accept_options.SetMinCompressSize(10)
	ws := newWebSocket(&cnn, &accept_options, nil)
	ws.useFlateWrite = true
	ws.WriteText([]byte("12345678901234567890"))
	opcode, data, err := ws.ReadMessage()
	if err != nil {
		t.Log(err)
	}
	t.Log("opcode=", opcode)
	t.Log(string(data))
}
