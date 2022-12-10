package worm

import (
	"fmt"
	"testing"
)

func TestSqlRawUpdate(t *testing.T) {
	InitEngine4Test()

	val, err := SQL("update user set age=22 where id=?", 1).Exec()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val.RowsAffected())
}

func TestSqlRawGetValues(t *testing.T) {
	InitEngine4Test()

	name := ""
	age := 0
	_, err := SQL("select name,age from user where id=?", 1).GetValues(&name, &age)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(name)
	t.Log(age)
}

func TestSqlRawGetString(t *testing.T) {
	InitEngine4Test()

	val, err := SQL("select name from user where id=?", 1).GetString()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val)
}

func TestSqlRawGetTime(t *testing.T) {
	InitEngine4Test()

	val, err := SQL("select created from user where id=?", 12).GetTime()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val)
}

func TestSqlRawGetModel(t *testing.T) {
	InitEngine4Test()

	var user User
	_, err := SQL("select * from user where id=?", 6).GetModel(&user)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}

func TestSqlRawRows(t *testing.T) {
	InitEngine4Test()

	rows, err := SQL("select name,age from user where id<?", 10).Rows()
	if err != nil {
		t.Error(err)
		return
	}
	for rows.Next() {
		var name string
		var age int
		err = Scan(rows, &name, &age)
		if err != nil {
			t.Error(err)
		}
		t.Log(name, age)
	}
	rows.Close()
}

func TestSqlRawFindValues(t *testing.T) {
	InitEngine4Test()

	var ids []int64
	var names []string
	num, err := SQL("select id,name from user where id>?", 0).FindValues(&ids, &names)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < num; i++ {
		str := fmt.Sprintf("id=%d, name=%s", ids[i], names[i])
		t.Log(str)
	}
}

func TestSqlRawFindString(t *testing.T) {
	InitEngine4Test()

	arr, err := SQL("select name from user where id>?", 0).FindString()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(arr)
}

func TestSqlRawFindModel(t *testing.T) {
	InitEngine4Test()

	var users []User
	err := SQL("select * from user where id>?", 5).FindModel(&users)
	if err != nil {
		t.Error(err)
		return
	}
	for _, user := range users {
		t.Log(user)
	}
}
