package wmd

import (
	"bytes"
	"fmt"
	"strings"
)

type tableInfo struct {
	border 		bool
	header 		bool
	width		string
	align		string
	cells		[]string
}

func (ti *tableInfo)getTableAttr(data string) {
	ti.border = true
	ti.header = true
	GetKeyValFromString(data, "&", "=", func(key string, val string) {
		key = strings.ToLower(key)
		if key == "border" {
			val = strings.ToLower(strings.Trim(val, " "))
			if val == "off" {
				ti.border = false
			} else if val == "false" {
				ti.border = false
			} else if val == "0" {
				ti.border = false
			} else {
				ti.border = true
			}
		} else if key == "header" {
			val = strings.ToLower(strings.Trim(val, " "))
			if val == "off" {
				ti.header = false
			} else if val == "false" {
				ti.header = false
			} else if val == "0" {
				ti.header = false
			} else {
				ti.header = true
			}
		} else if key == "width" {
			ti.width = strings.ToLower(strings.Trim(val, " "))
		}
	})
}

/*
```table
表头|表头
:----|----:
单元格|单元格
单元格|单元格
```

表格属性：
表格居中：表头|表头
表格靠左：<表头|表头
表格靠右：表头|表头>
表格拉伸：<表头|表头>

单元格属性：
单元格靠左 :----
单元格靠右 ----:
单元格居中 ----
*/
func writeCodeTable(data []byte, buff *bytes.Buffer, header string) []byte {
	ti := &tableInfo{}
	lines := []string{}
	dlen := len(data)
	for i:=0; i < dlen; i++ {
		if i == dlen -1 {
			ppp := data[0:dlen]
			data = data[dlen:]
			if len(ppp) > 0  {
				str_line := strings.TrimSpace(string(trimLineEnd(ppp)))
				if str_line != "" {
					lines = append(lines, str_line)
				}
			}
			dlen = len(data);
			i = -1
		} else if data[i] == '\n'  {
			ppp := data[0:i]
			data = data[i+1:]
			if len(ppp) > 0  {
				str_line := strings.TrimSpace(string(trimLineEnd(ppp)))
				if str_line != "" {
					lines = append(lines, str_line)
				}
			}
			dlen = len(data);
			i = -1
		} else if hasPrefix(data[i:], []byte("```")) {
			ppp := data[0:i]
			if len(ppp) > 0  {
				str_line := strings.TrimSpace(string(trimLineEnd(ppp)))
				if str_line != "" {
					lines = append(lines, str_line)
				}
			}
			writeTable(ti, buff, lines)
			data = data[i+3:]
			dlen = len(data);
			i = -1
			break
		}
	}
	return data
}

func getTableInfo1(ti *tableInfo, line string) string {
	ti.width = ""
	ti.align = ""

	num := len(line)
	if num < 2 {
		return line
	}

	if line[0] =='<' && line[num-1] == '>' {
		ti.width = "width:100%"
		ti.align = ""
	} else if line[0] =='<' {
		ti.width = ""
		ti.align = "justify-content:flex-start "
	} else if line[num-1] == '>' {
		ti.width = ""
		ti.align = "justify-content:flex-end"
	} else {
		ti.width = ""
		ti.align = "justify-content:center"
	}

	line = strings.Trim(line, "<")
	line = strings.Trim(line, ">")
	return line
}

func getTableInfo2(ti *tableInfo, fields []string) {
	fnum := len(fields)
	if fnum < 1 {
		return
	}

	str1 := fields[0]
	if len(str1) < 1 {
		return
	}
	if str1[0] != ':' && str1[0] != '-' {
		return
	}

	ti.header = true
	for i:=0; i < fnum && i< len(ti.cells); i++ {
		field := fields[i]
		field = strings.Trim(field, " ")
		nn := len(field)
		if nn < 2 {
			continue
		}

		if field[0] == ':' {
			ti.cells[i] = "left"
		} else if field[nn-1] == ':' {
			ti.cells[i] = "right"
		}
	}
}

func writeTable(ti *tableInfo, buff *bytes.Buffer, lines []string) {
	row_num := len(lines)
	if row_num < 1 {
		return
	}

	lines[0] = getTableInfo1(ti, lines[0])
	row_1 := strings.Split(lines[0], "|")
	if len(row_1) < 1 {
		return
	}

	col_num := len(row_1)
	ti.cells = make([]string, col_num)
	for i:=0; i < len(ti.cells); i++ {
		ti.cells[i] = "center"
	}

	ti.header = false
	ti.border = true
	var row_2 []string
	if row_num > 1 {
		lines[1] = strings.Trim(lines[1], "|")
		row_2 = strings.Split(lines[1], "|")
		getTableInfo2(ti, row_2)
	}

	buff.WriteString(fmt.Sprintf("<div style='display:flex;%s'>\n", ti.align))
	if ti.border {
		buff.WriteString(fmt.Sprintf("<table border='1' style='border-spacing: 0;border-collapse: collapse;%s'>\n", ti.width))
	} else {
		buff.WriteString(fmt.Sprintf("<table border='0' style='border-spacing: 0;border-collapse: collapse;%s'>\n", ti.width))
	}

	if ti.header == true {
		writeTableHead(ti, buff, row_1, col_num)
		buff.WriteString("<tbody>\n")
		for i:=2; i < len(lines); i++ {
			lines[i] = strings.Trim(lines[i], "|")
			row_i := strings.Split(lines[i], "|")
			writeTableRow(ti, buff, row_i, col_num)
		}
		buff.WriteString("</tbody>\n")
	} else {
		buff.WriteString("<tbody>\n")
		writeTableRow(ti, buff, row_1, col_num)
		if row_2 != nil {
			writeTableRow(ti, buff, row_2, col_num)
		}
		for i:=2; i < len(lines); i++ {
			lines[i] = strings.Trim(lines[i], "|")
			row_i := strings.Split(lines[i], "|")
			writeTableRow(ti, buff, row_i, col_num)
		}
		buff.WriteString("</tbody>\n")
	}

	buff.WriteString("</table>\n")
	buff.WriteString("</div>\n")
}

func writeTableHead(ti *tableInfo, buff *bytes.Buffer, fields []string, col_num int) {
	i:=0
	buff.WriteString("<thead>\n\t<tr>\n")
	for ; i < len(fields); i++ {
		field := fields[i]
		buff.WriteString("\t\t<th>")
		buff.WriteString(field)
		buff.WriteString("</th>\n")
	}
	for ; i < col_num; i++ {
		buff.WriteString("\t\t<th></th>")
	}
	buff.WriteString("\t</tr>\n</thead>\n")
}

func writeTableRow(ti *tableInfo, buff *bytes.Buffer, fields []string, col_num int) {
	buff.WriteString("\t<tr>\n")
	i:=0
	for ; i < len(fields) && i < col_num; i++ {
		field := fields[i]
		align := ti.cells[i]
		tag_beg := fmt.Sprintf("\t\t<td align='%s'>", align)
		buff.WriteString(tag_beg)
		buff.WriteString(field)
		buff.WriteString("</td>\n")
	}
	for ; i < col_num; i++ {
		buff.WriteString("\t\t<td></td>")
	}
	buff.WriteString("\t</tr>\n")
}
