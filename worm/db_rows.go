package worm

import (
	"database/sql"
	"errors"
	"reflect"
)

type ModelRows struct {
	rows    *sql.Rows
	model 	*DbModel
}

func (rs *ModelRows) Next() bool {
	if rs.rows != nil {
		return rs.rows.Next()
	}
	return false
}

func (rs *ModelRows) Close() error {
	if rs.rows != nil {
		err := rs.rows.Close()
		if err != nil {
			return err
		}
		rs.rows = nil
	}
	return nil
}

func (rs *ModelRows)ScanModel(ent_ptr interface{}) error {
	v_ent := reflect.ValueOf(ent_ptr)
	if v_ent.Kind() != reflect.Ptr {
		return  errors.New("ent_ptr must be reflect.Ptr")
	}
	if v_ent.IsNil() {
		return  errors.New("ent_ptr is nil")
	}
	t_ent_base := GetDirectType(v_ent.Type())
	v_ent_base := reflect.Indirect(v_ent)
	if t_ent_base.Kind() != reflect.Struct {
		return errors.New("ent base type muse be reflect.Struct")
	}

	var md = rs.model
	vals:= md.get_scan_valus()
	err := rs.rows.Scan(vals...)
	if err != nil {
		return err
	}

	if hook, ok := md.ent_ptr.(AfterQueryInterface); ok {
		hook.AfterQuery(md.ctx)
	}

	if md.ent_ptr != ent_ptr {
		v_ent_base_table := reflect.Indirect(reflect.ValueOf(md.ent_ptr))
		v_ent_base.Set(v_ent_base_table)
	}

	return nil
}