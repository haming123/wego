package wmd

import (
	"bytes"
	"fmt"
	"strings"
)

type lineInfo struct {
	span_tag string
	span_beg int
	span_end int
}

func (h *lineInfo)reset() {
	h.span_tag = ""
	h.span_beg = 0
	h.span_end = 0
}

/*
Markdown中的转义字符为\，转义的有：
\\ 反斜杠
\` 反引号
\* 星号
\_ 下划线
\{\} 大括号
\[\] 中括号
\(\) 小括号
\# 井号
\+ 加号
\- 减号
\! 感叹号
\/
*/
var esc_chars []byte = []byte{ '\\', '`', '*', '_', '{', '}', '[', ']', '(', ')', '+', '-', '#', '!', '/'}
func writeEscapeText(data []byte, buff *bytes.Buffer) {
	dlen := len(data)
	for i:=0; i < dlen; i++ {
		if data[i] == '\\' && i < dlen-1  {
			ccc := data[i+1]
			for j :=0; j < len(esc_chars); j++ {
				if ccc != esc_chars[j] {
					continue
				}
				ppp := data[0:i]
				if len(ppp) > 0 {
					buff.Write(ppp)
				}
				buff.WriteByte(ccc)
				data = data[i+2:]
				dlen=len(data); i = -1
				break
			}
		}
	}
	if len(data) > 0 {
		buff.Write(data)
	}
}

