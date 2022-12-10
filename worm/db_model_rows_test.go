package worm

import "testing"

func TestModelRows(t *testing.T) {
	InitEngine4Test()

	rows, err := Model(&User{}).Where("id>?", 0).Limit(10).Rows()
	if err != nil {
		t.Error(err)
		return
	}

	for rows.Next() {
		var user User
		err = rows.Scan(&user)
		if err != nil {
			t.Error(err)
		}
		t.Log(user)
	}
	rows.Close()
}

func BenchmarkModelRows(b *testing.B) {
	dbcnn, _ := OpenDb()
	InitEngine(&dialectMysql{}, dbcnn)
	ShowSqlLog(false)
	UsePrepare(true)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rows, _ := Model(&User{}).Where("id>?", 0).Limit(10).Rows()
		for rows.Next() {
			var user User
			rows.Scan(&user)
		}
		rows.Close()
	}
	b.StopTimer()
}

func BenchmarkModelRows2(b *testing.B) {
	dbcnn, _ := OpenDb()
	InitEngine(&dialectMysql{}, dbcnn)
	ShowSqlLog(false)
	UsePrepare(true)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rows, _ := SQL("select * from user where id<?", 10).Rows()
		for rows.Next() {
			var user User
			scanModel(rows, &user)
		}
		rows.Close()
	}
	b.StopTimer()
}
