package worm

import "reflect"

//获取stuct对象的字段地址数组
func getEntFieldAddrs(fields []FieldInfo, v_ent reflect.Value, flag bool) []FieldValue {
	f_num := len(fields)
	arr := make([]FieldValue, f_num)
	v_ent = reflect.Indirect(v_ent)
	for i := 0; i < f_num; i++ {
		var item FieldValue
		item.VAddr = nil
		if fields[i].FieldIndex < 0 {
			arr[i] = item
			continue
		}

		f_name := fields[i].DbName
		vv := v_ent.Field(fields[i].FieldIndex)
		v_ptr := vv.Addr().Interface()

		item.FName = f_name
		item.VAddr = v_ptr
		item.Flag = flag
		arr[i] = item
	}
	return arr
}

//为数据库提供scan要求的变量地址(用于rows.Scan)
//columns：数据库字段
//返回：DbField以及new(interface{})组成的数组
//说明：若没有数据库col对应的字段，则用new(interface{})代替
func genScanAddr4Columns(columns []string, v_ent reflect.Value) []interface{} {
	minfo := getModelInfo(v_ent.Type())
	values := make([]interface{}, len(columns))
	for i := 0; i < len(columns); i++ {
		var ptr interface{} = nil
		moindex := minfo.get_field_index_dbname(columns[i])
		if moindex >= 0 {
			vv := v_ent.Field(moindex)
			var item FieldValue
			item.FName = columns[i]
			item.VAddr = vv.Addr().Interface()
			ptr = &item
		}
		if ptr == nil {
			ptr = new(interface{})
		}
		values[i] = ptr
	}
	return values
}
