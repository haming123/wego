package worm

import (
	"database/sql"
	"errors"
)

type StringPair struct {
	Name  string
	Value string
}

type StringRows struct {
	*sql.Rows
	values    []*StringPair
	scan_vals []interface{}
}

func (rows *StringRows) Scan(ent_ptr *StringRow) error {
	if ent_ptr == nil {
		return errors.New("ent_ptr cannot equal nil")
	}

	//创建string数组用于保存行数据
	//在scan前将string变量的指针包装为&FieldValue
	//FieldValue实现了scanner接口用于接收数据库数据
	//FieldValue能够处理字段为null的情况
	if rows.values == nil {
		col_names, _ := rows.Columns()
		col_num := len(col_names)
		var values = make([]*StringPair, col_num)
		var scan_vals = make([]interface{}, col_num)
		for i := 0; i < col_num; i++ {
			val := StringPair{col_names[i], ""}
			cell := FieldValue{col_names[i], &val.Value, false}
			values[i] = &val
			scan_vals[i] = &cell
		}
		rows.values = values
		rows.scan_vals = scan_vals
	}

	err := rows.Rows.Scan(rows.scan_vals...)
	if err != nil {
		return err
	}

	//将stringv变量的值保存到map[string]string中
	for i := 0; i < len(rows.values); i++ {
		(*ent_ptr)[rows.values[i].Name] = rows.values[i].Value
	}
	return nil
}
