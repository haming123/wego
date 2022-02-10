package worm

import (
	"database/sql"
	"errors"
	"reflect"
)

//执行行数据的scan
//将数据库查询的结果拷贝到desc对应的指针变量中
//在scan前将变量的指针包装为&FieldValue
//FieldValue实现了scanner接口用于接收数据库数据
//FieldValue能够处理字段为null的情况
func Scan(rows*sql.Rows, dest ...interface{}) error {
	values := make([]interface{}, len(dest))
	for i := 0; i < len(dest); i++ {
		fld := &FieldValue{"", dest[i], false}
		values[i] = fld
	}
	return rows.Scan(values...)
}

//将行数据保存到stuct对象中
func ScanModel(rows*sql.Rows, ent_ptr interface{}) error {
	v_ent := reflect.ValueOf(ent_ptr)
	if v_ent.Kind() != reflect.Ptr {
		return  errors.New("ent_ptr must be reflect.Ptr")
	}
	if v_ent.IsNil() {
		return  errors.New("ent_ptr is nil")
	}
	t_ent_base := GetDirectType(v_ent.Type())
	if t_ent_base.Kind() != reflect.Struct {
		return errors.New("ent base type muse be reflect.Struct")
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	//为行scan提供变量指针数组
	values := genScanAddr4Columns(columns, v_ent)
	//将行数据拷贝到变量中
	err = rows.Scan(values...)
	if err != nil {
		return err
	}

	hook, has_hook := ent_ptr.(AfterQueryInterface)
	if has_hook {
		hook.AfterQuery(nil)
	}
	return nil
}

//将数据库查询结构保存到struct数组中
//arr_ptr是struct数组的地址
func ScanModelArray(rows*sql.Rows, arr_ptr interface{}) error {
	v_arr := reflect.ValueOf(arr_ptr)
	if v_arr.Kind() != reflect.Ptr {
		return  errors.New("arr_ptr must be *Slice")
	}
	t_arr := GetDirectType(v_arr.Type())
	if t_arr.Kind() != reflect.Slice {
		return  errors.New("arr_ptr must be *Slice")
	}
	t_item := GetDirectType(t_arr.Elem())
	if t_item.Kind() != reflect.Struct {
		return errors.New("array item type muse be reflect.Struct")
	}

	//获取字段名称数组
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	num := 0
	//创建数组成员对象
	ent_ptr := reflect.New(t_item)
	v_ent := reflect.Indirect(ent_ptr)
	//获取准备scan的变量的地址
	scan_vals := genScanAddr4Columns(columns, v_ent)
	//获取数组对象
	v_arr_base := reflect.Indirect(v_arr)
	hook, has_hook := ent_ptr.Interface().(AfterQueryInterface)
	for rows.Next() {
		err = rows.Scan(scan_vals...)
		if err != nil {
			return err
		}

		if has_hook {
			hook.AfterQuery(nil)
		}

		v_arr_base.Set(reflect.Append(v_arr_base, v_ent))
		num ++
	}

	return nil
}

//将行数据保存到StringMap(map[string]string)中
func ScanStringRow(rows*sql.Rows) (StringRow, error) {
	row_data := make(StringRow)
	col_names, err := rows.Columns()
	if err != nil {
		return row_data, err
	}

	//创建string数组用于保存行数据
	//在scan前将string变量的指针包装为&FieldValue
	//FieldValue实现了scanner接口用于接收数据库数据
	//FieldValue能够处理字段为null的情况
	col_num := len(col_names)
	result := make([]string, col_num)
	var scan_vals = make([]interface{}, col_num)
	for i := 0; i < col_num; i++ {
		cell := FieldValue{"", &result[i], false}
		scan_vals[i] = &cell
	}

	err = rows.Scan(scan_vals...)
	if err != nil {
		return nil, err
	}

	//将stringv变量的值保存到map[string]string中
	for i := 0; i < col_num; i++ {
		row_data[col_names[i]] = result[i]
	}

	return row_data, nil
}

//将数据库查询结果保存到StringTable的表对象中
//StringTable是一个一维字符串数组
func ScanStringTable(rows*sql.Rows) (*StringTable, error) {
	rdata := &StringTable{}
	col_names, err := rows.Columns()
	if err != nil {
		return rdata, err
	}
	rdata.columns = append(rdata.columns, col_names...)

	//创建string数组用于保存行数据
	//在scan前将string变量的指针包装为&FieldValue
	//FieldValue实现了scanner接口用于接收数据库数据
	//FieldValue能够处理字段为null的情况
	col_num := len(col_names)
	result := make([]string, col_num)
	var scan_vals = make([]interface{}, col_num)
	for i := 0; i < col_num; i++ {
		cell := FieldValue{"", &result[i], false}
		scan_vals[i] = &cell
	}

	num := 0
	for rows.Next() {
		err = rows.Scan(scan_vals...)
		if err != nil {
			return nil, err
		}

		//将string变量中的数据添加到StringTable的数组中
		rdata.data = append(rdata.data, result...)
		num += 1
	}

	return rdata, nil
}