package worm

import "testing"

func TestModelStatusGet(t *testing.T) {
	InitEngine4Test()

	var ent = DB_Book{}
	md := Model(&ent)
	GetInt64(md, &ent.DB_id)
	GetString(md, &ent.DB_name)
	_, err := md.Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ent)
}

func TestModelStatusUpdate(t *testing.T) {
	InitEngine4Test()

	var book = DB_Book{}
	md := Model(&book)
	SetValue(md, &book.DB_name, "c#")
	SetValue(md, &book.DB_author, 2)
	ret, err := md.ID(1).Update()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("update num=%d", ret)
}

type UserVo2 struct {
	DB_id   int64
	DB_name string
	Age     int
}

func (vo *UserVo2) LoadFromModel(md *DbModel, mo_ptr interface{}) {
	if mo, ok := mo_ptr.(*User); ok {
		CopyDataFromModel(md, vo, mo)
	}
}

func (vo *UserVo2) SaveToModel(md *DbModel, mo_ptr interface{}) {
	if mo, ok := mo_ptr.(*User); ok {
		CopyDataToModel(md, vo, mo)
	}
}

func TestModelInsertVo2(t *testing.T) {
	InitEngine4Test()

	vo := UserVo2{Age: 31, DB_name: "InsertVo2"}
	id, err := Model(&User{}).Insert(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(id)

	vo = UserVo2{Age: 31, DB_name: "UpdateVo2"}
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

func TestModelGetVo2(t *testing.T) {
	InitEngine4Test()

	var vo UserVo2
	_, err := Model(&User{}).Where("id=?", 1).Get(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(vo)
}
