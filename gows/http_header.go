package gows

import (
	"net/http"
	"strings"
)

type TokenInfo struct {
	Value string
	Index int
	Flag  int
}

//header token:k1;k2=v2; k3, k4;k5="v5"
//token.Flag == 0 : ended with ';'
//token.Flag == 1 : ended with ','
//token.Index : is the positoin of ‘=’
func getNextHeaderToken(str string) (token TokenInfo, rest string) {
	colon := false
	token.Flag = 0
	token.Index = -1
	for i := 0; i < len(str); i++ {
		if colon == false && str[i] == '"' {
			colon = true
		} else if colon == true && str[i] == '"' {
			colon = false
		} else if colon == false && str[i] == '=' {
			if token.Index < 0 {
				token.Index = i
			}
		} else if colon == false && str[i] == ';' {
			token.Flag = 0
			token.Value = str[0:i]
			return token, str[i+1:]
		} else if colon == false && str[i] == ',' {
			token.Flag = 1
			token.Value = str[0:i]
			return token, str[i+1:]
		}
	}
	token.Value = str
	token.Flag = 1
	return token, ""
}

func HeaderValueCheck(header http.Header, key string, val string) bool {
	for _, header_str := range header[key] {
		var token TokenInfo
		for header_str != "" {
			token, header_str = getNextHeaderToken(header_str)
			tk := trimSpace(token.Value)
			if strings.EqualFold(tk, val) {
				return true
			}
		}
	}
	return false
}

type ExtParam struct {
	name   string
	params []string
}

// From RFC 6455:
//  Sec-WebSocket-Extensions = extension-list
//  extension-list = 1#extension
//  extension = extension-token *( ";" extension-param )
//  extension-token = registered-token
//  registered-token = token
//  extension-param = token [ "=" (token | quoted-string) ]
func parseHeaderExtensions(header http.Header) []ExtParam {
	var items []ExtParam
	for _, ext_str := range header["Sec-Websocket-Extensions"] {
		var ti TokenInfo
		var item ExtParam
		for item_beg := true; ext_str != ""; {
			ti, ext_str = getNextHeaderToken(ext_str)
			ti.Value = trimSpace(ti.Value)

			if item_beg == true {
				item.name = ti.Value
			} else {
				item.params = append(item.params, ti.Value)
			}

			item_beg = false
			if ti.Flag == 1 {
				items = append(items, item)
				item.name = ""
				item.params = item.params[0:0]
				item_beg = true
			}
		}
	}
	return items
}

//token:k1
//token:k1=v1
//token:k2="v2"
func parseExtensionParam(ti TokenInfo) (string, string) {
	//split string of k=v
	key := ti.Value
	val := ""
	if ti.Index >= 0 {
		key, val = ti.Value[:ti.Index], ti.Value[ti.Index+1:]
	}
	key = trimSpace(key)

	//trime ' ' and '"'
	val = trimSpace(val)
	val_len := len(val)
	if val_len > 1 && val[0] == '"' && val[val_len-1] == '"' {
		val = trimChar(val, '"')
	}

	return key, val
}
