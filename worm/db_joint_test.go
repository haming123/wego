package worm

import (
	"testing"
)

type UserBookEo struct {
	User
	DB_Book
}

func TestModelJoinGet(t *testing.T) {
	InitEngine4Test()

	var data UserBookEo

	tb := Model(&User{}).Select("id", "name", "age").TableAlias("u")
	tb.Join(&DB_Book{}, "b", "b.author=u.id", "name")
	_, err := tb.Where("u.id=?", 1).Get(&data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestModelJoinFind(t *testing.T) {
	InitEngine4Test()

	var datas []UserBookEo

	tb := Model(&User{}).Select("id", "name", "age").TableAlias("u")
	tb.LeftJoin(&DB_Book{}, "b", "b.author=u.id", "name")
	err := tb.WhereIn("u.id", 1, 6).Find(&datas)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range datas {
		t.Log(item)
	}
}

type UserBookEo2 struct {
	User User
	Book DB_Book
}

func TestModelJoinGet2(t *testing.T) {
	InitEngine4Test()

	var data UserBookEo2

	tb := Model(&User{}).Select("id", "name", "age").TableAlias("u")
	tb.Join(&DB_Book{}, "b", "b.author=u.id", "name")
	_, err := tb.Where("u.id=?", 1).Get(&data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestModelJoinFind2(t *testing.T) {
	InitEngine4Test()

	var datas []UserBookEo2

	tb := Model(&User{}).Select("id", "name", "age").TableAlias("u")
	tb.LeftJoin(&DB_Book{}, "b", "b.author=u.id", "name")
	err := tb.WhereIn("u.id", 1, 6).Find(&datas)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range datas {
		t.Log(item)
	}
}

type UserBookEo3 struct {
	User
	AuthorId int64  `db:"book.author"`
	BookName string `db:"book.name"`
}

func TestModelJoinGet3(t *testing.T) {
	InitEngine4Test()

	var data UserBookEo3

	tb := Model(&User{}).Select("id", "name", "age").TableAlias("u")
	tb.Join(&DB_Book{}, "b", "b.author=u.id")
	_, err := tb.Where("u.id=?", 1).Get(&data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestModelJoinFind3(t *testing.T) {
	InitEngine4Test()

	var datas []UserBookEo3

	tb := Model(&User{}).Select("id", "name", "age").TableAlias("u")
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
