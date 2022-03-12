package worm

import (
	"database/sql"
	"errors"
	"reflect"
	"time"
)

type BatchResult struct {
	Count 	int64
	Err 	error
}

func (this *BatchResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (this *BatchResult) RowsAffected() (int64, error) {
	return this.Count, this.Err
}

func (md *DbModel)BatchInsert(arr_ptr interface{}) (sql.Result, error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}

	var result BatchResult

	//获取变量arr_ptr的类型
	v_arr := reflect.ValueOf(arr_ptr)
	if v_arr.Kind() != reflect.Ptr {
		return &result, errors.New("arr_ptr must be *Slice")
	}
	//获取变量arr_ptr指向的类型
	t_arr := GetDirectType(v_arr.Type())
	if t_arr.Kind() != reflect.Slice {
		return  &result, errors.New("arr_ptr must be *Slice")
	}
	//获取数组成员的类型
	t_item := GetDirectType(t_arr.Elem())
	if t_item.Kind() != reflect.Struct {
		return &result, errors.New("arr_ptr must be *Struct")
	}

	dbs := md.db_ptr
	sql_str := md.db_ptr.engine.db_dialect.GenModelInsertSql(md)
	vals:= []interface{}{}
	for i, item := range md.flds_addr {
		if md.GetFieldFlag4Insert(i) == false  {
			continue
		}
		vals = append(vals, item.VAddr)
	}

	log_info := &LogContex{}
	log_info.Session = dbs
	log_info.Start = time.Now()
	log_info.SqlType = "exec"
	log_info.SQL = sql_str
	log_info.Args = nil
	log_info.Ctx = md.ctx

	tx, err := dbs.engine.db_raw.Begin()
	if err != nil {
		return &result, err
	}

	stmt, err := tx.Prepare(sql_str)
	if err != nil {
		tx.Rollback()
		return &result, err
	}
	defer stmt.Close()

	v_arr = reflect.Indirect(v_arr)
	v_base := reflect.Indirect(reflect.ValueOf(md.ent_ptr))
	num := v_arr.Len()
	for i:=0; i < num; i++ {
		item := v_arr.Index(i)

		ent_ptr := item.Addr().Interface()
		if item.Type() == v_base.Type() {
			v_base.Set(item)
		} else if ptr, ok := ent_ptr.(VoSaver); ok {
			ptr.SaveToModel(nil, md.ent_ptr)
		} else  {
			CopyDataToModel(nil, ent_ptr, md.ent_ptr)
		}

		_, err = stmt.Exec(vals...)
		if err != nil {
			tx.Rollback()
			return &result, err
		}
		result.Count += 1
	}
	tx.Commit()

	if dbs.need_print_sql_log(&md.SqlContex) {
		log_info.ExeTime = time.Now().Sub(log_info.Start)
		log_info.Result = &result
		log_info.Err = nil
		log_info.IsTx = true
		log_info.UseStmt = STMT_EXE_PREPARE
		dbs.engine.sql_print_cb(log_info)
	}
	return &result, nil
}

