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