/*
强调：星号与下划线都可以，单是斜体，双是粗体
- **加粗**
- *倾斜*
- ~~删除线~~
- `Code 标记`
*/
func writeTextLine(data []byte, buff *bytes.Buffer) []byte {
	var line_info lineInfo
	dlen := len(data)
	for i:=0; i < dlen; i++ {
		if i == dlen -1 {
			ppp := trimLineEnd(data[0:dlen])
			data = data[dlen:]
			if len(ppp) > 0 {
				buff.Write(ppp)
			}
			break
		} else if data[i] > 127 {
			continue
		} else if data[i] >= 97 && data[i] <= 122 {
			//小写字母
			continue
		} else if data[i] >= 65 && data[i] <= 90 {
			//大写字母
			continue
		} else if data[i] >= 48 && data[i] <= 57 {
			//数字
			continue
		} else if data[i] == '\n'  {
				ppp := trimLineEnd(data[0:i])
				data = data[i+1:]
				if len(ppp) > 0  {
					buff.Write(ppp)
				}
				break
		} else if data[i] == '\\' && i < dlen-1  {
			ccc := data[i+1]
			for j :=0; j < len(esc_chars); j++ {
				if ccc != esc_chars[j] {
					continue
				}
				ppp := data[0:i]
				if len(ppp) > 0 {
					buff.Write(ppp)
				}
				buff.WriteByte(ccc)
				data = data[i+2:]
				dlen=len(data); i = -1
				break
			}
		} else if data[i]=='*' && i<dlen-1 && data[i+1]=='*' &&line_info.span_tag == "" {
			//加粗开始：**...**
			line_info.span_tag = "strong"
			data = writeFrontAndSkip(data, buff, i, 2); dlen=len(data); i = -1
			buff.WriteString("<strong>")
		} else if data[i]=='*' && i<dlen-1 && data[i+1]=='*' && line_info.span_tag == "strong" {
			//加粗结束：**...**
			line_info.span_tag = ""
			data = writeFrontAndSkip(data, buff, i, 2); dlen=len(data); i = -1
			buff.WriteString("</strong>")
		} else if data[i]=='/' && i<dlen-1 && data[i+1]=='/' &&line_info.span_tag == "" {
			//倾斜开始：//...//
			line_info.span_tag = "em"
			data = writeFrontAndSkip(data, buff, i, 2); dlen=len(data); i = -1
			buff.WriteString("<em>")
		} else if data[i]=='/' && i<dlen-1 && data[i+1]=='/' && line_info.span_tag == "em" {
			//倾斜结束：//...//
			line_info.span_tag = ""
			data = writeFrontAndSkip(data, buff, i, 2); dlen=len(data); i = -1
			buff.WriteString("</em>")
		} else if data[i]=='~' && i<dlen-1 && data[i+1]=='~' && line_info.span_tag == "" {
			//删除线开始：~~...~~
			line_info.span_tag = "del"
			data = writeFrontAndSkip(data, buff, i, 2); dlen=len(data); i = -1
			buff.WriteString("<del>")
		} else if data[i]=='~' && i<dlen-1 && data[i+1]=='~' && line_info.span_tag == "del" {
			//删除线结束：~~...~~
			line_info.span_tag = ""
			data = writeFrontAndSkip(data, buff, i, 2); dlen=len(data); i = -1
			buff.WriteString("</del>")
		} else if data[i]=='_' && i<dlen-1 && data[i+1]=='_' && line_info.span_tag == "" {
			//下划线开始：__...__
			line_info.span_tag = "ins"
			data = writeFrontAndSkip(data, buff, i, 2); dlen=len(data); i = -1
			buff.WriteString("<ins>")
		} else if data[i]=='_' && i<dlen-1 && data[i+1]=='_' && line_info.span_tag == "ins" {
			//下划线结束：__...__
			line_info.span_tag = ""
			data = writeFrontAndSkip(data, buff, i, 2); dlen=len(data); i = -1
			buff.WriteString("</ins>")
		} else if data[i]=='[' && i<dlen-1 && data[i+1]=='[' && line_info.span_tag == "" {
			//class开始：[[name]]{{.....}}
			data = writeFrontAndSkip(data, buff, i, 0); dlen=len(data); i = -1
			class_name := ""
			data, class_name = getClassInfo(data, buff)
			if class_name != "" {
				line_info.span_tag = "span"
				buff.WriteString(fmt.Sprintf("<span class='%s'>", class_name))
			}
		} else if data[i]=='}' && i<dlen-1 && data[i+1]=='}' && line_info.span_tag == "span" {
			//class结束：[[name]]{{.....}}
			line_info.span_tag = ""
			data = writeFrontAndSkip(data, buff, i, 2); dlen=len(data); i = -1
			buff.WriteString("</span>")
		} else if data[i] == '`' && line_info.span_tag == "" {
			//Code标记开始：`code`
			line_info.span_tag = "code"
			data = writeFrontAndSkip(data, buff, i, 1); dlen=len(data); i = -1
			buff.WriteString("<code>")
		} else if data[i] == '`' && line_info.span_tag == "code" {
			//Code标记结束：`code`
			line_info.span_tag = ""
			data = writeFrontAndSkip(data, buff, i, 1); dlen=len(data); i = -1
			buff.WriteString("</code>")
		} else if data[i] == '[' && line_info.span_tag == "" {
			//链接开始：[文字](addr)
			writeFrontAndSkip(data, buff, i, 0)
			data, _ = write_link(data[i:], buff)
			dlen=len(data); i = -1
		} else if data[i] =='!' && i<dlen-1 && data[i+1]=='[' && line_info.span_tag == "" {
			//图片开始：![文字](addr)
			writeFrontAndSkip(data, buff, i, 0)
			data, _ = write_image(data[i:], buff)
			dlen=len(data); i = -1
		}
	}
	if line_info.span_tag != "" {
		buff.WriteString(fmt.Sprintf("</%s>", line_info.span_tag))
	}
	return data
}

//将POS前的内容写入buff，并忽略指定数量的字符
func writeFrontAndSkip(data []byte, buff *bytes.Buffer, pos int, skip int) []byte {
	ppp := data[0:pos]
	data = data[pos+skip:]
	if len(ppp) > 0 {
		buff.Write(ppp)
	}
	return data
}

//获取class的名称
func getClassInfo(data []byte, buff *bytes.Buffer) ([]byte, string) {
	var name []byte
	data_old := data
	data = data[2:]
	dlen := len(data)
	for i:=0; i < dlen; i++ {
		if data[i] == '\n' {
			break
		} else if data[i]==']' && i<dlen-1 && data[i+1]==']' {
			name = data[0:i]
			data = data[i+2:]
			break
		}
	}

	dlen = len(data)
	if name != nil && dlen>=2 && data[0]=='{' && data[1]=='{' {
		return data[2:],string(name)
	} else {
		buff.WriteString("[[")
		return data_old[2:], ""
	}
}

