package wego

import (
	"testing"
)

func TestParseMediaType(t *testing.T) {
	ct := "multipart/form-data; boundary=WebKitFormBoundary7TMYhSONfkAM2z3a"
	t.Log(parseMediaType(ct))
	ct = "multipart/form-data ; boundary =WebKitFormBoundary7TMYhSONfkAM2z3a"
	t.Log(parseMediaType(ct))
}
