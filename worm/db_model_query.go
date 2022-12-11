package worm

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

func (md *DbModel) GetFieldFlag4Select(i int) bool {
	//没有被人工选择
	if md.flds_addr[i].Flag == false {
		return false
	}
	//该字段不用于Select
	if md.flds_info[i].NotSelect == true {
		return false
	}
	return true
}

func (md *DbModel) gen_select_fields() string {
	var buffer bytes.Buffer
	index := 0
	for i, item := range md.flds_addr {
		if md.GetFieldFlag4Select(i) == false {
			continue
		}
		if index > 0 {
			buffer.WriteString(",")
		}
		if len(md.table_alias) > 0 {
			buffer.WriteString(md.table_alias)
			buffer.WriteString(".")
			buffer.WriteString(item.FName)
			//buffer.WriteString(" as ")
			//buffer.WriteString(md.table_alias)
			//buffer.WriteString("_")
			//buffer.WriteString(item.FName)
		} else {
			buffer.WriteString(item.FName)
		}
		index += 1
	}
	return buffer.String()
}

func (md *DbModel) gen_count_sql(count_field string) string {
	var buffer bytes.Buffer

	buffer.WriteString("select ")
	buffer.WriteString(count_field)
	buffer.WriteString(" from ")
	buffer.WriteString(md.table_name)
	if len(md.table_alias) > 0 {
		buffer.WriteString(" ")
		buffer.WriteString(md.table_alias)
	}

	if len(md.db_where.Tpl_sql) > 0 {
		buffer.WriteString(" where ")
		buffer.WriteString(md.db_where.Tpl_sql)
	}

	return buffer.String()
}

func (md *DbModel) get_scan_valus() []interface{} {
	index := 0
	vals := make([]interface{}, len(md.flds_addr))
	for i, _ := range md.flds_addr {
		if md.GetFieldFlag4Select(i) == false {
			continue
		}
		//vals[index] = md.flds_addr[i].VAddr
		vals[index] = &md.flds_addr[i]
		index += 1
	}
	return vals[0:index]
}

