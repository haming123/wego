package wmd

import (
	"bytes"
	"strings"
)

//添加指定数量的：\t
func writeTabs(buff *bytes.Buffer, num_tab int)  {
	for i := 0; i < num_tab; i++ {
		buff.WriteString("\t")
	}
}

//是否存在指定的前缀字符
func hasPrefix(data []byte, prefix []byte) bool {
	if len(data) < len(prefix) {
		return false
	}
	for i:=0; i < len(prefix); i++ {
		if data[i] != prefix[i] {
			return false
		}
	}
	return true
}

//去掉line后面的\r\n
func trimLineEnd(data []byte) ([]byte) {
	dlen := len(data)
	if dlen < 1 {
		return data
	} else if dlen > 1 && data[dlen-2] == '\r' && data[dlen-1] == '\n' {
		return data[0:dlen-2]
	} else if dlen > 0 && data[dlen-1] == '\r' {
		return data[0:dlen-1]
	} else if dlen > 0 && data[dlen-1] == '\n' {
		return data[0:dlen-1]
	} else {
		return data
	}
}

//返回cc前数据以及cc后的数据
//若不存在cc，全部数据作为front的数据
func splitBufferByChar(data []byte, cc byte) ([]byte, []byte) {
	dlen := len(data)
	for i:=0; i < dlen; i++ {
		if data[i] == cc {
			front := data[0:i]
			data = data[i+1:]
			return front, data
		}
	}
	front := data[0:]
	data  = data[dlen:]
	return front, data
}

//从buff中获取一行数据
//返回\n前数据以及\n后的数据
//若不存在\n，全部数据作为line的数据
func getLineFromBuffer(buff []byte) ([]byte, []byte) {
	line, data := splitBufferByChar(buff, '\n')
	return trimLineEnd(line), data
}

func getKeyString(key []byte) string {
	if key == nil && len(key) < 1 {
		return ""
	}
	key_str := string(key)
	key_str = strings.Trim(key_str, " ")
	return key_str
}

func getValString(val []byte) string {
	if val == nil && len(val) < 1 {
		return ""
	}
	return string(val)
}

//k:v;k:v.....\n
//cs: item split char
//ce：key value split char
type SetKeyVal func(string, string)
func GetKeyValFromBuff(data []byte, cs byte, ce byte, fn SetKeyVal) {
	var key []byte
	pos := 0; dlen := len(data)
	for i:=0; i < dlen; i++ {
		if i== dlen - 1 {
			val := data[pos:dlen]
			if key == nil {
				key = val
				val = nil
			}
			fn(getKeyString(key), getValString(val))
			pos = dlen - 1
			key = nil
		} else if data[i] == cs {
			val := data[pos:i]
			if key == nil {
				key = val
				val = nil
			}
			fn(getKeyString(key), getValString(val))
			pos = i+1
			key = nil
		} else if data[i] == ce {
			key = data[pos:i]
			pos = i+1
		}
	}
}

func GetKeyValFromString(query string, cs string, ce string, fn SetKeyVal)  {
	for query != "" {
		key := query
		if i := strings.Index(key, cs); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}

		value := ""
		if i := strings.Index(key, ce); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		fn(key, value)
	}
}

