package worm

import (
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

func TestSqlRawGetRow(t *testing.T) {
	InitEngine4Test()

	val, err := SQL("select * from user where id=?", 1).GetRow()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(val)
}

func TestSqlRawRows(t *testing.T) {
	InitEngine4Test()

	rows, err := SQL("select * from user where id>?", 0).Rows()
	if err != nil {
		t.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = ScanModel(rows, &user)
		if err != nil {
			t.Error(err)
		}
		t.Log(user)
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

func TestSqlRawFindRow(t *testing.T) {
	InitEngine4Test()

	ret, err := SQL("select u.*, b.* from user u left join book b on b.author=u.id where u.id>?", 0).FindRow()
	if err != nil {
		t.Error(err)
		return
	}
	rr := ret.GetRowCount()
	t.Logf("row num = %d", rr)
	for i := 0; i < rr; i++ {
		t.Log(ret.GetRowData(i))
	}
}
