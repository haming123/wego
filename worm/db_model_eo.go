package worm

import (
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

//vo、mo字段交集缓存
var g_pubfield_cache map[string]*PublicFields = make(map[string]*PublicFields)
var g_pubfield_mutex sync.Mutex

//与vo对应的mo的model的selection缓存
var g_selection_cache map[string][]int = make(map[string][]int)
var g_selection_mutex sync.Mutex

//计算字段索引与Model选择集的交集
func genSelectionByFieldIndex(md *DbModel, fld_ext []int) {
	if fld_ext == nil {
		return
	}

	num := len(md.flds_addr)
	for i := 0; i < num; i++ {
		md.flds_addr[i].Flag = false
		if fld_ext[i] > 0 && md.flds_addr[i].Flag {
			md.flds_addr[i].Flag = true
		}
	}
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

//首先从缓存中获取字段交集
//若缓存中不存在，则生成字段交集
func getPubField4VoMo(cache_key string, t_vo reflect.Type, t_mo reflect.Type) (*PublicFields, error) {
	g_pubfield_mutex.Lock()
	defer g_pubfield_mutex.Unlock()

	if cache_key == "" {
		cache_key = t_vo.String() + t_mo.String()
	}
	val, ok := g_pubfield_cache[cache_key]
	if ok {
		return val, nil
	}

	pflds := NewPublicFields(t_mo.NumField())
	genPubField4VoMo(pflds, t_vo, t_mo)
	g_pubfield_cache[cache_key] = pflds

	return pflds, nil
}

//获取与Eo对象对应的mo的字段选中状态
func selectFieldsByEo(md *DbModel, vo_ptr interface{}) {
	g_selection_mutex.Lock()
	defer g_selection_mutex.Unlock()

	//获取选择集缓存
	//计算缓存选择集与Model选择集的交集
	t_vo := GetDirectType(reflect.TypeOf(vo_ptr))
	t_mo := GetDirectType(reflect.TypeOf(md.ent_ptr))
	cache_key := t_vo.String() + t_mo.String()
	if selection_ext, ok := g_selection_cache[cache_key]; ok {
		genSelectionByFieldIndex(md, selection_ext)
		return
	}

	//获取字段交集
	pflds, err := getPubField4VoMo(cache_key, t_vo, t_mo)
	if err != nil {
		return
	}

	//将公共字段添加到选扩展择集中
	//若存在model字段，不用设置扩展选择集
	if pflds.ModelField < 0 {
		for _, item := range pflds.Fields {
			md.add_field_ext_index(item.MoIndex)
		}
	}

	//缓存vo的选择集
	g_selection_cache[cache_key] = md.flds_ext
	//计算缓存选择集与Model选择集的交集
	genSelectionByFieldIndex(md, md.flds_ext)
	//清空临时选择集
	md.flds_ext = nil
}
