package wmd

import (
	"bytes"
	"fmt"
	"strings"
)

type Markdown struct {
	buff 		bytes.Buffer
	num_space 	int			//' '的数量
	num_tab 	int			//'\t'的数量
	line_beg 	int			//行的开始位置
	list_type 	int			//列表类型0 无序 1 有序
	list 		listInfo	//列表属性
	quote		quoteInfo	//引用属性
	btype 		LineType	//当前行类型
}

func (md *Markdown)GetBuff() * bytes.Buffer{
	return &md.buff
}

func (md *Markdown)closeTag() {
	if md.btype == LT_quote {
		md.quote.closeQuote(&md.buff)
	} else if md.btype == LT_list {
		md.list.closeList(&md.buff)
	}
}

func (md *Markdown)setBlockType(btype LineType) {
	if md.btype == btype {
		return
	}
	md.closeTag()
	md.btype = btype
}

func MarshalHtml(input []byte) []byte {
	var md Markdown
	md.marshalHtml(input)
	return md.buff.Bytes()
}

func (md *Markdown)marshalHtml(data []byte) {
	for len(data) > 0 {
		btype := md.getLineType(data)
		md.setBlockType(btype)
		switch btype {
		case LT_blank:
			//空行
			data = md.writeMdBlank(data[md.line_beg:])
		case LT_line:
			//分割线
			data = md.writeMdLine(data[md.line_beg:])
		case LT_title:
			//标题
			data = md.writeMdTitle(data[md.line_beg:])
		case LT_quote:
			//引用
			data = md.writeMdQuote(data[md.line_beg:])
		case LT_list:
			if md.list_type == 0 {
				//无序列表
				data = md.writeMdClist(data[md.line_beg:])
			} else {
				//有序列表
				data = md.writeMdNlist(data[md.line_beg:])
			}
		case LT_join:
			//续行符，输出文本
			data = md.writeMdText(data[md.line_beg+2:])
		case LT_newl:
			//换行符，输出文本
			data = md.writeMdText(data[md.line_beg+2:])
		case LT_code:
			//代码段
			data = md.writeMdCode(data[md.line_beg:])
		case LT_ntext:
			//数字列表
			data = md.writeMdNText(data[md.line_beg:])
		case LT_text:
			//文本段落
			data = md.writeMdText(data[md.line_beg:])
		}
	}
	md.closeTag()
}

//连续输出文本
func writeContinue(data []byte, buff *bytes.Buffer) []byte {
	var md_tmp Markdown
	for {
		btype := md_tmp.getLineType(data)
		if btype == LT_join {
			buff.WriteString("\n")
			data = writeTextLine(data[md_tmp.line_beg+2:], buff)
		} else if btype == LT_newl {
			buff.WriteString("</br>\n")
			num_space := md_tmp.num_space + md_tmp.num_tab*4
			for i:= 0; i < num_space; i++ {
				buff.WriteString("&nbsp;")
			}
			data = writeTextLine(data[md_tmp.line_beg+2:], buff)
			for {
				btype := md_tmp.getLineType(data)
				if btype == LT_join {
					buff.WriteString("\n")
					data = writeTextLine(data[md_tmp.line_beg+2:], buff)
				} else {
					break
				}
			}
		} else {
			break
		}
	}
	return data
}

//数字列表
func (md *Markdown)writeMdNText(data []byte) []byte {
	buff := &md.buff
	buff.WriteString("<div style='padding: 2px 0px'>")
	num_space := md.num_space + md.num_tab*4
	for i:= 0; i < num_space; i++ {
		buff.WriteString("&nbsp;")
	}
	data = writeTextLine(data, buff)
	data = writeContinue(data, buff)
	buff.WriteString("</div>\n")
	return data
}

//文本段落
func (md *Markdown)writeMdText(data []byte) []byte {
	buff := &md.buff
	buff.WriteString("<p>")
	num_space := md.num_space + md.num_tab*4
	for i:= 0; i < num_space; i++ {
		buff.WriteString("&nbsp;")
	}
	data = writeTextLine(data, buff)
	data = writeContinue(data, buff)
	buff.WriteString("</p>\n")
	return data
}

//空行
func (md *Markdown)writeMdBlank(data []byte) []byte {
	if data[0] == '\r' {
		data = data[1:]
	}
	if data[0] == '\n' {
		data = data[1:]
	}
	return data
}

//分割线
func (md *Markdown)writeMdLine(data []byte) []byte {
	tag := data[0]
	num_tag := 0; dlen := len(data)
	for i:=0; i < dlen; i++ {
		if data[i] == tag {
			num_tag += 1
		} else {
			break
		}
	}

	buff := &md.buff
	buff.WriteString("<hr/>\n")
	data = data[num_tag:]
	if data[0] == '\r' {
		data = data[1:]
	}
	if data[0] == '\n' {
		data = data[1:]
	}

	return data
}

//标题
func (md *Markdown)writeMdTitle(data []byte) []byte {
	//获取“#”号的数量，最多支持6个“#”号
	num_tt := 0
	dlen := len(data)
	for i:=0; i < dlen && i < 6; i++ {
		if data[i] == '#' {
			num_tt += 1
		} else {
			break
		}
	}

	buff := &md.buff
	buff.WriteString(fmt.Sprintf("<H%d>", num_tt))
	data = writeTextLine(data[num_tt:], buff)
	buff.WriteString(fmt.Sprintf("</H%d>\n", num_tt))
	return data
}

//引用
func (md *Markdown)writeMdQuote(data []byte) []byte {
	num_tt := 0
	dlen := len(data)
	for i:=0; i < dlen; i++ {
		if data[i] == '>' {
			num_tt += 1
		} else {
			break
		}
	}

	buff := &md.buff
	return md.quote.writeQuote(buff, data[num_tt:], num_tt)
}

//无序列表
func (md *Markdown)writeMdClist(data []byte) []byte {
	data = data[2:]
	num_tab := md.num_tab + md.num_space/4
	buff := &md.buff
	return md.list.writeList(buff, data, '*', num_tab)
}

//有序列表
func (md *Markdown)writeMdNlist(data []byte) []byte {
	data = data[2:]
	num_tab := md.num_tab + md.num_space/4
	buff := &md.buff
	return md.list.writeList(buff, data, 'o', num_tab)
}

//有序列表
func (md *Markdown)writeMdNlist2(data []byte) []byte {
	num_tt := 0
	dlen := len(data)
	for i:=0; i < dlen && i < 10; i++ {
		if data[i] >= '0' && data[i] <= '9' {
			num_tt += 1
		} else {
			break
		}
	}

	data = data[num_tt+2:]
	num_tab := md.num_tab + md.num_space/4
	buff := &md.buff
	return md.list.writeList(buff, data, 'o', num_tab)
}

//代码段
func (md *Markdown)writeMdCode(data []byte) []byte {
	var line []byte
	line, data = getLineFromBuffer(data[3:])

	lang_name := string(line)
	lang_name = strings.Trim(lang_name, " ")
	lang_name = strings.ToLower(lang_name)

	buff := &md.buff
	if strings.HasPrefix(lang_name, "table") {
		return writeCodeTable(data, buff, lang_name)
	} else if strings.HasPrefix(lang_name, "image") {
		return writeCodeImage(data, buff, lang_name)
	} else {
		return wiriteCodeContent(data, buff, lang_name)
	}
}
