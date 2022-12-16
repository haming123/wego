package worm

import (
	"reflect"
	"strings"
)

//获取字段的tag信息, tag的结构为：`db:"table.field"`, 或者是：`db:"table.*"`
func (lk *DbJoint) getVoFieldTagInfo(i_field reflect.StructField) (string, string) {
	table_db := ""
	field_db := ""
	tag := i_field.Tag.Get("db")
	parts := strings.Split(tag, ";")
	part0 := strings.Trim(parts[0], " ")
	if index := strings.Index(part0, "."); index >= 0 {
		table_db = part0[0:index]
		field_db = part0[index+1:]
	}
	return table_db, field_db
}

//通过数据库名称查询model的位置
func (lk *DbJoint) getModelIndexByDbTable(table_db string) int {
	for mm := 0; mm < len(lk.tables); mm++ {
		md := lk.tables[mm]
		if md.table_name == table_db {
			return mm
		}
	}
	return -1
}

//通过数据库字段名称查询model中字段的位置
func (lk *DbJoint) getModelFieldIndexByDbField(cache *JointEoFieldCache, table_inedex int, t_field reflect.Type, field_db string) (int, int) {
	var mo_field int = -1
	var mo_table int = -1
	for mm := 0; mm < len(lk.tables); mm++ {
		md := lk.tables[mm]
		if table_inedex >= 0 && mm != table_inedex {
			continue
		}
		//通过数据库字段名称查询model中字段的位置
		mo_index := md.get_field_index_dbname(field_db)
		if mo_index < 0 {
			continue
		}
		//只有类型与名称一致才算匹配上
		if md.flds_info[mo_index].FieldType != t_field {
			continue
		}
		mo_table = mm
		mo_field = mo_index
		break
	}

	return mo_table, mo_field
}

//查找第一个与Model结构体名称类型匹配的字段
func (lk *DbJoint) getModelFieldIndexByGoField(cache *JointEoFieldCache, table_inedex int, i_field reflect.StructField, t_field reflect.Type) (int, int) {
	var mo_field int = -1
	var mo_table int = -1
	for mm := 0; mm < len(lk.tables); mm++ {
		md := lk.tables[mm]
		if table_inedex >= 0 && mm != table_inedex {
			continue
		}
		//通过eo字段名称查询model中字段的位置
		mo_index := md.get_field_index_goname(i_field.Name)
		if mo_index < 0 {
			continue
		}
		//只有类型与名称一致才算匹配上
		if md.flds_info[mo_index].FieldType != t_field {
			continue
		}
		mo_table = mm
		mo_field = mo_index
		break
	}
	return mo_table, mo_field
}

//查找与Model名称、类型一致的字段，选中该字段，记录该字段的索引位置
func (lk *DbJoint) genPubField4VoMoNest(cache *JointEoFieldCache, t_vo reflect.Type, pos FieldPos, deep int, table_inedex int) {
	//超过最大层次，则退出
	if deep >= len(pos) {
		return
	}

	f_num := t_vo.NumField()
	for ff := 0; ff < f_num; ff++ {
		i_field := t_vo.Field(ff)
		t_field := i_field.Type

		//首先查找与Model类型一致的字段，并将字段的索引赋值给Model
		//只有第1层(deep=0)字段才判断是否为Model类型
		if deep < 1 && t_field.Kind() == reflect.Struct {
			is_model_field := false
			for mm := 0; mm < len(lk.tables); mm++ {
				md := lk.tables[mm]
				if t_field == md.ent_type {
					cache.models[mm].ModelField = ff
					is_model_field = true
					break
				}
			}
			if is_model_field {
				continue
			}
		}

		//只有第1层(deep=0)字段才获取字段的tag
		table_db := ""
		field_db := ""
		if deep < 1 {
			table_db, field_db = lk.getVoFieldTagInfo(i_field)
			if field_db == "*" {
				field_db = ""
			}
			table_inedex = -1
			if table_db != "" {
				table_inedex = lk.getModelIndexByDbTable(table_db)
			}
		}
		//每个字段，必须有对应的表名称
		if table_inedex < 0 {
			continue
		}

		//若是匿名字段，并且没有tag字段名称,则递归调用
		if i_field.Anonymous && field_db == "" {
			pos[deep] = ff
			lk.genPubField4VoMoNest(cache, t_field, pos, deep+1, table_inedex)
			continue
		}

		var mo_field int = -1
		var mo_table int = table_inedex
		//通过数据库字段名称查询model中字段的位置
		if field_db != "" {
			mo_table, mo_field = lk.getModelFieldIndexByDbField(cache, table_inedex, t_field, field_db)
		}
		//查找第一个与Model结构体名称类型匹配的字段
		if mo_field < 0 {
			mo_table, mo_field = lk.getModelFieldIndexByGoField(cache, table_inedex, i_field, t_field)
		}
		for mo_field >= 0 {
			var item FieldIndex
			item.FieldName = i_field.Name
			item.MoIndex = mo_field
			item.VoField = pos
			item.VoField[deep] = ff
			item.VoIndex = item.VoField[0 : deep+1]
			cache.models[mo_table].Fields = append(cache.models[mo_table].Fields, item)
			break
		}
	}
}

