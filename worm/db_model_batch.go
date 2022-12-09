package worm

import (
	"database/sql"
	"errors"
	"reflect"
	"time"
)

type BatchResult struct {
	Count int64
	Err   error
}

func (this *BatchResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (this *BatchResult) RowsAffected() (int64, error) {
	return this.Count, this.Err
}

func (md *DbModel) BatchInsert(arr_ptr interface{}) (sql.Result, error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}

	var result BatchResult
	//获取变量arr_ptr的类型
	v_arr_ptr := reflect.ValueOf(arr_ptr)
	if v_arr_ptr.Kind() != reflect.Ptr {
		return &result, errors.New("arr_ptr must be *Slice")
	}
	//获取变量arr_ptr指向的类型
	t_arr := GetDirectType(v_arr_ptr.Type())
	if t_arr.Kind() != reflect.Slice {
		return &result, errors.New("arr_ptr must be *Slice")
	}
	//获取数组成员的类型
	t_item := GetDirectType(t_arr.Elem())
	if t_item.Kind() != reflect.Struct {
		return &result, errors.New("arr_ptr must be *Struct")
	}
	//数组长度<1，直接返回
	v_arr := reflect.Indirect(v_arr_ptr)
	arr_len := v_arr.Len()
	if arr_len < 1 {
		return &result, nil
	}

	//若t_item与Model类型不一致，则获取第一个数组成员的指针：v_item_ptr，并通过v_item_ptr来选择字段
	//若t_item是一个vo，item_ptr_vo指向：t_item，调用SaveToModel来给md.ent_value赋值，并选择字段
	//若t_item是一个eo，item_ptr_eo指向：t_item，则调用copyDataToModel来给md.ent_value赋值，不并选择字段
	var item_ptr_vo VoSaver = nil
	if t_item != md.ent_type {
		v_item_ptr := v_arr.Index(0).Addr()
		item_ptr := v_item_ptr.Interface()
		if vo_ptr, ok := item_ptr.(VoSaver); ok {
			item_ptr_vo = vo_ptr
			vo_ptr.SaveToModel(md, md.ent_ptr)
		} else {
			var v_item = v_item_ptr.Elem()
			copyDataToModel(md, v_item, md.ent_value)
		}
	}

	dbs := md.db_ptr
	sql_str := md.db_ptr.engine.db_dialect.GenModelInsertSql(md)
	vals := []interface{}{}
	for i, item := range md.flds_addr {
		if md.GetFieldFlag4Insert(i) == false {
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

	for i := 0; i < arr_len; i++ {
		item := v_arr.Index(i)
		if t_item == md.ent_type {
			md.ent_value.Set(item)
		} else if item_ptr_vo != nil {
			item_ptr_vo.SaveToModel(nil, md.ent_ptr)
		} else {
			copyDataToModel(md, item, md.ent_value)
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
