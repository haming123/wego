package wmd

import (
	"bytes"
	"fmt"
	"strings"
)

type imageInfo struct {
	width 		string
	height 		string
	align 		string
}

func (ti *imageInfo)getImageAttr(data string) {
	ti.width = ""
	ti.height = ""
	ti.align = "center"
	GetKeyValFromString(data, "&", "=", func(key string, val string) {
		key = strings.ToLower(key)
		if key == "width" {
			ti.width = val
		} else if key == "height" {
			ti.height = val
		} else if key == "align" {
			ti.align = val
		}
	})
}

/*
```image
	http:"//xxxxxx#align=center&width=80%
```
*/
func writeCodeImage(data []byte, buff *bytes.Buffer, attr string) []byte {
	var img_addr []byte
	dlen := len(data)
	for i:=0; i < dlen; i++ {
		if i == dlen -1 {
			ppp := data[0:dlen]
			data = data[dlen:]
			if len(ppp) > 0  {
				img_addr = trimLineEnd(ppp)
			}
			dlen = len(data);
			i = -1
		} else if data[i] == '\n' {
			ppp := data[0:i]
			data = data[i+1:]
			if len(ppp) > 0  {
				img_addr = trimLineEnd(ppp)
			}
			dlen = len(data);
			i = -1
		} else if hasPrefix(data[i:], []byte("```")) {
			data = data[i+3:]
			dlen = len(data);
			i = -1
			break
		}
	}

	var img_info imageInfo
	addr := string(img_addr); info := ""
	index := strings.Index(addr, "#")
	if index > 0 {
		info = addr[index+1:]
		addr = addr[0:index]
	}
	img_info.getImageAttr(info)

	style_p := ""
	if img_info.align != "" {
		style_p = fmt.Sprintf("text-align:%s;", img_info.align)
	}
	style_img := ""
	if img_info.width != "" {
		style_img = fmt.Sprintf("width:%s", img_info.width)
	}
	if img_info.height != "" {
		style_img = fmt.Sprintf("%s;height:%s;", style_img, img_info.height)
	}
	buff.WriteString(fmt.Sprintf("<p style='%s'>\n", style_p))
	buff.WriteString(fmt.Sprintf("\t<img src='%s' style='%s'/>\n", addr, style_img))
	buff.WriteString("</p>\n")
	return data
}
