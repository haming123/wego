package worm

import (
	"errors"
	"reflect"
	"sync"
)

type FieldIndex struct {
	FieldName string
	VoIndex   []int
	MoIndex   int
	//ValFlag   bool
}

type PublicFields struct {
	//Model字段索引号
	ModelField int
	//其他字段索引信息
	Fields []FieldIndex
}

func NewPublicFields(num int) *PublicFields {
	var flds PublicFields
	flds.ModelField = -1
	flds.Fields = make([]FieldIndex, num)
	flds.Fields = flds.Fields[:0]
	return &flds
}

//vo、mo字段交集信息缓存
var g_pflds_cache map[string]*PublicFields = make(map[string]*PublicFields)
var g_pflds_mutex sync.Mutex

//首先从缓存中获取字段交集
//若缓存中不存在，则生成字段交集
func getPubField4VoMo(vo_ptr interface{}, mo_ptr interface{}) (*PublicFields, error) {
	g_pflds_mutex.Lock()
	defer g_pflds_mutex.Unlock()

	var pflds *PublicFields = nil
	if vo_ptr == nil {
		return pflds, errors.New("vo_ptr is nil")
	}
	t_vo := reflect.TypeOf(vo_ptr)
	if t_vo.Kind() != reflect.Ptr {
		return pflds, errors.New("vo_ptr must be Pointer")
	}
	t_vo = GetDirectType(t_vo)
	if t_vo.Kind() != reflect.Struct {
		return pflds, errors.New("vo_ptr  muse be Struct")
	}

	if mo_ptr == nil {
		return pflds, errors.New("mo_ptr is nil")
	}
	t_mo := reflect.TypeOf(mo_ptr)
	if t_mo.Kind() != reflect.Ptr {
		return pflds, errors.New("mo_ptr must be Pointer")
	}
	t_mo = GetDirectType(t_mo)
	if t_mo.Kind() != reflect.Struct {
		return pflds, errors.New("mo_ptr  muse be Struct")
	}

	key := t_vo.String() + t_mo.String()
	val, ok := g_pflds_cache[key]
	if ok {
		return val, nil
	}

	pflds = NewPublicFields(t_mo.NumField())
	genPubField4VoMo(pflds, t_vo, t_mo)
	g_pflds_cache[key] = pflds

	return pflds, nil
}

//生成vo与mo的字段交集信息
//只有名称与类型相同的字段才属于字段交集
func genPubField4VoMo(pflds *PublicFields, t_vo reflect.Type, t_mo reflect.Type) {
	//遍历vo的结构体，看看是否有Model类型的字段
	//Model类型的字段，则获取字段索引，并退出（意味着选中全部Model的字段）
	//只进行第一级字段的检查
	f_num := t_vo.NumField()
	for i := 0; i < f_num; i++ {
		ft_vo := t_vo.Field(i)
		if ft_vo.Type == t_mo {
			pflds.ModelField = i
			return
		}
	}

	//获取vo结构体与Model共同的字段索引
	f_num = t_mo.NumField()
	for i := 0; i < f_num; i++ {
		ft_mo := t_mo.Field(i)
		ft_vo, ok := t_vo.FieldByName(ft_mo.Name)
		if !ok {
			continue
		}
		if ft_vo.Type != ft_mo.Type {
			continue
		}

		var item FieldIndex
		item.FieldName = ft_vo.Name
		item.VoIndex = ft_vo.Index
		item.MoIndex = i
		//item.ValFlag = true
		pflds.Fields = append(pflds.Fields, item)
	}
}

//从Eo对象（一般性的struct）获取与model对象字段名称类型相同的字段集合
func selectFieldsByEo(md *DbModel, vo_ptr interface{}) {
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
	pflds, err := getPubField4VoMo(vo_ptr, md.ent_ptr)
	if err != nil {
		return
	}

	//若vo中存在Model字段，意味着需要选择全部的Model字段
	if pflds.ModelField >= 0 {
		md.SelectALL()
		return
	}

	//根据字段索引获取model字段的地址，然后根据地址选中对应的字段
	for _, item := range pflds.Fields {
		md.set_flag_by_index(item.MoIndex, true)
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
	pflds, err := getPubField4VoMo(vo_ptr, mo_ptr)
	if err != nil {
		return 0, err
	}
	//fmt.Println(pflds)

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

		if md != nil {
			md.set_flag_by_index(item.MoIndex, true)
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
	pflds, err := getPubField4VoMo(vo_ptr, mo_ptr)
	if err != nil {
		return 0, err
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

		if md != nil {
			md.set_flag_by_index(item.MoIndex, true)
		}
	}
	return count, nil
}
