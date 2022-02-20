package wmd

import (
	"bytes"
)

func writeHtml(buff *bytes.Buffer, data []byte) {
	dlen := len(data)
	for i:=0; i < dlen; i++ {
		if data[i] > 127 {
			continue
		} else if data[i] == '&' {
			data = writeFrontAndSkip(data, buff, i, 1);
			buff.WriteString("&amp;")
			dlen = len(data);
			i = -1
		} else if data[i] == '\'' {
			data = writeFrontAndSkip(data, buff, i, 1);
			buff.WriteString("&#39;")
			dlen = len(data);
			i = -1
		} else if data[i] == '<' {
			data = writeFrontAndSkip(data, buff, i, 1);
			buff.WriteString("&lt;")
			dlen = len(data);
			i = -1
		} else if data[i] == '>' {
			data = writeFrontAndSkip(data, buff, i, 1);
			buff.WriteString("&gt;")
			dlen = len(data);
			i = -1
		} else if data[i] == '"' {
			data = writeFrontAndSkip(data, buff, i, 1);
			buff.WriteString("&#34;")
			dlen = len(data);
			i = -1
		}
	}
	if len(data) > 0 {
		buff.Write(data)
	}
}

