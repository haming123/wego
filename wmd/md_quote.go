package wmd

import "bytes"

type quoteInfo struct {
	level		int			//层次
}

//关闭根列表
func (qt *quoteInfo)closeQuote(buff *bytes.Buffer) {
	qt.closeQuoteLevel(buff, 1)
}

//关闭到指定层次，>lv的层此会被关闭掉
func (qt *quoteInfo)closeQuoteLevel(buff *bytes.Buffer, lv int) {
	if lv < 1 {
		lv = 1
	}
	for i := qt.level; i >= lv; i-- {
		writeTabs(buff, i-1)
		buff.WriteString("</blockquote>\n")
	}
	qt.level = lv-1
}

func (qt *quoteInfo)addQuoteTag(buff *bytes.Buffer, lv int) {
	//只能比当前层级深1层
	if lv > qt.level + 1 {
		lv = qt.level + 1
	}
	//关闭深度>lv的层级
	if lv < qt.level {
		qt.closeQuoteLevel(buff, lv+1)
	}
	//同级则退出
	if qt.level == lv {
		return
	}

	//新增一个下级列表
	qt.level = lv
	writeTabs(buff, qt.level - 1)
	buff.WriteString("<blockquote>\n")
}

//列表
func (qt *quoteInfo)writeQuote(buff *bytes.Buffer, data []byte, lv int) []byte {
	qt.addQuoteTag(buff, lv)
	writeTabs(buff, lv)
	buff.WriteString("<p>")
	data = writeTextLine(data, buff)
	data = writeContinue(data, buff)
	buff.WriteString("</p>\n")
	return data
}

