package worm

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type LogContex struct {
	Ctx         context.Context
	Session     *DbSession
	DbName     	string
	SqlType    	string
	Start       time.Time
	SQL         string
	Args        []interface{}
	Err         error
	Result      sql.Result
	ExeTime 	time.Duration
	IsTx      	bool
	UseStmt    	int
}

//打印SQL日志
type SqlPrintCB func(ctx *LogContex)
func print_sql_log(ctx *LogContex) {
	if ctx.Err != nil {
		debug_log.Output("[E]", ctx.GetPrintSQL())
		debug_log.Output("[E]", ctx.GetPrintResult())
	} else {
		debug_log.Output("[S]", ctx.GetPrintSQL())
		debug_log.Output("[S]", ctx.GetPrintResult())
	}
}

//生成执行结果日志
func (ctx *LogContex)GetPrintResult() string {
	if ctx.SqlType == "exec" {
		return ctx.get_exec_print_info()
	} else {
		return ctx.get_query_print_info()
	}
}

//生成db.Exec执行结果日志
func (ctx *LogContex)get_exec_print_info() string {
	//打印错误信息
	if ctx.Err != nil {
		return "failed: err=" + ctx.Err.Error()
	}
	if ctx.Result == nil {
		return "failed"
	}

	result_info := ""
	insert_id, _ := ctx.Result.LastInsertId()
	if insert_id > 0 {
		result_info = fmt.Sprintf(" insertId=%d", insert_id)
	} else {
		affected,_ := ctx.Result.RowsAffected()
		result_info = fmt.Sprintf(" affected=%d", affected)
	}

	//打印执行信息
	var buff bytes.Buffer
	if ctx.IsTx {
		str := fmt.Sprintf("TX: time=%0.3fms", float64(ctx.ExeTime.Nanoseconds())/float64(1e6))
		buff.WriteString(str)
	} else {
		str := fmt.Sprintf("DB: time=%0.3fms", float64(ctx.ExeTime.Nanoseconds())/float64(1e6))
		buff.WriteString(str)
	}
	buff.WriteString(result_info)
	if ctx.UseStmt == STMT_EXE_PREPARE {
		buff.WriteString("; prepare")
	} else if ctx.UseStmt == STMT_USE_CACHE {
		buff.WriteString("; stmt_cache")
	}

	return buff.String()
}

//生成db.Query执行结果日志
func (ctx *LogContex)get_query_print_info() string {
	//打印错误信息
	if ctx.Err != nil {
		return "failed: err=" + ctx.Err.Error()
	}

	//打印执行信息
	var buff bytes.Buffer
	if ctx.IsTx {
		str := fmt.Sprintf("TX: time=%0.3fms", float64(ctx.ExeTime.Nanoseconds())/float64(1e6))
		buff.WriteString(str)
	} else {
		str := fmt.Sprintf("DB: time=%0.3fms", float64(ctx.ExeTime.Nanoseconds())/float64(1e6))
		buff.WriteString(str)
	}
	if ctx.UseStmt == STMT_EXE_PREPARE {
		buff.WriteString("; prepare")
	} else if ctx.UseStmt == STMT_USE_CACHE {
		buff.WriteString("; stmt_cache")
	}
	if ctx.DbName != "" {
		buff.WriteString("; db=")
		buff.WriteString(ctx.DbName)
	}

	return buff.String()
}

//生成SQL请求日志
func (ctx *LogContex)GetPrintSQL() string {
	sql_tpl := ctx.SQL
	vals := ctx.Args
	max_field_len := ctx.Session.engine.max_log_field_len
	max_select_len := ctx.Session.engine.select_log_len
	show_pretty_log := ctx.Session.engine.show_pretty_log

	if show_pretty_log {
		if strings.HasPrefix(sql_tpl,"insert") {
			sql_tpl = parsel_debug_sql_insert(sql_tpl)
		}
	}

	sql_str := parsel_debug_sql_holder(sql_tpl, vals, max_field_len)

	//若select语句中的字段太多，则只显示部分字段，其他的以'...'代替
	if strings.HasPrefix(sql_str, "select") {
		if max_select_len > 0 {
			index := strings.Index(sql_str,"from")
			if index > max_select_len {
				sql_str1 := sql_str[0:max_select_len]
				sql_str2 := sql_str[index:]
				sql_str = sql_str1 + "... " + sql_str2
			}
		}
	}

	return sql_str
}

//将inser into table (name,name) values (val,val)
//修改为：inser into table set name=val,name=val的形式
func parsel_debug_sql_insert(sql_str string) string {
	ret_str := sql_str
	if !strings.HasPrefix(sql_str, "insert") {
		return ret_str
	}
	var buffer bytes.Buffer
	field_beg := strings.Index(sql_str,"(")
	if field_beg < 1 {
		return ret_str
	}
	buffer.WriteString(sql_str[0:field_beg])
	sql_str = sql_str[field_beg+1:]

	field_end := strings.Index(sql_str,")")
	if field_end < 1 {
		return ret_str
	}
	str_field := sql_str[0:field_end]
	sql_str = sql_str[field_end+1:]

	value_beg := strings.Index(sql_str,"(")
	if value_beg < 1 {
		return ret_str
	}
	sql_str = sql_str[value_beg+1:]

	value_end := strings.Index(sql_str,")")
	if value_end < 1 {
		return ret_str
	}
	str_value := sql_str[0:value_end]

	arr_field := strings.Split(str_field, ",")
	arr_value := strings.Split(str_value, ",")
	if len(arr_field) != len(arr_value) {
		return ret_str
	}

	buffer.WriteString("set ")
	for i, name := range arr_field {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(name)
		buffer.WriteString("=")
		buffer.WriteString(arr_value[i])
	}

	return buffer.String()
}

//将sql中的占位符替换为真正的值
//若字段的长度>max_value_len，则以'...'代替
func parsel_debug_sql_holder(sql_tpl string, vals []interface{}, max_value_len int) string {
	var buffer bytes.Buffer
	for i:=0; i < len(vals); i++ {
		index := strings.Index(sql_tpl, "?")
		if index < 0 {
			break;
		}
		txt_str := sql_tpl[0:index]
		sql_tpl = sql_tpl[index+1:]

		val := reflect.ValueOf(vals[i])
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		val_str := ""
		switch value := val.Interface().(type) {
		case string:
			val_str = fmt.Sprintf("'%v'", value)
		case time.Time:
			val_str = value.Format("'2006-01-02 15:04:05'")
		default:
			val_str = fmt.Sprintf("%v", value)
		}
		if max_value_len > 0 && len(val_str) > max_value_len {
			val_str = "'...'"
		}

		buffer.WriteString(txt_str)
		buffer.WriteString(val_str)
	}
	if len(sql_tpl) > 0 {
		buffer.WriteString(sql_tpl)
	}

	return buffer.String()
}