func (md *DbModel) Scan() (bool, error) {
	if md.Err != nil {
		return false, md.Err
	}

	if hook, isHook := md.ent_ptr.(BeforeQueryInterface); isHook {
		hook.BeforeQuery(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelGetSql(md)
	rows, err := md.db_ptr.ExecQuery(&md.SqlContex, sql_str, md.db_where.Values...)
	if err != nil {
		return false, err
	}

	if hook, isAfter := md.ent_ptr.(AfterQueryInterface); isAfter {
		hook.AfterQuery(md.ctx)
	}

	if !rows.Next() {
		rows.Close()
		return false, nil
	}

	scan_vals := md.get_scan_valus()
	err = rows.Scan(scan_vals...)
	if err != nil {
		rows.Close()
		return false, err
	}

	rows.Close()
	return true, nil
}

//若flag==true，则正常执行
//若flag==false，则退出执行
func (md *DbModel) GetIf(flag bool, args ...interface{}) (bool, error) {
	if flag == false {
		return false, nil
	}
	return md.Get(args...)
}

func (md *DbModel) Get(args ...interface{}) (bool, error) {
	if md.Err != nil {
		return false, md.Err
	}

	if len(args) > 1 {
		return false, errors.New("arg number can not great 1")
	}

	//参数为空, 调用Scan()
	if len(args) < 1 {
		return md.Scan()
	}

	ent_ptr := args[0]
	if ent_ptr == nil {
		return false, errors.New("ent_ptr must be *Struct")
	}
	//ent_ptr必须是一个指针
	v_ent := reflect.ValueOf(ent_ptr)
	if v_ent.Kind() != reflect.Ptr {
		return false, errors.New("ent_ptr must be *Struct")
	}
	//ent_ptr必须是一个结构体指针
	v_ent = reflect.Indirect(v_ent)
	if v_ent.Kind() != reflect.Struct {
		return false, errors.New("ent_ptr must be *Struct")
	}

	//若args是model类型，则将mo_ptr指向ent_ptr
	//若args是vo，则调用selectFieldsByVo选择对应的字段，并将vo_ptr指向ent_ptr
	//若args是eo，则调用selectFieldsByEo选择对应的字段
	var mo_ptr interface{} = nil
	var vo_ptr VoLoader = nil
	var pflds *PublicFields = nil
	t_ent := v_ent.Type()
	if t_ent == md.ent_type {
		mo_ptr = ent_ptr
	} else if ptr, ok := ent_ptr.(VoLoader); ok {
		vo_ptr = ptr
		selectFieldsByVo(md, vo_ptr)
	} else {
		pflds = md.selectFieldsByEo(t_ent)
	}

	has, err := md.Scan()
	if err != nil {
		return has, err
	}

	//mo_ptr != nil，说明ent_ptr是一个Model类型, 则调用Value.Set给ent_ptr赋值
	//vo_ptr != nil，说明ent_ptr是一个vo，则调用LoadFromModel给ent_ptr赋值
	//mo_ptr == nil && vo_ptr == nil，说明ent_ptr是一个eo，则调用CopyDataFromModel给ent_ptr赋值
	if mo_ptr != nil {
		v_ent.Set(md.ent_value)
	} else if vo_ptr != nil {
		vo_ptr.LoadFromModel(nil, md.ent_ptr)
	} else {
		md.copyModelData2Eo(pflds, v_ent)
	}
	return true, nil
}

func (md *DbModel) Exist() (bool, error) {
	if md.Err != nil {
		return false, md.Err
	}

	if hook, isHook := md.ent_ptr.(BeforeQueryInterface); isHook {
		hook.BeforeQuery(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelGetSql(md)
	rows, err := md.db_ptr.ExecQuery(&md.SqlContex, sql_str, md.db_where.Values...)
	if err != nil {
		return false, err
	}

	if !rows.Next() {
		rows.Close()
		return false, nil
	}

	rows.Close()
	return true, nil
}

func (md *DbModel) Count(field ...string) (int64, error) {
	if md.Err != nil {
		return 0, md.Err
	}
	if len(field) > 1 {
		return 0, errors.New("field vumber > 0")
	}

	if hook, isHook := md.ent_ptr.(BeforeQueryInterface); isHook {
		hook.BeforeQuery(md.ctx)
	}

	count_field := "count(1)"
	if len(field) == 1 {
		count_field = fmt.Sprintf("count(%s)", field[0])
	}

	sql_str := md.gen_count_sql(count_field)
	if len(md.group_by) > 0 {
		sub_sql := md.db_ptr.engine.db_dialect.GenModelFindSql(md)
		sql_str = fmt.Sprintf("select %s from (%s) tmp", count_field, sub_sql)
	}

	rows, err := md.db_ptr.ExecQuery(&md.SqlContex, sql_str, md.db_where.Values...)
	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		rows.Close()
		return 0, nil
	}

	var total int64
	err = rows_scan(rows, &total)
	rows.Close()

	return total, nil
}

func (md *DbModel) DistinctCount(field string) (int64, error) {
	if md.Err != nil {
		return 0, md.Err
	}

	if hook, isHook := md.ent_ptr.(BeforeQueryInterface); isHook {
		hook.BeforeQuery(md.ctx)
	}

	count_field := fmt.Sprintf("count(distinct %s)", field)
	sql_str := md.gen_count_sql(count_field)
	if len(md.group_by) > 0 {
		sub_sql := md.db_ptr.engine.db_dialect.GenModelFindSql(md)
		sql_str = fmt.Sprintf("select %s from (%s) tmp", count_field, sub_sql)
	}

	rows, err := md.db_ptr.ExecQuery(&md.SqlContex, sql_str, md.db_where.Values...)
	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		rows.Close()
		return 0, nil
	}

	var total int64
	err = rows_scan(rows, &total)
	rows.Close()

	return total, nil
}

func (md *DbModel) Rows() (ModelRows, error) {
	var rs ModelRows
	rs.model = md

	if md.Err != nil {
		return rs, md.Err
	}

	if hook, isHook := md.ent_ptr.(BeforeQueryInterface); isHook {
		hook.BeforeQuery(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelFindSql(md)
	rows, err := md.db_ptr.ExecQuery(&md.SqlContex, sql_str, md.db_where.Values...)
	if err != nil {
		return rs, err
	}

	rs.Rows = rows
	return rs, nil
}

func (md *DbModel) Find(arr_ptr interface{}) error {
	if md.Err != nil {
		return md.Err
	}

	//获取变量arr_ptr的类型
	v_arr := reflect.ValueOf(arr_ptr)
	if v_arr.Kind() != reflect.Ptr {
		return errors.New("arr_ptr must be *Slice")
	}
	//获取变量arr_ptr指向的类型
	t_arr := GetDirectType(v_arr.Type())
	if t_arr.Kind() != reflect.Slice {
		return errors.New("arr_ptr must be *Slice")
	}
	//获取数组成员的类型
	t_item := GetDirectType(t_arr.Elem())
	if t_item.Kind() != reflect.Struct {
		return errors.New("arr_ptr must be *Struct")
	}

	//若数组成员与model相同，则v_item指向md.ent_value（不用创建item对象，直接使用md.ent_value）
	//若数组成员与model不同，则创建一个成员对象v_item_ptr, 并使v_item指向v_item_ptr
	//将item_eo指向v_item_ptr, 若v_item_ptr是一个vo，则item_vo指向v_item_ptr
	var item_vo VoLoader = nil
	var item_eo interface{} = nil
	v_item := md.ent_value
	if t_item != md.ent_type {
		v_item_ptr := reflect.New(t_item)
		item_eo = v_item_ptr.Interface()
		v_item = v_item_ptr.Elem()
		if vo_ptr, ok := item_eo.(VoLoader); ok {
			item_vo = vo_ptr
		}
	}

	//若数组成员是vo，则调用selectFieldsByVo选择对应的字段
	//若数组成员是eo，则调用selectFieldsByEo选择对应的字段
	var pflds *PublicFields = nil
	if item_vo != nil {
		selectFieldsByVo(md, item_vo)
	} else if item_eo != nil {
		pflds = md.selectFieldsByEo(t_item)
	}

	if hook, isHook := md.ent_ptr.(BeforeQueryInterface); isHook {
		hook.BeforeQuery(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelFindSql(md)
	rows, err := md.db_ptr.ExecQuery(&md.SqlContex, sql_str, md.db_where.Values...)
	if err != nil {
		return err
	}

	vals := md.get_scan_valus()
	v_arr_base := reflect.Indirect(v_arr)
	hook, isAfter := md.ent_ptr.(AfterQueryInterface)
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			rows.Close()
			return err
		}
		if isAfter {
			hook.AfterQuery(md.ctx)
		}

		//若数组成员与model相同，不用进行赋值操作，直接将v_item添加到数组
		//若数组成员与model不同，则调用相应的函数给v_item赋值，然后将v_item添加到数组
		//若数组成员是vo，则调用LoadFromModel给v_item赋值
		//若数组成员是eo，则调用CopyModelData2Eo给v_item赋值
		if item_vo != nil {
			item_vo.LoadFromModel(nil, md.ent_ptr)
		} else if item_eo != nil {
			//md.CopyModelData2Eo(v_item)
			md.copyModelData2Eo(pflds, v_item)
		}
		v_arr_base.Set(reflect.Append(v_arr_base, v_item))
	}

	rows.Close()
	return nil
}

/*
//返回动态数组
func (md *DbModel)FindArray() (interface{}, error) {
	if md.Err != nil {
		return nil, md.Err
	}

	//数组成员的类型
	v_ent := reflect.ValueOf(md.ent_ptr)
	v_ent = reflect.Indirect(v_ent)
	t_ent := v_ent.Type()

	//动态创建数组
	t_arr := reflect.SliceOf(t_ent)
	v_arr_tmp := reflect.MakeSlice(t_arr, 0, 10)

	//创建一个变量指向数组
	v_arr_var := reflect.New(v_arr_tmp.Type()).Elem()
	v_arr_var.Set(v_arr_tmp)

	sql_str := md.gen_select()
	rows, err := md.db_ptr.ExecQuery(md.GetContext(), sql_str, md.db_where.Values...)
	if err != nil {
		return nil, md.Err
	}

	v_arr_base := v_arr_var
	v_item_base := v_ent
	vals:= md.get_scan_valus()
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			rows.Close()
			return nil, md.Err
		}
		v_arr_base.Set(reflect.Append(v_arr_base, v_item_base))
	}

	rows.Close()
	return v_arr_var.Interface(), nil
}
*/
