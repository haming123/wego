package worm

import (
	"database/sql"
	"errors"
	"reflect"
)

type ModelRows struct {
	*sql.Rows
	model     *DbModel
	scan_vals []interface{}
}

func (rows *ModelRows) Scan(ent_ptr interface{}) error {
	if ent_ptr == nil {
		return errors.New("ent_ptr must be reflect.Ptr")
	}
	//ent_ptr必须是一个指针
	v_ent := reflect.ValueOf(ent_ptr)
	if v_ent.Kind() != reflect.Ptr {
		return errors.New("ent_ptr must be reflect.Ptr")
	}
	//ent_ptr必须是一个结构体指针
	v_ent = reflect.Indirect(v_ent)
	if v_ent.Kind() != reflect.Struct {
		return errors.New("ent_ptr must be reflect.Ptr")
	}

	var md = rows.model
	if rows.scan_vals == nil {
		rows.scan_vals = md.get_scan_valus()
	}
	err := rows.Rows.Scan(rows.scan_vals...)
	if err != nil {
		return err
	}

	if hook, ok := md.ent_ptr.(AfterQueryInterface); ok {
		hook.AfterQuery(md.ctx)
	}

	if md.ent_ptr != ent_ptr {
		v_ent.Set(md.ent_value)
	}
	return nil
}
