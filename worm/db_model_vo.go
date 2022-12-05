package worm

import (
	"reflect"
)

type VoSaver interface {
	SaveToModel(md *DbModel, mo interface{})
}

type VoLoader interface {
	LoadFromModel(md *DbModel, mo interface{})
}

//获取与vo对应的mo的字段选中状态
//若不存在缓存，则执行Vo的LoadFromModel来生成字段选中状态
//LoadFromModel通常会调用:CopyDataFromModel/GetXXX函数来生成字段选中状态
//CopyDataFromModel会调用getPubField4VoMo来获取字段交集
//生成字段选中状态后，将字段选中状态缓存起来
func selectFieldsByVo(md *DbModel, vo_ptr VoLoader) {
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

	//通过LoadFromModel来设置字段的选中状态
	vo_ptr.LoadFromModel(md, md.ent_ptr)
	//缓存vo的选择集
	g_selection_cache[cache_key] = md.flds_ext

	//计算缓存选择集与Model选择集的交集
	genSelectionByFieldIndex(md, md.flds_ext)
	//清空临时选择集
	md.flds_ext = nil
}
