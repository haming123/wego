package wmd

import (
	"bytes"
	"fmt"
)

type codeInfo struct {
	num_space 	int		//' '的数量
	num_ss_t 	int		//'\t'的数量
}

func (h *codeInfo)reset() {
	h.num_ss_t = 0
	h.num_space = 0
}

//https://highlightjs.org/usage/
//HTML里的pre元素，可定义预格式化的文本。在pre元素中的文本会保留空格和换行符。
//可以导致段落断开的标签（例如标题、<p> 和 <address> 标签）绝不能包含在 <pre> 所定义的块里。
//当pre元素来展示源代码的时候最好的方式是用code元素来包裹代码，这样既可以保持格式又可以代表语义
func wiriteCodeContent(data []byte, buff *bytes.Buffer, lang_name string) []byte {
	buff.WriteString(fmt.Sprintf("<pre><code class='language-%s'>\n", lang_name))

	dlen := len(data)
	for i:=0; i < dlen; i++ {
		if hasPrefix(data[i:], []byte("```")) {
			ppp := data[0:i]
			writeHtml(buff, ppp)
			data = data[i+3:]
			break
		}
	}
	buff.WriteString("</code></pre>\n")
	return data
}
