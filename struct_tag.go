package wego

import (
	"reflect"
	"strings"
)

type TagInfo struct {
	FieldName	string
	HasValue 	bool
	DefValue   	string
}

func SplitString(str string, sep string) (string, string, bool) {
	if str ==  "" {
		return "", "", false
	}

	index := strings.Index(str, sep)
	if index < 0 {
		return str, "", false
	}

	key := str[0:index]
	index += len(sep)
	val := str[index:]
	return key, val, true
}

func SplitAndTrim(str string, sep string) (string, string, bool) {
	key, val, has := SplitString(str, sep)
	key = strings.Trim(key, " ")
	val = strings.Trim(val, " ")
	return key, val, has
}

//获取struct的字段信息
func GetTagInfo(ff reflect.StructField, tag string) TagInfo {
	var info TagInfo
	info.FieldName = ff.Name

	tag_str := ff.Tag.Get(tag)
	key, val, has := SplitString(tag_str, ";")
	if key = strings.Trim(key, " "); key == "-" {
		info.FieldName = ""
		return info
	} else if key != "" {
		info.FieldName = key
	}

	if len(val) < 1 {
		return  info
	}

	key, val, has = SplitString(val, "=")
	if  key=strings.Trim(key, " "); key == "default" {
		info.HasValue = has
		info.DefValue = val
	}
	return info
}