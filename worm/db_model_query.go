package worm

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

func (md *DbModel)GetFieldFlag4Select(i int) bool {
	if md.flds_addr[i].Flag == false {
		return false
	}
	if md.flds_info[i].NotSelect == true {
		return false
	}
	return true
}

func (md *DbModel)gen_select_fields() string {
	var buffer bytes.Buffer
	index := 0;
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
			buffer.WriteString(" as ")
			buffer.WriteString(md.table_alias)
			buffer.WriteString("_")
			buffer.WriteString(item.FName)
		} else {
			buffer.WriteString(item.FName)
		}
		index += 1
	}
	return buffer.String()
}

func (md *DbModel)gen_count_sql(count_field string) string {
	var buffer bytes.Buffer

	buffer.WriteString("select ")
	buffer.WriteString(count_field)
	buffer.WriteString(" from ")
	buffer.WriteString(md.table_name)
	if len(md.table_alias) > 0 {
		buffer.WriteString(" ")
		buffer.WriteString(md.table_alias)
	}

	if len(md.db_where.Tpl_sql)>0 {
		buffer.WriteString(" where ")
		buffer.WriteString(md.db_where.Tpl_sql)
	}

	return buffer.String()
}

func (md *DbModel)get_scan_valus() []interface{} {
	index := 0
	vals:= make([]interface{}, len(md.flds_addr))
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

func (md *DbModel)Scan() (bool, error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return false, md.Err
	}

	if hook, isHook := md.ent_ptr.(BeforeQueryInterface); isHook {
		hook.BeforeQuery(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelGet(md)
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

	scan_vals:= md.get_scan_valus()
	err = rows.Scan(scan_vals...)
	if err != nil {
		rows.Close()
		return false, err
	}

	rows.Close()
	return true, nil
}

func (md *DbModel)Get(args ...interface{}) (bool, error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
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

	//若数组成员与model不同，则调用相应的函数选择对应的字段
	//若数组成员是vo，则调用getSelectFieldsByVo选择对应的字段
	//若数组成员是struct，则调用getSelectFieldsByEo选择对应的字段
	var mo_ptr interface{} = nil
	var vo_ptr VoLoader = nil
	var eo_ptr interface{} = nil
	v_base := reflect.Indirect(reflect.ValueOf(md.ent_ptr))
	if v_ent.Type() == v_base.Type() {
		mo_ptr = ent_ptr
	} else if ptr, ok := ent_ptr.(VoLoader); ok {
		vo_ptr = ptr
		getSelectFieldsByVo(md, vo_ptr)
	} else {
		eo_ptr = ent_ptr
		getSelectFieldsByEo(md, ent_ptr)
	}

	has, err := md.Scan()
	if err != nil {
		return has, err
	}

	//若数组成员与model不同，则调用相应的函数给ent_ptr赋值
	//若数组成员是vo，则调用LoadFromModel给ent_ptr赋值
	//若数组成员是struct，则调用CopyDataFromModel给ent_ptr赋值
	if mo_ptr != nil {
		v_ent.Set(v_base)
	} else if vo_ptr != nil {
		vo_ptr.LoadFromModel(nil, md.ent_ptr)
	} else {
		CopyDataFromModel(nil, eo_ptr, md.ent_ptr)
	}

	return true, nil
}

func (md *DbModel)Exist() (bool, error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return false, md.Err
	}

	if hook, isHook := md.ent_ptr.(BeforeQueryInterface); isHook {
		hook.BeforeQuery(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelGet(md)
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

func (md *DbModel)Count(field ...string) (int64, error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
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
		sub_sql := md.db_ptr.engine.db_dialect.GenModelFind(md)
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
	err = Scan(rows, &total)
	rows.Close()

	return total, nil
}

func (md *DbModel)DistinctCount(field string) (int64, error) {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return 0, md.Err
	}

	if hook, isHook := md.ent_ptr.(BeforeQueryInterface); isHook {
		hook.BeforeQuery(md.ctx)
	}

	count_field := fmt.Sprintf("count(distinct %s)", field)
	sql_str := md.gen_count_sql(count_field)
	if len(md.group_by) > 0 {
		sub_sql := md.db_ptr.engine.db_dialect.GenModelFind(md)
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
	err = Scan(rows, &total)
	rows.Close()

	return total, nil
}

func (md *DbModel)Rows() (*ModelRows, error) {
	if md.Err != nil {
		return nil, md.Err
	}

	if hook, isHook := md.ent_ptr.(BeforeQueryInterface); isHook {
		hook.BeforeQuery(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelFind(md)
	rows, err := md.db_ptr.ExecQuery(&md.SqlContex, sql_str, md.db_where.Values...)
	if err != nil {
		return nil, err
	}

	rs := new(ModelRows)
	rs.rows = rows
	rs.model = md
	return rs, nil
}

func (md *DbModel)Find(arr_ptr interface{}) error {
	if md.auto_put && md.md_pool != nil {
		pool := md.split_pool()
		defer md.put_pool(pool)
	}
	if md.Err != nil {
		return md.Err
	}

	//获取变量arr_ptr的类型
	v_arr := reflect.ValueOf(arr_ptr)
	if v_arr.Kind() != reflect.Ptr {
		return  errors.New("arr_ptr must be *Slice")
	}
	//获取变量arr_ptr指向的类型
	t_arr := GetDirectType(v_arr.Type())
	if t_arr.Kind() != reflect.Slice {
		return  errors.New("arr_ptr must be *Slice")
	}
	//获取数组成员的类型
	t_item := GetDirectType(t_arr.Elem())
	if t_item.Kind() != reflect.Struct {
		return errors.New("arr_ptr must be *Struct")
	}

	//若数组成员与model相同，则v_item_base指向md.ent_ptr
	//若数组成员与model不同，则创建一个成员对象v_item, 并使v_item_base指向v_item
	var item_vo VoLoader = nil
	var item_ptr interface{} = nil
	v_item_base := reflect.Indirect(reflect.ValueOf(md.ent_ptr))
	if t_item != GetDirectType(reflect.TypeOf(md.ent_ptr)) {
		v_item := reflect.New(t_item)
		item_ptr = v_item.Interface()
		vo_ptr, ok := item_ptr.(VoLoader)
		if ok {
			item_vo = vo_ptr
		}
		v_item_base = reflect.Indirect(v_item)
	}

	//若数组成员与model不同，则调用相应的函数选择对应的字段
	//若数组成员是vo，则调用getSelectFieldsByVo选择对应的字段
	//若数组成员是struct，则调用getSelectFieldsByEo选择对应的字段
	if item_vo != nil {
		getSelectFieldsByVo(md, item_vo)
	} else if item_ptr != nil {
		getSelectFieldsByEo(md, item_ptr)
	}

	if hook, isHook := md.ent_ptr.(BeforeQueryInterface); isHook {
		hook.BeforeQuery(md.ctx)
	}

	sql_str := md.db_ptr.engine.db_dialect.GenModelFind(md)
	rows, err := md.db_ptr.ExecQuery(&md.SqlContex, sql_str, md.db_where.Values...)
	if err != nil {
		return err
	}

	vals:= md.get_scan_valus()
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

		//若数组成员与model不同，则调用相应的函数给v_item赋值
		//若数组成员是vo，则调用LoadFromModel给v_item赋值
		//若数组成员是struct，则调用CopyDataFromModel给v_item赋值
		if item_vo != nil {
			item_vo.LoadFromModel(nil, md.ent_ptr)
		} else if item_ptr != nil {
			CopyDataFromModel(nil, item_ptr, md.ent_ptr)
		}
		v_arr_base.Set(reflect.Append(v_arr_base, v_item_base))
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