//[链接文字](链接地址 "链接title") title可加可不加
//处理 [baidu](http://baidu.com "xxx")
//<a href="http://baidu.com" title="xxx">baidu</a>
func write_link(data []byte, buff *bytes.Buffer) ([]byte, error) {
	var text []byte
	var address []byte
	var title []byte
	data_old := data
	data = data[1:]
	dlen := len(data)
	for i:=0; i < dlen; i++ {
		if data[i] == '\n' {
			break
		} else if text== nil && hasPrefix(data[i:], []byte("](")) {
			text = data[0:i]
			data = data[i+2:]; dlen = len(data); i=-1
		} else if text != nil && address == nil && data[i] == ' ' {
			address = data[0:i]
			data = data[i+1:]; dlen = len(data); i=-1
		} else if text != nil && address == nil && data[i] == ')' {
			address = data[0:i]
			data = data[i+1:]; dlen = len(data); i=-1
			break
		} else if address != nil && data[i] == ')' {
			title = data[0:i]
			data = data[i+1:]; dlen = len(data); i=-1
			break
		}
	}

	if text == nil || address == nil {
		buff.WriteString("[")
		return data_old[1:], nil
	} else {
		str_title := ""
		if title != nil {
			str_title = string(title)
			str_title = strings.Trim(str_title, "\"")
		}
		if len(str_title) > 0 {
			buff.WriteString(fmt.Sprintf("<a href='%s' title='%s'>", string(address), str_title))
		} else {
			buff.WriteString(fmt.Sprintf("<a href='%s'>", string(address)))
		}
		writeEscapeText(text, buff)
		buff.WriteString("</a>")
	}

	return data, nil
}

//![图片alt](图片地址 ''图片title'')
//图片alt就是显示在图片下面的文字，相当于对图片内容的解释。
//图片title是图片的标题，当鼠标移到图片上时显示的内容。title可加可不加
func write_image(data []byte, buff *bytes.Buffer) ([]byte, error) {
	var text []byte
	var address []byte
	var title []byte
	data_old := data
	data = data[2:]
	dlen := len(data)
	for i:=0; i < dlen; i++ {
		if data[i] == '\n' {
			break
		} else if text== nil && hasPrefix(data[i:], []byte("](")) {
			text = data[0:i]
			data = data[i+2:]; dlen = len(data); i=-1
		} else if text != nil && address == nil && data[i] == ' ' {
			address = data[0:i]
			data = data[i+1:]; dlen = len(data); i=-1
		} else if text != nil && address == nil && data[i] == ')' {
			address = data[0:i]
			data = data[i+1:]; dlen = len(data); i=-1
			break
		} else if address != nil && data[i] == ')' {
			title = data[0:i]
			data = data[i+1:]; dlen = len(data); i=-1
			break
		}
	}

	if text == nil || address == nil {
		buff.WriteString("![")
		return data_old[2:], nil
	} else {
		str_title := ""
		if title != nil {
			str_title = string(title)
			str_title = strings.Trim(str_title, "\"")
		}

		var img_info imageInfo
		addr := string(address); info := ""
		index := strings.Index(addr, "#")
		if index > 0 {
			info = addr[index+1:]
			addr = addr[0:index]
		}
		img_info.getImageAttr(info)

		style_img := ""
		if img_info.width != "" {
			style_img = fmt.Sprintf("width:%s", img_info.width)
		}
		if img_info.height != "" {
			style_img = fmt.Sprintf("%s;height:%s;", style_img, img_info.height)
		}

		if len(str_title) > 0 {
			buff.WriteString(fmt.Sprintf("<img src='%s' alt='%s' title='%s' style='%s' />",addr, text, str_title, style_img))
		} else {
			buff.WriteString(fmt.Sprintf("<img src='%s' alt='%s' style='%s' />", addr, text, style_img))
		}
	}

	return data, nil
}
