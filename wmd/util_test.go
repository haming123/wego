package wmd

import (
	"testing"
)

func TestGetBeforeChar(t *testing.T) {
	data := []byte("table=aaa:1;bbb:2")
	data1, data2 := splitBufferByChar(data, ':')
	t.Log(string(data1))
	t.Log(string(data2))
}

func TestGetBeforeChar2(t *testing.T) {
	data := []byte("table=aaa:1;bbb:2")
	data1, data2 := splitBufferByChar(data, '@')
	t.Log(string(data1))
	t.Log(len(data2))
}

func TestGetLineFromBuffer(t *testing.T) {
	data := []byte("table=aaa:1;bbb:2\n999")
	data1, data2 := getLineFromBuffer(data)
	t.Log(string(data1))
	t.Log(string(data2))
}

func TestGetLineFromBuffer2(t *testing.T) {
	data := []byte("table=aaa:1;bbb:2")
	data1, data2 := getLineFromBuffer(data)
	t.Log(string(data1))
	t.Log(len(data2))
}
