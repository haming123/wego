package worm

import "testing"

func TestSQLBuilderJoinRows2(t *testing.T) {
	InitEngine4Test()

	rows, err := Table("user").Select("*").Where("id>?", 0).Limit(10).StringRows()
	if err != nil {
		t.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		data := make(StringRow)
		err := rows.Scan(&data)
		if err != nil {
			t.Error(err)
		}
		t.Log(data)
	}
}

func TestSQLBuilderGetStringRows(t *testing.T) {
	InitEngine4Test()

	rows, err := Table("user").Select("*").Where("id>?", 0).Limit(10).StringRows()
	if err != nil {
		t.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		data := make(StringRow)
		err := rows.Scan(&data)
		if err != nil {
			t.Error(err)
		}
		t.Log(data)
	}
}

func TestSQLBuilderJoinRows(t *testing.T) {
	InitEngine4Test()

	tb := Table("user").Alias("u").Select("*").Where("u.id>?", 0).Limit(10)
	rows, err := tb.Join("book", "b", "b.author=u.id").StringRows()
	if err != nil {
		t.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		data := make(StringRow)
		err := rows.Scan(&data)
		if err != nil {
			t.Error(err)
		}
		t.Log(data)
	}
}

func BenchmarkStringRows(b *testing.B) {
	dbcnn, _ := OpenDb()
	InitEngine(&dialectMysql{}, dbcnn)
	ShowSqlLog(false)
	UsePrepare(true)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		rows, _ := Table("user").Select("*").Where("id>?", 0).Limit(10).StringRows()
		for rows.Next() {
			data := make(StringRow)
			rows.Scan(&data)
		}
		rows.Close()
	}
	b.StopTimer()
}
