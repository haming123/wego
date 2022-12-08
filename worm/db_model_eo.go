package worm

import (
	"reflect"
	"sync"
)

type FieldIndex2 struct {
	FieldName string
	PIndex    int
	FIndex    int
	MoIndex   int
}

//匿名字段嵌套的最大深度
const MaxNestDeep int = 10

type FieldPos [MaxNestDeep]int
type FieldIndex struct {
	FieldName string
	VoIndex   []int
	MoIndex   int
	VoField   FieldPos
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
		pflds.Fields = append(pflds.Fields, item)
	}
}

//生成vo与mo的字段交集信息
//只有名称与类型相同的字段才属于字段交集
//deep必须从0开始
func genPubField4VoMoNest(pflds *PublicFields, md *DbModel, t_vo reflect.Type, pos FieldPos, deep int) {
	//超过最大层次，则退出
	if deep >= len(pos) {
		return
	}

	f_num := t_vo.NumField()
	for ff := 0; ff < f_num; ff++ {
		ft_vo := t_vo.Field(ff)

		//只有第1层(deep=0)字段才判断是否为Model类型
		//若存在Model类型的字段, 则直接退出（意味着选中全部Model的字段）
		if deep < 1 {
			if ft_vo.Type == md.ent_type {
				pflds.ModelField = ff
				return
			}
		}

		//若是匿名字段,则递归调用
		if ft_vo.Anonymous {
			pos[deep] = ff
			genPubField4VoMoNest(pflds, md, ft_vo.Type, pos, deep+1)
			continue
		}

		//通过eo字段名称查询model中字段的位置
		mo_index := md.get_field_index_byname2(ft_vo.Name)
		if mo_index < 0 {
			continue
		}

		//只有类型与名称一致才算匹配上
		if md.flds_info[mo_index].FieldType != ft_vo.Type {
			continue
		}

		var item FieldIndex
		item.FieldName = ft_vo.Name
		item.MoIndex = mo_index
		item.VoField = pos
		item.VoField[deep] = ff
		item.VoIndex = item.VoField[0 : deep+1]
		pflds.Fields = append(pflds.Fields, item)
	}
}

//首先从缓存中获取字段交集
//若缓存中不存在，则生成字段交集
func getPubField4VoMo(md *DbModel, cache_key string, t_vo reflect.Type, t_mo reflect.Type) (*PublicFields, error) {
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
	//genPubField4VoMo(pflds, t_vo, t_mo)
	var pos FieldPos
	genPubField4VoMoNest(pflds, md, t_vo, pos, 0)
	g_pubfield_cache[cache_key] = pflds

	return pflds, nil
}

//获取与Eo对象对应的mo的字段选中状态
func selectFieldsByEo(md *DbModel, vo_ptr interface{}) {
	//获取字段交集
	t_vo := GetDirectType(reflect.TypeOf(vo_ptr))
	t_mo := GetDirectType(reflect.TypeOf(md.ent_ptr))
	cache_key := t_vo.String() + t_mo.String()
	pflds, err := getPubField4VoMo(md, cache_key, t_vo, t_mo)
	if err != nil {
		return
	}

	//若进行了字段的人工选择，则不需要进行字段的自动选择
	if md.flag_edit == true {
		return
	}

	//将公共字段添加到选择集中
	//若存在model类型一致的字段，不用额外选择字段（缺省选择全部）
	if pflds.ModelField < 0 {
		for _, item := range pflds.Fields {
			md.auto_add_field_index(item.MoIndex)
		}
	}
}
