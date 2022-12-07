package worm

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

//自动模式追加一批字段(通过地址字段匹配)
//添加字段前需要清空当前状态（缺省选择全部）
func (md *DbModel) auto_add_field_addr(fields ...interface{}) *DbModel {
	if md.flag_auto == false {
		md.OmitALL()
		md.flag_auto = true
	}

	for _, fld_ptr := range fields {
		if fld_ptr == nil {
			md.Err = errors.New("field addr is nil")
			return md
		}
		ret := md.set_flag_by_addr(fld_ptr, true)
		if ret == false {
			md.Err = errors.New("field not found")
			return md
		}
	}

	return md
}

//自动模式追加一批字段(通过index匹配)
//添加字段前需要清空当前状态（缺省选择全部）
func (md *DbModel) auto_add_field_index(fields ...int) *DbModel {
	if md.flag_auto == false {
		md.OmitALL()
		md.flag_auto = true
	}

	for _, index := range fields {
		ret := md.set_flag_by_index(index, true)
		if ret == false {
			md.Err = errors.New("field not found")
			return md
		}
	}

	return md
}

func (md *DbModel) auto_add_field_all() *DbModel {
	if md.flag_auto == false {
		md.flag_auto = true
	}
	num := len(md.flds_addr)
	for i := 0; i < num; i++ {
		md.flds_addr[i].Flag = true
	}
	return md
}

//执行vo=mo的赋值操作，只有名称相同、类型相同的字段才能赋值
//若md != nil，则获取mo的字段地址，并调用md的set_flag_by_addr函数来选中该字段
//只有被选中的字段才需要从数据库中查询
func CopyDataFromModel(md *DbModel, vo_ptr interface{}, mo_ptr interface{}) (int, error) {
	if vo_ptr == nil {
		return 0, errors.New("vo_ptr is nil")
	}
	v_vo := reflect.ValueOf(vo_ptr)
	if v_vo.Kind() != reflect.Ptr {
		return 0, errors.New("vo_ptr must be Pointer")
	}
	v_vo = reflect.Indirect(v_vo)
	if v_vo.Kind() != reflect.Struct {
		return 0, errors.New("vo_ptr  muse be Struct")
	}

	if mo_ptr == nil {
		return 0, errors.New("mo_ptr is nil")
	}
	v_mo := reflect.ValueOf(mo_ptr)
	if v_mo.Kind() != reflect.Ptr {
		return 0, errors.New("mo_ptr must be Pointer")
	}
	v_mo = reflect.Indirect(v_mo)
	if v_mo.Kind() != reflect.Struct {
		return 0, errors.New("mo_ptr  muse be Struct")
	}

	//获取字段交集
	pflds, err := getPubField4VoMo("", v_vo.Type(), v_mo.Type())
	if err != nil {
		return 0, err
	}

	//将公共字段添加到选扩展择集中
	//若存在model字段，不用设置扩展选择集
	if md != nil {
		if pflds.ModelField < 0 {
			for _, item := range pflds.Fields {
				md.auto_add_field_index(item.MoIndex)
			}
		}
	}

	//若vo中存在Model字段，只需要赋值Model对应的字段即可
	if pflds.ModelField >= 0 {
		fv_vo := v_vo.Field(pflds.ModelField)
		if fv_vo.CanSet() == true {
			fv_vo.Set(v_mo)
			return 1, nil
		}
	}

	count := 0
	for _, item := range pflds.Fields {
		fv_vo := v_vo.FieldByIndex(item.VoIndex)
		fv_mo := v_mo.Field(item.MoIndex)
		if fv_vo.CanSet() == false {
			continue
		}
		fv_vo.Set(fv_mo)
		count += 1
	}
	return count, nil
}

//执行mo=vo的赋值操作，只有名称相同、类型相同的字段才能赋值
//若md != nil，则获取mo的字段地址，并调用md的set_flag_by_addr函数来选中该字段
//只有被选中的字段才能更新到数据库中
func CopyDataToModel(md *DbModel, vo_ptr interface{}, mo_ptr interface{}) (int, error) {
	if vo_ptr == nil {
		return 0, errors.New("vo_ptr is nil")
	}
	v_vo := reflect.ValueOf(vo_ptr)
	if v_vo.Kind() != reflect.Ptr {
		return 0, errors.New("vo_ptr must be Pointer")
	}
	v_vo = reflect.Indirect(v_vo)
	if v_vo.Kind() != reflect.Struct {
		return 0, errors.New("vo_ptr  muse be Struct")
	}

	if mo_ptr == nil {
		return 0, errors.New("mo_ptr is nil")
	}
	v_mo := reflect.ValueOf(mo_ptr)
	if v_mo.Kind() != reflect.Ptr {
		return 0, errors.New("mo_ptr must be Pointer")
	}
	v_mo = reflect.Indirect(v_mo)
	if v_mo.Kind() != reflect.Struct {
		return 0, errors.New("mo_ptr  muse be Struct")
	}

	//获取字段交集
	pflds, err := getPubField4VoMo("", v_vo.Type(), v_mo.Type())
	if err != nil {
		return 0, err
	}

	//将公共字段添加到选扩展择集中
	//若存在model字段，不用设置扩展选择集
	if md != nil {
		if pflds.ModelField < 0 {
			for _, item := range pflds.Fields {
				md.auto_add_field_index(item.MoIndex)
			}
		}
	}

	//若vo中存在Model字段，只需要赋值Model对应的字段即可
	if pflds.ModelField >= 0 {
		fv_vo := v_vo.Field(pflds.ModelField)
		if v_mo.CanSet() == true {
			v_mo.Set(fv_vo)
			return 1, nil
		}
	}

	count := 0
	for _, item := range pflds.Fields {
		fv_vo := v_vo.FieldByIndex(item.VoIndex)
		fv_mo := v_mo.Field(item.MoIndex)
		if fv_mo.CanSet() == false {
			continue
		}
		fv_mo.Set(fv_vo)
		count += 1
	}
	return count, nil
}

