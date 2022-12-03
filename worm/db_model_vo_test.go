package worm

import (
	log "github.com/haming123/wego/dlog"
	"testing"
)

type UserVo struct {
	DB_id   int64
	DB_name string
	Age     int
}

type UserVo2 struct {
	DB_id   int64
	DB_name string
	Age     int
}

func (vo *UserVo2) LoadFromModel(md *DbModel, mo_ptr interface{}) {
	if mo, ok := mo_ptr.(*User); ok {
		vo.DB_id = GetInt64(md, &mo.DB_id)
		vo.DB_name = GetString(md, &mo.DB_name)
		vo.Age = GetInt(md, &mo.Age)
	}
}
func (vo *UserVo2) SaveToModel(md *DbModel, mo_ptr interface{}) {
	if mo, ok := mo_ptr.(*User); ok {
		SetValue(md, &mo.DB_id, vo.DB_id)
		SetValue(md, &mo.DB_name, vo.DB_name)
		SetValue(md, &mo.Age, vo.Age)
	}
}

type UserBookVo struct {
	UserId   int64
	UserName string
	UserAge  int
	BookName string
}

func (vo *UserBookVo) LoadFromModel(md *DbModel, mo_ptr interface{}) {
	if mo, ok := mo_ptr.(*User); ok {
		vo.UserId = GetInt64(md, &mo.DB_id)
		vo.UserName = GetString(md, &mo.DB_name)
		vo.UserAge = GetInt(md, &mo.Age)
	} else if mo, ok := mo_ptr.(*DB_Book); ok {
		vo.BookName = GetString(md, &mo.DB_name)
	}
}

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

func TestModelInsertVo(t *testing.T) {
	InitEngine4Test()

	vo := UserVo{Age: 31, DB_name: "InsertVo"}
	id, err := Model(&User{}).Insert(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(id)

	num, err := Model(&User{}).Where("id=?", id).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("delete num=%d", num)
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

func TestModelUpdateVo(t *testing.T) {
	InitEngine4Test()

	vo := UserVo{Age: 31, DB_name: "UpdateVo"}
	ret, err := Model(&User{}).ID(1).Update(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ret)
}

func TestModelUpdateVo2(t *testing.T) {
	InitEngine4Test()

	vo := UserVo2{Age: 31, DB_name: "UpdateVo2"}
	ret, err := Model(&User{}).ID(1).Update(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ret)
}

func TestModelFindVo1(t *testing.T) {
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

func TestModelFindVo2(t *testing.T) {
	InitEngine4Test()

	var arr []UserVo2
	err := Model(&User{}).Where("id>?", 0).Find(&arr)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range arr {
		t.Log(item)
	}
}

type BookVo1 struct {
	DB_id     int64
	DB_author int64
	DB_name   string
}

type BookVo2 struct {
	BookVo1
	DB_remark string
}

func TestModelGetBookVo1(t *testing.T) {
	InitEngine4Test()
	log.ShowIndent(true)

	var vo BookVo2
	_, err := Model(&DB_Book{}).Where("id=?", 1).Get(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(vo)
}
