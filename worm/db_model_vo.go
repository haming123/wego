package worm

import (
	"reflect"
	"sync"
)

type VoSaver interface {
	SaveToModel(md *DbModel, mo interface{})
}

type VoLoader interface {
	LoadFromModel(md *DbModel, mo interface{})
}

//与vo对应的mo的model的字段选中状态缓存
var g_voselect_cache map[string][]FieldValue = make(map[string][]FieldValue)
var g_voselect_mutex sync.Mutex

//获取与vo对应的mo的字段选中状态
//若不存在缓存，则执行Vo的LoadFromModel来生成字段选中状态
//LoadFromModel通常会调用:CopyDataFromModel/GetXXX函数来生成字段选中状态
//CopyDataFromModel会调用getVoInfoJoinMo来获取字段交集
//生成字段选中状态后，将字段选中状态缓存起来
func selectFieldsByVo(md *DbModel, vo_ptr VoLoader) {
	g_voselect_mutex.Lock()
	defer g_voselect_mutex.Unlock()

	//从缓存中获取mo的字段状态
	//若存在缓存，则将字段状态设置到model的字段状态中
	t_vo := GetDirectType(reflect.TypeOf(vo_ptr))
	t_mo := GetDirectType(reflect.TypeOf(md.ent_ptr))
	key := t_vo.String() + t_mo.String()
	if flds_vo, ok := g_voselect_cache[key]; ok {
		if len(flds_vo) == len(md.flds_addr) {
			for i := 0; i < len(flds_vo); i++ {
				md.flds_addr[i].Flag = flds_vo[i].Flag
			}
			return
		}
	}

	//备份md原来的字段状态
	flds_md := make([]FieldValue, len(md.flds_addr))
	copy(flds_md, md.flds_addr)

	//清空字段的选择状态
	md.OmitALL()
	//通过vo来设置字段的选中状态
	vo_ptr.LoadFromModel(md, md.ent_ptr)
	//获取vo的字段选择状态
	flds_vo := md.flds_addr
	//恢复md的字段状态
	md.flds_addr = flds_md

	//缓存vo的字段选择状态
	g_voselect_cache[key] = flds_vo
	//将vo字段的选择状态叠加到model的字段状态上
	if len(flds_vo) == len(md.flds_addr) {
		for i := 0; i < len(flds_vo); i++ {
			md.flds_addr[i].Flag = flds_vo[i].Flag
		}
	}
}