func set_value(fld_ptr interface{}, val interface{}) error {
	if fld_ptr == nil {
		return errors.New("fld_ptr must be Pointer")
	}
	v_fld := reflect.ValueOf(fld_ptr)
	if v_fld.Kind() != reflect.Ptr {
		return errors.New("fld_ptr must be Pointer")
	}

	v_fld = reflect.Indirect(v_fld)
	v_val := reflect.ValueOf(val)
	t_fld := v_fld.Type()
	t_val := v_val.Type()
	k_fld := t_fld.Kind()
	k_val := t_val.Kind()
	if k_fld == k_val {
		v_fld.Set(v_val)
		return nil
	}

	//类型不相同，但可以进行转换，转换后赋值
	if k_val == reflect.String {
		switch k_fld {
		case reflect.Int, reflect.Int32, reflect.Int64:
			tmp, err := strconv.ParseInt(val.(string), 10, 64)
			if err != nil {
				return err
			}
			v_fld.SetInt(tmp)
			return nil
		case reflect.Float32, reflect.Float64:
			tmp, err := strconv.ParseFloat(val.(string), 64)
			if err != nil {
				return err
			}
			v_fld.SetFloat(tmp)
			return nil
		}
	} else if k_fld == reflect.Int64 {
		switch k_val {
		case reflect.Int, reflect.Int32:
			v_fld.SetInt(v_val.Int())
			return nil
		}
	} else if k_fld == reflect.Int32 || k_fld == reflect.Int {
		switch k_val {
		case reflect.Int, reflect.Int32:
			v_fld.SetInt(v_val.Int())
			return nil
		}
	} else if k_fld == reflect.Float64 {
		switch k_val {
		case reflect.Float32:
			v_fld.SetFloat(v_val.Float())
			return nil
		}
	}

	return errors.New(fmt.Sprintf("incorrect data type: %v != %v", t_fld, t_val))
}

//设置字段的值，并选中该字段
func (md *DbModel) SetValue(fld_ptr interface{}, val interface{}) error {
	err := set_value(fld_ptr, val)
	if err != nil {
		md.Err = err
		return err
	}
	md.auto_add_field_addr(fld_ptr)
	return nil
}

func SetValue(md *DbModel, fld_ptr interface{}, val interface{}) error {
	if md != nil {
		return md.SetValue(fld_ptr, val)
	} else {
		return set_value(fld_ptr, val)
	}
}

func GetPointer(md *DbModel, fld_ptr interface{}) interface{} {
	if md == nil {
		return fld_ptr
	} else {
		md.auto_add_field_addr(fld_ptr)
		return fld_ptr
	}
}

func GetIndirect(md *DbModel, fld_ptr interface{}) interface{} {
	if fld_ptr == nil {
		return nil
	}

	v_fld := reflect.Indirect(reflect.ValueOf(fld_ptr))
	if md == nil {
		return v_fld.Interface()
	} else {
		md.auto_add_field_addr(fld_ptr)
		return v_fld.Interface()
	}
}

func GetBool(md *DbModel, fld_ptr *bool) bool {
	if md == nil {
		return *fld_ptr
	} else {
		md.auto_add_field_addr(fld_ptr)
		return *fld_ptr
	}
}

func GetInt(md *DbModel, fld_ptr *int) int {
	if md == nil {
		return *fld_ptr
	} else {
		md.auto_add_field_addr(fld_ptr)
		return *fld_ptr
	}
}

func GetInt32(md *DbModel, fld_ptr *int32) int32 {
	if md == nil {
		return *fld_ptr
	} else {
		md.auto_add_field_addr(fld_ptr)
		return *fld_ptr
	}
}

func GetInt64(md *DbModel, fld_ptr *int64) int64 {
	if md == nil {
		return *fld_ptr
	} else {
		md.auto_add_field_addr(fld_ptr)
		return *fld_ptr
	}
}

func GetFloat32(md *DbModel, fld_ptr *float32) float32 {
	if md == nil {
		return *fld_ptr
	} else {
		md.auto_add_field_addr(fld_ptr)
		return *fld_ptr
	}
}

func GetFloat64(md *DbModel, fld_ptr *float64) float64 {
	if md == nil {
		return *fld_ptr
	} else {
		md.auto_add_field_addr(fld_ptr)
		return *fld_ptr
	}
}

func GetString(md *DbModel, fld_ptr *string) string {
	if md == nil {
		return *fld_ptr
	} else {
		md.auto_add_field_addr(fld_ptr)
		return *fld_ptr
	}
}

func GetTime(md *DbModel, fld_ptr *time.Time) time.Time {
	if md == nil {
		return *fld_ptr
	} else {
		md.auto_add_field_addr(fld_ptr)
		return *fld_ptr
	}
}
