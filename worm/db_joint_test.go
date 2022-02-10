package worm

import (
	"testing"
)

func TestModelJoinGet (t *testing.T) {
	InitEngine4Test()

	type Data struct {
		User      	User
		Book      	DB_Book
	}
	var data Data

	tb := Model(&User{}).Select("id","name","age").TableAlias("u")
	tb.Join(&DB_Book{}, "b", "b.author=u.id", "name")
	_, err := tb.Where("u.id=?", 1).Get(&data)
	if err != nil{
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestModelJoinGetVo (t *testing.T) {
	InitEngine4Test()

	var vo UserBookVo
	tb := Model(&User{}).Select("id","name","age").TableAlias("u")
	tb.Join(&DB_Book{}, "b", "b.author=u.id", "")
	_, err := tb.Where("u.id=?", 1).Get(&vo)
	if err != nil{
		t.Error(err)
		return
	}
	t.Log(vo)
}

func TestModelJoinFind (t *testing.T) {
	InitEngine4Test()

	type Data struct {
		User      	User
		Book      	DB_Book
	}
	var datas []Data

	tb := Model(&User{}).Select("id","name","age").TableAlias("u")
	tb.LeftJoin(&DB_Book{}, "b", "b.author=u.id", "name")
	err := tb.WhereIn("u.id", 1,6).Find(&datas)
	if err != nil{
		t.Error(err)
		return
	}
	for _, item := range datas {
		t.Log(item)
	}
}

func TestModelJoinFindVo (t *testing.T) {
	InitEngine4Test()

	var datas []UserBookVo
	tb := Model(&User{}).Select("id","name","age").TableAlias("u")
	tb.LeftJoin(&DB_Book{}, "b", "b.author=u.id", "")
	err := tb.WhereIn("u.id", 1,6).Find(&datas)
	if err != nil{
		t.Error(err)
		return
	}
	for _, item := range datas {
		t.Log(item)
	}
}


