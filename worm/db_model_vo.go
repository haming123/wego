package worm

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

type VoSaver interface {
	SaveToModel(md *DbModel, mo interface{})
}

type VoLoader interface {
	LoadFromModel(md *DbModel, mo interface{})
}

type VoInfo struct {
	VoName  string
	VoIndex []int
	MoIndex []int
	ValFlag bool
}
type VoFields []VoInfo

//vo、mo字段交集信息缓存
var g_voCopy_cache map[string]VoFields = make(map[string]VoFields)
var g_voCopy_mutex sync.Mutex

//与vo对应的mo的model的字段选中状态缓存
var g_voLoad_cache map[string][]FieldValue = make(map[string][]FieldValue)
var g_voLoad_mutex sync.Mutex

//首先从缓存中获取字段交集
//若缓存中不存在，则生成字段交集
func getPubField4VoMo(vo_ptr interface{}, mo_ptr interface{}) (VoFields, error) {
	g_voCopy_mutex.Lock()
	defer g_voCopy_mutex.Unlock()

	var arr VoFields = nil
	if vo_ptr == nil {
		return arr, errors.New("vo_ptr is nil")
	}
	t_vo := reflect.TypeOf(vo_ptr)
	if t_vo.Kind() != reflect.Ptr {
		return arr, errors.New("vo_ptr must be Pointer")
	}
	t_vo = GetDirectType(t_vo)
	if t_vo.Kind() != reflect.Struct {
		return arr, errors.New("vo_ptr  muse be Struct")
	}

	if mo_ptr == nil {
		return arr, errors.New("mo_ptr is nil")
	}
	t_mo := reflect.TypeOf(mo_ptr)
	if t_mo.Kind() != reflect.Ptr {
		return arr, errors.New("mo_ptr must be Pointer")
	}
	t_mo = GetDirectType(t_mo)
	if t_mo.Kind() != reflect.Struct {
		return arr, errors.New("mo_ptr  muse be Struct")
	}

	key := t_vo.String() + t_mo.String()
	arr, ok := g_voCopy_cache[key]
	if ok {
		return arr, nil
	}

	arr = make(VoFields, t_mo.NumField())
	arr = arr[:0]
	arr = getPubField4VoMoFunc(arr, t_vo, t_mo)
	g_voCopy_cache[key] = arr
	//fmt.Println("GenVoInfoJoinMo")

	return arr, nil
}

//生成vo与mo的字段交集信息
//只有名称与类型相同的字段才属于字段交集
func getPubField4VoMoFunc(arr VoFields, t_vo reflect.Type, t_mo reflect.Type) VoFields {
	f_num := t_mo.NumField()
	for i := 0; i < f_num; i++ {
		ft_mo := t_mo.Field(i)
		ft_vo, ok := t_vo.FieldByName(ft_mo.Name)
		if !ok {
			continue
		}
		if ft_vo.Type != ft_mo.Type {
			continue
		}

		var item VoInfo
		item.VoName = ft_vo.Name
		item.VoIndex = ft_vo.Index
		item.MoIndex = ft_mo.Index
		item.ValFlag = true
		arr = append(arr, item)
	}
	return arr
}

/*
//生成vo与mo的字段交集信息
//只有名称与类型相同的字段才属于字段交集
func getPubField4VoMoNest(arr VoFields, t_vo reflect.Type, t_mo reflect.Type, pos []int) VoFields {
	f_num := t_vo.NumField()
	for i := 0; i < f_num; i++ {
		ft_vo := t_vo.Field(i)
		if ft_vo.Anonymous == true {
			arr = getPubField4VoMoNest(arr, ft_vo.Type, t_mo, append(pos, i))
			continue
		}

		ft_mo, ok := t_mo.FieldByName(ft_vo.Name)
		if !ok {
			continue
		}
		if ft_vo.Type != ft_mo.Type {
			continue
		}

		var item VoInfo
		item.VoName = ft_vo.Name
		item.VoIndex = append(pos, i)
		item.MoIndex = ft_mo.Index
		item.ValFlag = true
		arr = append(arr, item)
	}
	return arr
}
*/

