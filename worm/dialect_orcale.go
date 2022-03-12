package worm

import (
	"bytes"
	"fmt"
	"strings"
)

type dialectOrcale struct {
	DialectBase
}

func (db *dialectOrcale) GetName() string {
	return "oracle"
}

func (db *dialectOrcale) LimitSql(offset int64, limit int64) string  {
	return fmt.Sprintf(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY ", offset, limit)
}

func (db *dialectOrcale) ParsePlaceholder(sql_tpl string) string {
	tpl_str := sql_tpl
	var buffer bytes.Buffer
	for i:=0; i < len(sql_tpl); i++ {
		index := strings.Index(tpl_str, "?")
		if index < 0 {
			break;
		}
		txt_str := tpl_str[0:index]
		tpl_str = tpl_str[index+1:]
		bindvar := fmt.Sprintf(":%d", i+1)
		buffer.WriteString(txt_str)
		buffer.WriteString(bindvar)
	}
	if len(tpl_str) > 0 {
		buffer.WriteString(tpl_str)
	}
	return buffer.String()
}
