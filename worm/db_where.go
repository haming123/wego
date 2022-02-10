package worm

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

//https://gitee.com/fifsky/gosql
type DbWhere struct {
	Tpl_sql string
	Values []interface{}
}

func SQLW(sql string, vals ...interface{}) *DbWhere {
	sqlw := new (DbWhere)
	sqlw.Init(sql, vals...)
	return sqlw
}

func (sqlw *DbWhere)GetTpl() string {
	return sqlw.Tpl_sql
}

func (sqlw *DbWhere)Reset() {
	sqlw.Tpl_sql = ""
	if sqlw.Values != nil {
		sqlw.Values = sqlw.Values[:0]
	}
}

func (sqlw *DbWhere)GetValues() []interface{} {
	return sqlw.Values
}

func (sqlw *DbWhere)Init(sql string, vals ...interface{}) *DbWhere {
	size := 10
	if size < len(vals) {size = len(vals)}
	if sqlw.Values == nil || cap(sqlw.Values) < size {
		sqlw.Values = make([]interface{}, size)
	}

	sqlw.Tpl_sql = sql
	sqlw.Values = append(sqlw.Values[:0], vals...)
	//sqlw.Values = append([]interface{}{}, vals...)
	return sqlw
}

func (sqlw *DbWhere)And(sql string, vals ...interface{}) *DbWhere {
	if len(sqlw.Tpl_sql) > 0 {
		sqlw.Tpl_sql += " and " + sql
	} else {
		sqlw.Tpl_sql = sql
	}
	sqlw.Values = append(sqlw.Values, vals...)
	return sqlw
}

func (sqlw *DbWhere)Or(sql string, vals ...interface{}) *DbWhere {
	if len(sqlw.Tpl_sql) > 0{
		sqlw.Tpl_sql += " or " + sql
	} else {
		sqlw.Tpl_sql = sql
	}
	sqlw.Values = append(sqlw.Values, vals...)
	return sqlw
}

func (sqlw *DbWhere)AndIf(cond bool, sql string, vals ...interface{}) *DbWhere {
	if cond == false {
		return sqlw
	}
	return sqlw.And(sql, vals...)
}

func (sqlw *DbWhere)OrIf(cond bool, sql string, vals ...interface{}) *DbWhere {
	if cond == false {
		return sqlw
	}
	return sqlw.Or(sql, vals...)
}

func (sqlw *DbWhere)AndExp(sqlw_sub *DbWhere) *DbWhere {
	sql := sqlw_sub.Tpl_sql
	if len(sqlw.Tpl_sql) > 0 {
		sql = "("  + sqlw_sub.Tpl_sql + ")"
	}
	return sqlw.And(sql, sqlw_sub.Values...)
}

func (sqlw *DbWhere)OrExp(sqlw_sub *DbWhere) *DbWhere {
	sql := sqlw_sub.Tpl_sql
	if len(sqlw.Tpl_sql) > 0 {
		sql = "("  + sqlw_sub.Tpl_sql + ")"
	}
	return sqlw.Or(sql, sqlw_sub.Values...)
}

func (sqlw *DbWhere)AndIn(field string, vals ...interface{}) *DbWhere {
	tpl, arr := parselParam4In(vals...)
	sql := fmt.Sprintf("%s in (%s)", field, tpl)
	return sqlw.And(sql, arr...)
}

func (sqlw *DbWhere)AndNotIn(field string, vals ...interface{}) *DbWhere {
	tpl, arr := parselParam4In(vals...)
	sql := fmt.Sprintf("%s not in (%s)", field, tpl)
	return sqlw.And(sql, arr...)
}

func (sqlw *DbWhere)OrIn(field string, vals ...interface{}) *DbWhere {
	tpl, arr := parselParam4In(vals...)
	sql := fmt.Sprintf("%s in (%s)", field, tpl)
	return sqlw.Or(sql, arr...)
}

func (sqlw *DbWhere)OrNotIn(field string, vals ...interface{}) *DbWhere {
	tpl, arr := parselParam4In(vals...)
	sql := fmt.Sprintf("%s not in (%s)", field, tpl)
	return sqlw.Or(sql, arr...)
}

func parselParam4In(vals ...interface{}) (string, []interface{}) {
	if len(vals) < 1 {
		return "", []interface{}{}
	}

	arr := vals
	if len(vals) == 1 {
		v_arr := reflect.ValueOf(vals[0])
		if v_arr.Kind() == reflect.Ptr {
			v_arr = v_arr.Elem()
		}
		if v_arr.Type().Kind() == reflect.Slice {
			num := v_arr.Len()
			arr = make([]interface{}, num)
			for i:=0; i < num; i++{
				arr[i] = v_arr.Index(i).Interface()
			}
		}
	}

	var buff bytes.Buffer
	for i:=0; i < len(arr); i++{
		if i > 0 {
			buff.WriteString(",")
		}
		buff.WriteString("?")
	}

	return buff.String(), arr
}

func (sqlw *DbWhere)GenPlatSQl() string {
	var buffer bytes.Buffer
	num := len(sqlw.Values)
	tpl_str := sqlw.Tpl_sql
	for i:=0; i < num; i++ {
		index := strings.Index(tpl_str, "?")
		if index < 0 {
			break;
		}
		txt_str := tpl_str[0:index]
		tpl_str = tpl_str[index+1:]

		val := reflect.ValueOf(sqlw.Values[i])
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		val_str := ""
		_, ok := val.Interface().(string)
		if ok {
			val_str = fmt.Sprintf("'%v'", val)
		} else {
			val_str = fmt.Sprintf("%v", val)
		}

		buffer.WriteString(txt_str)
		buffer.WriteString(val_str)
	}
	if len(tpl_str) > 0 {
		buffer.WriteString(tpl_str)
	}

	return buffer.String()
}

