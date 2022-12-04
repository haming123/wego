package worm

import (
	"testing"
)

type UserEo struct {
	DB_id   int64
	DB_name string
	Age     int
}

func TestModelEoIUD(t *testing.T) {
	InitEngine4Test()

	vo := UserEo{Age: 31, DB_name: "InsertEo"}
	id, err := Model(&User{}).Insert(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(id)

	vo = UserEo{Age: 31, DB_name: "UpdateEo"}
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

func TestModelGetEo(t *testing.T) {
	InitEngine4Test()

	var vo UserEo
	_, err := Model(&User{}).Where("id=?", 1).Get(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(vo)
}

func TestModelFindEo(t *testing.T) {
	InitEngine4Test()

	var arr []UserEo
	err := Model(&User{}).Where("id>?", 0).Find(&arr)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range arr {
		t.Log(item)
	}
}

type UserEo2 struct {
	User
	UserAttr string
}

func TestModelGetUserEo2(t *testing.T) {
	InitEngine4Test()

	var vo UserEo2
	_, err := Model(&User{}).Where("id=?", 1).Get(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(vo)
}

func TestModelFindEo2(t *testing.T) {
	InitEngine4Test()

	var arr []UserEo2
	err := Model(&User{}).Where("id>?", 0).Find(&arr)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range arr {
		t.Log(item)
	}
}

type UserEo3 struct {
	UserEo
	UserAttr string
}

func TestModelGetUserEo3(t *testing.T) {
	InitEngine4Test()

	var vo UserEo3
	_, err := Model(&User{}).Where("id=?", 1).Get(&vo)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(vo)
}

func TestModelFindEo3(t *testing.T) {
	InitEngine4Test()

	var arr []UserEo3
	err := Model(&User{}).Where("id>?", 0).Find(&arr)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range arr {
		t.Log(item)
	}
}
