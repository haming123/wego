package worm

import "testing"

func TestSqlGetModelRows(t *testing.T) {
	InitEngine4Test()

	rows, err := SQL("select * from user where id<?", 10).ModelRows()
	if err != nil {
		t.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user)
		if err != nil {
			t.Error(err)
		}
		t.Log(user)
	}
}

func TestSQLBuilderGetModelRows(t *testing.T) {
	InitEngine4Test()

	rows, err := Table("user").Select("*").Where("id>?", 0).Limit(10).ModelRows()
	if err != nil {
		t.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(&user)
		if err != nil {
			t.Error(err)
		}
		t.Log(user)
	}
}

func BenchmarkGetModelRows(b *testing.B) {
	dbcnn, _ := OpenDb()
	InitEngine(&dialectMysql{}, dbcnn)
	ShowSqlLog(false)
	UsePrepare(true)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rows, _ := SQL("select * from user where id<?", 10).ModelRows()
		for rows.Next() {
			var user User
			rows.Scan(&user)
		}
		rows.Close()
	}
	b.StopTimer()
}
