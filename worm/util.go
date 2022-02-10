package worm

import (
	"errors"
	"reflect"
)

//获取直接类型
func GetDirectType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

//从对象数组中获取整形字段的值
func CreateIntArray(arr_ptr interface{}, field string, def_val ...int64) ([]int64, error) {
	v_arr := reflect.Indirect(reflect.ValueOf(arr_ptr))
	if v_arr.Kind() != reflect.Slice {
		return nil, errors.New("arr_ptr must be *Slice")
	}
	t_item := v_arr.Type().Elem()
	if t_item.Kind() != reflect.Struct {
		return nil, errors.New("array item type muse be reflect.Struct")
	}
	fi_info, ok := t_item.FieldByName(field)
	if !ok {
		return nil, errors.New("not find field in struct")
	}

	var arr []int64
	num := v_arr.Len()
	for i:=0; i < num; i++ {
		v_item := v_arr.Index(i)
		v_ff := v_item.FieldByIndex(fi_info.Index)
		arr = append(arr, v_ff.Int())
	}

	if len(arr) < 1 && len(def_val) > 0 {
		arr = append(arr, def_val...)
	}

	return arr, nil
}