//获取与vo对应的mo的字段选中状态
//首先查询缓存，若不存在缓存，则调用LoadFromModel来生成字段选中状态
func getSelectFieldsByVo(md *DbModel, vo_ptr VoLoader) {
	g_voLoad_mutex.Lock()
	defer g_voLoad_mutex.Unlock()

	//从缓存中获取mo的字段状态
	//若存在缓存，则将字段状态设置到model的字段状态中
	t_vo := GetDirectType(reflect.TypeOf(vo_ptr))
	t_mo := GetDirectType(reflect.TypeOf(md.ent_ptr))
	key := t_vo.String() + t_mo.String()
	arr, ok := g_voLoad_cache[key]
	if ok {
		fld_num := len(arr)
		if fld_num == len(md.flds_addr) {
			for i := 0; i < len(arr); i++ {
				md.flds_addr[i].Flag = arr[i].Flag
			}
			return
		}
	}

	//若不存在缓存，则执行Vo的LoadFromModel来生成字段选中状态
	//LoadFromModel通常会调用:CopyDataFromModel来生成字段选中状态
	//CopyDataFromModel会调用getVoInfoJoinMo来获取字段交集
	//成字段选中状态后，将字段选中状态缓存起来
	vo_ptr.LoadFromModel(md, md.ent_ptr)
	fld_num := len(md.flds_addr)
	arr = make([]FieldValue, fld_num)
	for i := 0; i < fld_num; i++ {
		arr[i].Flag = md.flds_addr[i].Flag
	}
	g_voLoad_cache[key] = arr
	//fmt.Println("LoadFromModel")
}

//从struct（ent对象）获取与model对象字段名称类型相同的字段集合
func getSelectFieldsByEo(md *DbModel, vo_ptr interface{}) {
	if md.flag_omit == false {
		md.OmitALL()
		md.flag_omit = true
	}

	if vo_ptr == nil {
		return
	}
	v_vo := reflect.ValueOf(vo_ptr)
	if v_vo.Kind() != reflect.Ptr {
		return
	}
	v_vo = reflect.Indirect(v_vo)
	if v_vo.Kind() != reflect.Struct {
		return
	}

	//获取字段交集
	arr_join, err := getPubField4VoMo(vo_ptr, md.ent_ptr)
	if err != nil {
		return
	}

	v_mo := reflect.Indirect(reflect.ValueOf(md.ent_ptr))
	for _, item := range arr_join {
		fv_mo := v_mo.FieldByIndex(item.MoIndex)
		fld_ptr := fv_mo.Addr().Interface()
		md.set_flag_by_addr(fld_ptr, true)
	}
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

	if md != nil && md.flag_omit == false {
		md.OmitALL()
		md.flag_omit = true
	}

	//获取字段交集
	arr_join, err := getPubField4VoMo(vo_ptr, mo_ptr)
	if err != nil {
		return 0, err
	}
	//fmt.Println(arr_join)

	count := 0
	for _, item := range arr_join {
		fv_vo := v_vo.FieldByIndex(item.VoIndex)
		fv_mo := v_mo.FieldByIndex(item.MoIndex)
		if fv_vo.CanSet() == false {
			continue
		}
		fv_vo.Set(fv_mo)
		count += 1

		if md != nil {
			fld_ptr := fv_mo.Addr().Interface()
			md.set_flag_by_addr(fld_ptr, true)
		}
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

	if md != nil && md.flag_omit == false {
		md.OmitALL()
		md.flag_omit = true
	}

	//获取字段交集
	arr_join, err := getPubField4VoMo(vo_ptr, mo_ptr)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, item := range arr_join {
		fv_vo := v_vo.FieldByIndex(item.VoIndex)
		fv_mo := v_mo.FieldByIndex(item.MoIndex)
		if fv_mo.CanSet() == false {
			continue
		}
		fv_mo.Set(fv_vo)
		count += 1

		if md != nil {
			fld_ptr := fv_mo.Addr().Interface()
			md.set_flag_by_addr(fld_ptr, true)
		}
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
		md.SelectX(fld_ptr)
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
		md.SelectX(fld_ptr)
		return v_fld.Interface()
	}
}

func GetBool(md *DbModel, fld_ptr *bool) bool {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetInt(md *DbModel, fld_ptr *int) int {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetInt32(md *DbModel, fld_ptr *int32) int32 {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetInt64(md *DbModel, fld_ptr *int64) int64 {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetFloat32(md *DbModel, fld_ptr *float32) float32 {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetFloat64(md *DbModel, fld_ptr *float64) float64 {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetString(md *DbModel, fld_ptr *string) string {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetTime(md *DbModel, fld_ptr *time.Time) time.Time {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}
