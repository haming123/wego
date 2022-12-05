package worm

import (
	"errors"
	"reflect"
)

func (md *DbModel) GetFieldFlag4Insert(i int) bool {
	if md.flds_addr[i].Flag == false {
		return false
	}
	if md.flds_info[i].AutoIncr == true {
		return false
	}
	if md.flds_info[i].NotInsert == true {
		return false
	}
	return true
}

func (md *DbModel) get_fieldaddr_insert() []interface{} {
	cc := 0
	for i, _ := range md.flds_addr {
		if md.GetFieldFlag4Insert(i) == false {
			continue
		}
		cc += 1
	}

	index := 0
	vals := make([]interface{}, cc)
	for i, _ := range md.flds_addr {
		if md.GetFieldFlag4Insert(i) == false {
			continue
		}
		vals[index] = md.flds_addr[i].VAddr
		index += 1
	}
	return vals
}

func (md *DbModel) insertWithOutput() (int64, error) {
	if md.Err != nil {
		return 0, md.Err
	}

	if md.flds_ent != nil {
		genSelectionByFieldIndex(md, md.flds_ent)
	}

	if hook, ok := md.ent_ptr.(BeforeInsertInterface); ok {
		hook.BeforeInsert(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelInsertSql(md)
	values := md.get_fieldaddr_insert()
	rows, err := md.db_ptr.ExecQuery(&md.SqlContex, sql_str, values...)
	if err != nil {
		return 0, err
	}

	if hook, ok := md.ent_ptr.(AfterInsertInterface); ok {
		hook.AfterInsert(md.ctx)
	}

	if !rows.Next() {
		rows.Close()
		return 0, nil
	}

	var id int64 = 0
	err = rows.Scan(&id)
	if err != nil {
		rows.Close()
		return 0, err
	}
	rows.Close()
	return id, nil
}

func (md *DbModel) exec_insert() (int64, error) {
	if md.db_ptr.engine.db_dialect.ModelInsertHasOutput(md) {
		return md.insertWithOutput()
	}

	if md.Err != nil {
		return 0, md.Err
	}

	if md.flds_ent != nil {
		genSelectionByFieldIndex(md, md.flds_ent)
		md.flds_ent = nil
	}

	if hook, ok := md.ent_ptr.(BeforeInsertInterface); ok {
		hook.BeforeInsert(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelInsertSql(md)
	values := md.get_fieldaddr_insert()
	res, err := md.db_ptr.ExecSQL(&md.SqlContex, sql_str, values...)
	if err != nil {
		return 0, err
	}

	if hook, ok := md.ent_ptr.(AfterInsertInterface); ok {
		hook.AfterInsert(md.ctx)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (md *DbModel) Insert(args ...interface{}) (int64, error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return 0, md.Err
	}

	if len(args) > 1 {
		return 0, errors.New("arg number can not great 1")
	}

	if len(args) < 1 {
		return md.exec_insert()
	}

	ent_ptr := args[0]
	if ent_ptr == nil {
		return 0, errors.New("ent_ptr must be *Struct")
	}
	//ent_ptr必须是一个指针
	v_ent := reflect.ValueOf(ent_ptr)
	if v_ent.Kind() != reflect.Ptr {
		return 0, errors.New("ent_ptr must be *Struct")
	}
	//ent_ptr必须是一个结构体指针
	v_ent = reflect.Indirect(v_ent)
	if v_ent.Kind() != reflect.Struct {
		return 0, errors.New("ent_ptr must be *Struct")
	}

	v_base := reflect.Indirect(reflect.ValueOf(md.ent_ptr))
	if v_ent.Type() == v_base.Type() {
		v_base.Set(v_ent)
	} else if ptr, ok := ent_ptr.(VoSaver); ok {
		ptr.SaveToModel(md, md.ent_ptr)
	} else {
		CopyDataToModel(md, ent_ptr, md.ent_ptr)
	}

	return md.exec_insert()
}

func (md *DbModel) GetFieldFlag4Update(i int) bool {
	if md.flds_addr[i].Flag == false {
		return false
	}
	if md.flds_info[i].AutoIncr == true {
		return false
	}
	if md.flds_info[i].NotUpdate == true {
		return false
	}
	return true
}

func (md *DbModel) get_fieldaddr_update() []interface{} {
	cc := 0
	for i, _ := range md.flds_addr {
		if md.GetFieldFlag4Update(i) == false {
			continue
		}
		cc += 1
	}

	index := 0
	vals := make([]interface{}, cc)
	for i, _ := range md.flds_addr {
		if md.GetFieldFlag4Update(i) == false {
			continue
		}
		vals[index] = md.flds_addr[i].VAddr
		index += 1
	}
	return vals
}

func (md *DbModel) exec_update() (int64, error) {
	if md.Err != nil {
		return 0, md.Err
	}

	if md.flds_ent != nil {
		genSelectionByFieldIndex(md, md.flds_ent)
		md.flds_ent = nil
	}

	if len(md.db_where.Tpl_sql) < 1 {
		return 0, errors.New("no where clause")
	}

	if hook, ok := md.ent_ptr.(BeforeUpdateInterface); ok {
		hook.BeforeUpdate(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelUpdateSql(md)
	values := md.get_fieldaddr_update()
	values = append(values, md.db_where.Values...)
	res, err := md.db_ptr.ExecSQL(&md.SqlContex, sql_str, values...)
	if err != nil {
		return 0, err
	}

	if hook, ok := md.ent_ptr.(AfterUpdateInterface); ok {
		hook.AfterUpdate(md.ctx)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (md *DbModel) Update(args ...interface{}) (int64, error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return 0, md.Err
	}

	if len(args) > 1 {
		return 0, errors.New("arg number can not great 1")
	}

	if len(args) < 1 {
		return md.exec_update()
	}

	ent_ptr := args[0]
	if ent_ptr == nil {
		return 0, errors.New("ent_ptr must be *Struct")
	}
	//ent_ptr必须是一个指针
	v_ent := reflect.ValueOf(ent_ptr)
	if v_ent.Kind() != reflect.Ptr {
		return 0, errors.New("ent_ptr must be *Struct")
	}
	//ent_ptr必须是一个结构体指针
	v_ent = reflect.Indirect(v_ent)
	if v_ent.Kind() != reflect.Struct {
		return 0, errors.New("ent_ptr must be *Struct")
	}

	v_base := reflect.Indirect(reflect.ValueOf(md.ent_ptr))
	if v_ent.Type() == v_base.Type() {
		v_base.Set(v_ent)
	} else if ptr, ok := ent_ptr.(VoSaver); ok {
		ptr.SaveToModel(md, md.ent_ptr)
	} else {
		CopyDataToModel(md, ent_ptr, md.ent_ptr)
	}

	return md.exec_update()
}

func (md *DbModel) Delete() (int64, error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return 0, md.Err
	}

	if len(md.db_where.Tpl_sql) < 1 {
		return 0, errors.New("no where clause")
	}

	if hook, ok := md.ent_ptr.(BeforeDeleteInterface); ok {
		hook.BeforeDelete(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelDeleteSql(md)
	res, err := md.db_ptr.ExecSQL(&md.SqlContex, sql_str, md.db_where.Values...)
	if err != nil {
		return 0, err
	}

	if hook, ok := md.ent_ptr.(AfterDeleteInterface); ok {
		hook.AfterDelete(md.ctx)
	}

	num, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return num, nil
}

//id>0调用Update, 否则调用Insert
func (md *DbModel) UpdateOrInsert(id int64, args ...interface{}) (affected int64, entId int64, err error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return 0, id, md.Err
	}

	if id > 0 {
		num, err := md.ID(id).Update(args...)
		return num, id, err
	} else {
		id, err := md.Insert(args...)
		return 1, id, err
	}
}

//若存在记录，则调用Update，否则调用Insert
func (md *DbModel) Save(args ...interface{}) (affected int64, insertId int64, err error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return 0, 0, md.Err
	}

	has, err := md.Exist()
	if err != nil {
		return 0, 0, err
	}

	if has == true {
		num, err := md.Update(args...)
		return num, 0, err
	} else {
		id, err := md.Insert(args...)
		return 0, id, err
	}
}

//若不存在记录，则调用Insert
func (md *DbModel) InsertIfNotExist(args ...interface{}) (affected int64, insertId int64, err error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return 0, 0, md.Err
	}

	has, err := md.Exist()
	if err != nil {
		return 0, 0, err
	}

	if has == true {
		return 0, 0, nil
	}

	id, err := md.Insert(args...)
	return 0, id, err
}
