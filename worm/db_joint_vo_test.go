package worm

import "testing"

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

func TestModelJoinGetVo(t *testing.T) {
	InitEngine4Test()

	var vo UserBookVo
	tb := Model(&User{}).Select("id", "name", "age").TableAlias("u")
	tb.Join(&DB_Book{}, "b", "b.author=u.id", "")
	_, err := tb.Where("u.id=?", 1).Get(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(vo)
}

func TestModelJoinFindVo(t *testing.T) {
	InitEngine4Test()

	var datas []UserBookVo
	tb := Model(&User{}).Select("id", "name", "age").TableAlias("u")
	tb.LeftJoin(&DB_Book{}, "b", "b.author=u.id", "")
	err := tb.WhereIn("u.id", 1, 6).Find(&datas)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range datas {
		t.Log(item)
	}
}

type UserBookVo2 struct {
	User
	BookName string
}

func (vo *UserBookVo2) LoadFromModel(md *DbModel, mo_ptr interface{}) {
	if mo, ok := mo_ptr.(*User); ok {
		CopyDataFromModel(md, vo, mo)
	} else if mo, ok := mo_ptr.(*DB_Book); ok {
		vo.BookName = GetString(md, &mo.DB_name)
	}
}

func TestModelJoinGetVo2(t *testing.T) {
	InitEngine4Test()

	var vo UserBookVo2
	tb := Model(&User{}).TableAlias("u")
	tb.Join(&DB_Book{}, "b", "b.author=u.id")
	_, err := tb.Where("u.id=?", 1).Get(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(vo)
}

func TestModelJoinFindVo2(t *testing.T) {
	InitEngine4Test()

	var datas []UserBookVo
	tb := Model(&User{}).TableAlias("u")
	tb.LeftJoin(&DB_Book{}, "b", "b.author=u.id")
	err := tb.WhereIn("u.id", 1, 6).Find(&datas)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range datas {
		t.Log(item)
	}
}