//通过eo(struct)对象来选择需要查询的字段
//首先从缓存中获取字段交集
//若缓存中不存在，则生成字段交集
func (lk *DbJoint) select_field_by_eo(t_vo reflect.Type) *JointEoFieldCache {
	g_joint_field_mutex.Lock()
	defer g_joint_field_mutex.Unlock()

	//从缓存中获取字段交集
	//若没有字段交集的缓存，调用genPubField4VoMoNest生成字段交集
	var cache *JointEoFieldCache = nil
	md_num := len(lk.tables)
	if md_num <= 3 {
		cache_key := fieldCacheKey4Join{}
		cache_key.t_vo = t_vo
		if md_num >= 1 {
			cache_key.t_mo0 = lk.tables[0].ent_type
		}
		if md_num >= 2 {
			cache_key.t_mo1 = lk.tables[1].ent_type
		}
		if md_num >= 3 {
			cache_key.t_mo2 = lk.tables[2].ent_type
		}
		cache, _ = g_joint_field_cache[cache_key]
		if cache == nil {
			var pos FieldPos
			cache = newJointEoFieldCache(lk.tables)
			lk.genPubField4VoMoNest(cache, t_vo, pos, 0, -1)
			g_joint_field_cache[cache_key] = cache
		}
	}

	if cache == nil {
		var pos FieldPos
		cache = newJointEoFieldCache(lk.tables)
		lk.genPubField4VoMoNest(cache, t_vo, pos, 0, -1)
	}

	//将公共字段添加到选择集中
	for mm := 0; mm < len(lk.tables) && mm < len(cache.models); mm++ {
		md := lk.tables[mm]
		pflds := &cache.models[mm]

		//若进行了字段的人工选择，则不需要进行字段的自动选择
		if md.flag_edit {
			continue
		}

		//若存在model类型一致的字段，不用额外选择字段（缺省选择全部）
		if pflds.ModelField < 0 {
			for _, item := range pflds.Fields {
				md.auto_add_field_index(item.MoIndex)
			}
		}
	}

	return cache
}

//把Model中地址的值赋值给vo对象
func (lk *DbJoint) CopyModelData2Eo(cache *JointEoFieldCache, v_vo reflect.Value) {
	for mm := 0; mm < len(lk.tables) && mm < len(cache.models); mm++ {
		md := lk.tables[mm]
		pflds := &cache.models[mm]

		//若vo中存在Model字段，只需要赋值Model对应的字段即可
		if pflds.ModelField >= 0 {
			fv_vo := v_vo.Field(pflds.ModelField)
			if fv_vo.CanSet() == true {
				fv_vo.Set(md.ent_value)
			}
		}

		for _, item := range pflds.Fields {
			fv_vo := v_vo.FieldByIndex(item.VoIndex)
			fv_mo := md.ent_value.Field(item.MoIndex)
			if fv_vo.CanSet() == false {
				continue
			}
			fv_vo.Set(fv_mo)
		}
	}
}

//通过vo对象来选择需要查询的字段
func (lk *DbJoint) select_field_by_vo(vo_ptr VoLoader) {
	for _, table := range lk.tables {
		selectFieldsByVo(table, vo_ptr)
	}
}

//把Model中地址的值赋值给vo对象
func (lk *DbJoint) CopyModelData2Vo(vo_ptr VoLoader) {
	for _, table := range lk.tables {
		vo_ptr.LoadFromModel(nil, table.ent_ptr)
	}
}
