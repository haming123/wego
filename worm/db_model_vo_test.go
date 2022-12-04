package worm

import (
	"testing"
)

type UserVo struct {
	DB_id   int64
	DB_name string
	Age     int
}

func (vo *UserVo) LoadFromModel(md *DbModel, mo_ptr interface{}) {
	if mo, ok := mo_ptr.(*User); ok {
		vo.DB_id = GetInt64(md, &mo.DB_id)
		vo.DB_name = GetString(md, &mo.DB_name)
		vo.Age = GetInt(md, &mo.Age)
	}
}
func (vo *UserVo) SaveToModel(md *DbModel, mo_ptr interface{}) {
	if mo, ok := mo_ptr.(*User); ok {
		SetValue(md, &mo.DB_id, vo.DB_id)
		SetValue(md, &mo.DB_name, vo.DB_name)
		SetValue(md, &mo.Age, vo.Age)
	}
}

func TestModelInsertVo(t *testing.T) {
	InitEngine4Test()

	vo := UserVo{Age: 31, DB_name: "InsertVo2"}
	id, err := Model(&User{}).Insert(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(id)

	vo = UserVo{Age: 31, DB_name: "UpdateVo2"}
	ret, err := Model(&User{}).ID(id).Update(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ret)

	num, err := Model(&User{}).Where("id=?", id).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("delete num=%d", num)
}

func TestModelGetVo(t *testing.T) {
	InitEngine4Test()

	var vo UserVo
	_, err := Model(&User{}).Where("id=?", 1).Get(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(vo)
}

func TestModelFindVo(t *testing.T) {
	InitEngine4Test()

	var arr []UserVo
	err := Model(&User{}).Where("id>?", 0).Find(&arr)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range arr {
		t.Log(item)
	}
}
