package wmd

import "bytes"

type listInfo struct {
	list_tt  	[8]byte 	//列表标签
	list_hh		int			//列表层次
}

func (li *listInfo)getCurTag() byte {
	if li.list_hh < 1 {
		return 0
	}
	return li.list_tt[li.list_hh-1]
}

//关闭根列表
func (li *listInfo)closeList(buff *bytes.Buffer) {
	li.closeListLevel(buff, 1)
}

//关闭到指定层次，>lh的层此会被关闭掉
func (li *listInfo)closeListLevel(buff *bytes.Buffer, lh int) {
	if lh < 1 {
		lh = 1
	}

	for i := li.list_hh; i >= lh; i-- {
		if li.list_tt[i-1] == 'o' {
			writeTabs(buff, i-1)
			buff.WriteString("</ol>\n")
		} else {
			writeTabs(buff, i-1)
			buff.WriteString("</ul>\n")
		}
	}

	li.list_hh = lh-1
}

func (li *listInfo)addListTag(buff *bytes.Buffer, tag byte, num_tab int) {
	//层级=\t的数量+1
	lh := num_tab + 1
	if lh > len(li.list_tt) {
		lh = len(li.list_tt)
	}

	//关闭深度>lh的层级
	if lh < li.list_hh {
		li.closeListLevel(buff, lh+1)
	}

	//同级列表，标签相同，则退出
	//同级列表，但标签不同，则关闭当前层级
	if li.list_hh == lh && tag == li.getCurTag() {
		return
	} else if li.list_hh == lh && tag != li.getCurTag() {
		li.closeListLevel(buff, lh)
	}

	//新增下级列表
	for i := li.list_hh; i < lh; i++ {
		if tag == 'o' {
			writeTabs(buff, i)
			buff.WriteString("<ol>\n")
		} else {
			writeTabs(buff, i)
			buff.WriteString("<ul>\n")
		}
	}
	li.list_hh = lh
	li.list_tt[lh-1] = tag
}

//列表
func (li *listInfo)writeList(buff *bytes.Buffer, data []byte, tag byte, num_tab int) []byte {
	li.addListTag(buff, tag, num_tab)
	writeTabs(buff, num_tab+1)
	buff.WriteString("<li>")
	data = writeTextLine(data, buff)
	data = writeContinue(data, buff)
	buff.WriteString("</li>\n")
	return data
}

