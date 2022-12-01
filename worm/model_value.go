package worm

import "reflect"

//获取stuct对象的字段地址数组
func getEntFieldAddrs(fields []FieldInfo, v_ent reflect.Value, flag bool) []FieldValue {
	f_num := len(fields)
	arr := make([]FieldValue, f_num)
	v_ent = reflect.Indirect(v_ent)
	for i := 0; i < f_num; i++ {
		f_name := fields[i].DbName
		vv := v_ent.Field(fields[i].FieldIndex)
		v_ptr := vv.Addr().Interface()
		/*
			var v_ptr interface{}
			if fields[i].FieldPos == nil {
				vv := v_ent.Field(fields[i].FieldIndex)
				v_ptr = vv.Addr().Interface()
			} else {
				vv := v_ent.FieldByIndex(fields[i].FieldPos)
				v_ptr = vv.Addr().Interface()
			}*/
		var item FieldValue
		item.FName = f_name
		item.VAddr = v_ptr
		item.Flag = flag
		arr[i] = item
	}
	return arr
}

//将values的地址替换为当前对象的字段地址
func rebindEntAddrs(fields []FieldInfo, v_ent reflect.Value, values []FieldValue) {
	f_num := len(fields)
	v_ent = reflect.Indirect(v_ent)
	for i := 0; i < f_num; i++ {
		vv := v_ent.Field(fields[i].FieldIndex)
		v_ptr := vv.Addr().Interface()
		/*
			var v_ptr interface{}
			if fields[i].FieldPos == nil {
				vv := v_ent.Field(fields[i].FieldIndex)
				v_ptr = vv.Addr().Interface()
			} else {
				vv := v_ent.FieldByIndex(fields[i].FieldPos)
				v_ptr = vv.Addr().Interface()
			}*/
		values[i].VAddr = v_ptr
	}
}

//为数据库提供scan要求的变量地址(用于rows.Scan)
//columns：数据库字段
//返回：DbField以及new(interface{})组成的数组
//说明：若没有数据库col对应的字段，则用new(interface{})代替
func genScanAddr4Columns(columns []string, v_ent reflect.Value) []interface{} {
	minfo := getModelInfoUseCache(v_ent)
	ent_flds := getEntFieldAddrs(minfo.Fields, v_ent, true)
	values := make([]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		var ptr interface{} = nil
		for j := 0; j < len(ent_flds); j++ {
			if ent_flds[j].FName == columns[i] {
				//ptr = ent_flds[j].VAddr
				ptr = &ent_flds[j]
				break
			}
		}
		if ptr == nil {
			ptr = new(interface{})
		}
		values[i] = ptr
	}
	return values
}
