package worm

import (
	"reflect"
	"sync"
)

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
//deep必须从0开始
func (pflds *PublicFields) genPubField4VoMoNest(moinfo *ModelInfo, t_mo reflect.Type, t_vo reflect.Type, pos FieldPos, deep int) {
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
			if ft_vo.Type == t_mo {
				pflds.ModelField = ff
				return
			}
		}

		//若是匿名字段,则递归调用
		if ft_vo.Anonymous {
			pos[deep] = ff
			pflds.genPubField4VoMoNest(moinfo, t_mo, ft_vo.Type, pos, deep+1)
			continue
		}

		//通过eo字段名称查询model中字段的位置
		mo_index := moinfo.get_field_index_goname(ft_vo.Name)
		if mo_index < 0 {
			continue
		}
		//只有类型与名称一致才算匹配上
		if moinfo.Fields[mo_index].FieldType != ft_vo.Type {
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

//获取与Eo对象对应的mo的字段交集
//首先从缓存中获取字段交集
//若缓存中不存在，则生成字段交集
func getPubField4VoMo(t_mo reflect.Type, t_vo reflect.Type) *PublicFields {
	g_pubfield_mutex.Lock()
	defer g_pubfield_mutex.Unlock()

	//获取字段交集
	cache_key := t_vo.String() + t_mo.String()
	pflds, ok := g_pubfield_cache[cache_key]
	if ok {
		return pflds
	}

	//获取model信息，创建字段交集对象，并生成字段交集
	var pos FieldPos
	moinfo := getModelInfo(t_mo)
	pflds = NewPublicFields(len(moinfo.Fields))
	pflds.genPubField4VoMoNest(moinfo, t_mo, t_vo, pos, 0)
	g_pubfield_cache[cache_key] = pflds
	return pflds
}

//获取与Eo对象对应的mo的字段选中状态
//首先从缓存中获取字段交集
//若缓存中不存在，则生成字段交集
func (md *DbModel) selectFieldsByEo(t_vo reflect.Type) {
	//获取字段交集
	cache := getPubField4VoMo(md.ent_type, t_vo)
	md.VoFields = cache

	//若进行了字段的人工选择，则不需要进行字段的自动选择
	if md.flag_edit == true {
		return
	}

	//将公共字段添加到选择集中
	//若存在model类型一致的字段，不用额外选择字段（缺省选择全部）
	if cache.ModelField < 0 {
		for _, item := range cache.Fields {
			md.auto_add_field_index(item.MoIndex)
		}
	}
}

//把Model中地址的值赋值给vo对象
func (md *DbModel) CopyModelData2Eo(v_vo reflect.Value) {
	if md.VoFields == nil {
		panic(" md.VoFields == nil")
	}

	//若vo中存在Model字段，只需要赋值Model对应的字段即可
	pflds := md.VoFields
	if pflds.ModelField >= 0 {
		fv_vo := v_vo.Field(pflds.ModelField)
		if fv_vo.CanSet() == true {
			fv_vo.Set(md.ent_value)
			return
		}
	}

	//遍历字段交集，逐个给vo的对象赋值
	for _, item := range pflds.Fields {
		fv_vo := v_vo.FieldByIndex(item.VoIndex)
		fv_mo := md.ent_value.Field(item.MoIndex)
		if fv_vo.CanSet() == false {
			continue
		}
		fv_vo.Set(fv_mo)
	}
}
