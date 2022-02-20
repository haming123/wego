package wmd

type LineType int
const (
	LT_null  LineType = iota
	LT_blank	//空行
	LT_line 	//分割线
	LT_text 	//文本段落
	LT_title 	//标题
	LT_quote 	//引用
	LT_code  	//代码段
	LT_list  	//列表
	LT_ntext 	//数字列表
	LT_join 	//续行符
	LT_newl 	//换行符
)

func (md *Markdown)getLineType(data []byte) LineType {
	md.num_tab = 0
	md.num_space = 0
	md.line_beg = 0
	dlen := len(data)
	var first_char byte = 0
	for i:=0; i < dlen; i++ {
		//获取空格数量、\t的数量以及不是空格、\t、\r的第一个字符
		if first_char == 0 && data[i] == ' ' {
			md.num_space += 1
			continue
		} else if first_char == 0 && data[i] == '\t' {
			md.num_tab += 1
			continue
		} else if first_char == 0 && data[i] == '\r' {
			//若\r字符是第一个，忽略该字符
			continue
		}

		md.line_beg = i
		first_char = data[i]
		if first_char > 127 {
			//正文
			return LT_text
		} else if first_char == '\n' {
			//空行
			return LT_blank
		} else if first_char == '#' {
			//标题
			return LT_title
		} else if first_char == '>' {
			//引用
			return LT_quote
		} else if first_char == '<' {
			//续行符、换行符
			return line_type_join(md, data[i:])
		} else if first_char == '+' {
			//有序列表、正文
			return line_type_add(md, data[i:])
		} else if first_char == '-' {
			//分割线、正文
			return line_type_sub(md, data[i:])
		} else if first_char == '*' {
			//无序列表、正文
			return line_type_star(md, data[i:])
		} else if first_char >= '0' && first_char <= '9' {
			//数字列表、正文
			return line_type_number(md, data[i:])
		} else if hasPrefix(data[i:], []byte("```")) {
			//代码块
			return line_type_code(md, data[i:])
		} else {
			//正文
			return LT_text
		}
	}

	return LT_null
}

//"< "号, 作为续行符、新行符
func line_type_join(md *Markdown, data []byte) LineType {
	dlen := len(data)
	if dlen > 1 && data[1] == ' ' {
		//续行符
		return LT_join
	} else if dlen > 1 && data[1] == '<' {
		//换行符
		return LT_newl
	} else {
		return LT_text
	}
}

//"+"号作为有序列表
func line_type_add(md *Markdown, data []byte) LineType {
	var tag byte = '+'

	num_tt := 0
	dlen := len(data)
	for i:=0; i < dlen && i < 3; i++ {
		if data[i] == tag {
			num_tt += 1
		} else {
			break
		}
	}

	if num_tt == 1 && dlen > 1 && data[1] == ' ' {
		//有序列表
		md.list_type = 1
		return LT_list
	} else {
		return LT_text
	}
}

//"-"号, 作为水平分割线
func line_type_sub(md *Markdown, data []byte) LineType {
	var tag byte = '-'

	num_tt := 0
	dlen := len(data)
	for i:=0; i < dlen && i < 3; i++ {
		if data[i] == tag {
			num_tt += 1
		} else {
			break
		}
	}

	if num_tt >= 3 {
		//分割线
		return LT_line
	} else {
		return LT_text
	}
}

//"*"号, 作为无序列表
func line_type_star(md *Markdown, data []byte) LineType {
	var tag byte = '*'

	num_tt := 0
	dlen := len(data)
	for i:=0; i < dlen && i < 3; i++ {
		if data[i] == tag {
			num_tt += 1
		} else {
			break
		}
	}

	if num_tt == 1 && dlen > 1 && data[1] == ' ' {
		//无序列表
		md.list_type = 0
		return LT_list
	} else {
		return LT_text
	}
}

//"n. "数字列表
func line_type_number(md *Markdown, data []byte) LineType {
	num_tt := 0
	dlen := len(data)
	for i:=0; i < dlen && i < 10; i++ {
		if data[i] >= '0' && data[i] <= '9' {
			num_tt += 1
		} else {
			break
		}
	}

	if num_tt > 0 && dlen > num_tt+2 && data[num_tt] == '.' && data[num_tt+1] == ' ' {
		return LT_ntext
	} else if num_tt > 0 && dlen > num_tt+2 && data[num_tt] == ')' && data[num_tt+1] == ' ' {
		return LT_ntext
	} else {
		return LT_text
	}
}

//"`"号, 作为代码段开始
func line_type_code(md *Markdown, data []byte) LineType {
	num_tt := 0
	dlen := len(data)
	for i:=0; i < dlen && i < 3; i++ {
		if data[i] == '`' {
			num_tt += 1
		} else {
			break
		}
	}

	if num_tt == 3 {
		//代码段
		return LT_code
	} else {
		return LT_text
	}
}
