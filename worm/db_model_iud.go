package worm

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

func (md *DbModel)get_feild_flag_insert(i int) bool {
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

func (md *DbModel)gen_sql_insert() string {
	var buffer bytes.Buffer
	index := 0;
	buffer.WriteString(fmt.Sprintf("insert into %s (", md.table_name))
	for i, item := range md.flds_addr {
		if md.get_feild_flag_insert(i) == false {
			continue
		}
		if index > 0{
			buffer.WriteString(",")
		}
		buffer.WriteString(item.FName)
		index += 1
	}
	buffer.WriteString(")")

	index = 0;
	buffer.WriteString(" values (")
	for i, _ := range md.flds_addr {
		if md.get_feild_flag_insert(i) == false {
			continue
		}
		if index > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString("?")
		index += 1
	}
	buffer.WriteString(")")

	return buffer.String()
}

func (md *DbModel) get_fieldaddr_insert() []interface{} {
	cc := 0
	for i, _ := range md.flds_addr {
		if md.get_feild_flag_insert(i) == false {
			continue
		}
		cc+=1
	}

	index := 0
	vals:= make([]interface{}, cc)
	for i, _ := range md.flds_addr {
		if md.get_feild_flag_insert(i) == false {
			continue
		}
		vals[index] = md.flds_addr[i].VAddr
		index += 1
	}
	return vals
}

func (md *DbModel)exec_insert() (int64, error) {
	if md.Err != nil {
		return 0, md.Err
	}

	if hook, ok := md.ent_ptr.(BeforeInsertInterface); ok {
		hook.BeforeInsert(md.ctx)
	}

	sql_str := md.gen_sql_insert()
	values:= md.get_fieldaddr_insert()
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

func (md *DbModel)Insert(args ...interface{}) (int64, error) {
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

func (md *DbModel)get_feild_flag_update(i int) bool {
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

func (md *DbModel)gen_sql_update() string {
	var buffer bytes.Buffer
	buffer.WriteString("update ")
	buffer.WriteString(md.table_name)
	buffer.WriteString(" set ")
	index := 0;
	for i, item := range md.flds_addr {
		if md.get_feild_flag_update(i) == false {
			continue
		}
		if index > 0{
			buffer.WriteString(",")
		}
		buffer.WriteString(item.FName)
		buffer.WriteString("=?")
		index += 1
	}

	if len(md.db_where.Tpl_sql)>0 {
		buffer.WriteString(" where ")
		buffer.WriteString(md.db_where.Tpl_sql)
	}

	return buffer.String()
}

func (md *DbModel)get_fieldaddr_update() []interface{} {
	cc := 0
	for i, _ := range md.flds_addr {
		if md.get_feild_flag_update(i) == false {
			continue
		}
		cc+=1
	}

	index := 0
	vals:= make([]interface{}, cc)
	for i, _ := range md.flds_addr {
		if md.get_feild_flag_update(i) == false {
			continue
		}
		vals[index] = md.flds_addr[i].VAddr
		index += 1
	}
	return vals
}

func (md *DbModel)exec_update() (int64, error) {
	if md.Err != nil {
		return 0, md.Err
	}

	if len(md.db_where.Tpl_sql) < 1 {
		return  0, errors.New("no where clause")
	}

	if hook, ok := md.ent_ptr.(BeforeUpdateInterface); ok {
		hook.BeforeUpdate(md.ctx)
	}

	sql_str := md.gen_sql_update()
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
	if err != nil{
		return 0, err
	}
	return num, nil
}

func (md *DbModel)Update(args ...interface{}) (int64, error) {
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

func (md *DbModel)gen_delete() string {
	var buffer bytes.Buffer
	buffer.WriteString("delete from ")
	buffer.WriteString(md.table_name)
	if len(md.db_where.Tpl_sql)>0 {
		buffer.WriteString(" where ")
		buffer.WriteString(md.db_where.Tpl_sql)
	}
	return buffer.String()
}

func (md *DbModel)Delete() (int64, error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return 0, md.Err
	}

	if len(md.db_where.Tpl_sql) < 1 {
		return  0, errors.New("no where clause")
	}

	if hook, ok := md.ent_ptr.(BeforeDeleteInterface); ok {
		hook.BeforeDelete(md.ctx)
	}

	sql_str := md.gen_delete()
	res, err := md.db_ptr.ExecSQL(&md.SqlContex, sql_str, md.db_where.Values...)
	if err != nil {
		return 0, err
	}

	if hook, ok := md.ent_ptr.(AfterDeleteInterface); ok {
		hook.AfterDelete(md.ctx)
	}

	num, err := res.RowsAffected()
	if err != nil{
		return 0, err
	}
	return num, nil
}

//id>0调用Update, 否则调用Insert
func (md *DbModel)UpdateOrInsert(id int64, args ...interface{}) (affected int64, insertId int64, err error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return 0, 0, md.Err
	}

	if id > 0 {
		num, err := md.ID(id).Update(args...)
		return num, 0, err
	} else {
		id, err := md.Insert(args...)
		return 0, id, err
	}
}

//若存在记录，则调用Update，否则调用Insert
func (md *DbModel)Save(args ...interface{}) (affected int64, insertId int64, err error) {
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
func (md *DbModel)InsertIfNotExist(args ...interface{}) (affected int64, insertId int64, err error) {
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

