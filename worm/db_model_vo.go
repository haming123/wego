package worm

type VoSaver interface {
	SaveToModel(md *DbModel, mo interface{})
}

type VoLoader interface {
	LoadFromModel(md *DbModel, mo interface{})
}

//获取与vo对应的mo的字段选中状态
//LoadFromModel通常会调用:CopyDataFromModel/GetXXX函数来生成字段选中状态
//CopyDataFromModel会调用getPubField4VoMo来获取字段交集
func selectFieldsByVo(md *DbModel, vo_ptr VoLoader) {
	//若进行了字段的人工选择，则不需要进行字段的自动选择
	if md.flag_edit == true {
		return
	}
	//通过LoadFromModel来设置字段的选中状态
	vo_ptr.LoadFromModel(md, md.ent_ptr)
}
