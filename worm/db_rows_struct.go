package worm

import (
	"database/sql"
	"errors"
	"reflect"
)

type StructRows struct {
	*sql.Rows
	base_ent  reflect.Value
	base_type reflect.Type
	scan_vals []interface{}
}

func (rows *StructRows) Scan(ent_ptr interface{}) error {
	if ent_ptr == nil {
		return errors.New("ent_ptr must be reflect.Ptr")
	}
	v_ent := reflect.ValueOf(ent_ptr)
	if v_ent.Kind() != reflect.Ptr {
		return errors.New("ent_ptr must be reflect.Ptr")
	}
	//ent_ptr必须是一个结构体指针
	v_ent = reflect.Indirect(v_ent)
	if v_ent.Kind() != reflect.Struct {
		return errors.New("ent_ptr must be reflect.Ptr")
	}

	t_ent := v_ent.Type()
	if rows.scan_vals == nil || rows.base_type != t_ent {
		//获取返回的数据库的全部字段
		columns, _ := rows.Columns()
		//创建一个对象，用于scan
		rows.base_type = t_ent
		rows.base_ent = reflect.Indirect(reflect.New(t_ent))
		//为行scan提供变量指针数组
		rows.scan_vals = genScanAddr4Columns(columns, rows.base_ent)
	}

	//将行数据拷贝到变量中
	err := rows.Rows.Scan(rows.scan_vals...)
	if err != nil {
		return err
	}
	v_ent.Set(rows.base_ent)

	hook, has_hook := ent_ptr.(AfterQueryInterface)
	if has_hook {
		hook.AfterQuery(nil)
	}
	return nil
}
