package worm

import (
	"testing"
)

func BenchmarkDbQueryRow(b *testing.B) {
	_, err := OpenDb()
	if err != nil {
		b.Error(err)
		return
	}
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var ent User
		rows, err := DbConn.Query("select id,name,age from user where id=? limit 1", 9)
		if err != nil {
			b.Error(err)
			return
		}
		defer rows.Close()

		if !rows.Next() {
			b.Error(err)
		}
		err = rows.Scan(&ent.DB_id, &ent.DB_name, &ent.Age)
		if err != nil {
			b.Error(err)
			return
		}
		rows.Close()
	}
	b.StopTimer()
}

func BenchmarkDbQueryRows(b *testing.B) {
	_, err := OpenDb()
	if err != nil {
		b.Error(err)
		return
	}
	b.StopTimer()

	b.StartTimer()
	var arr []User
	for i := 0; i < b.N; i++ {
		var ent User
		rows, err := DbConn.Query("select id,name,age from user where id>? and name is not null", 0)
		if err != nil {
			b.Error(err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&ent.DB_id, &ent.DB_name, &ent.Age)
			if err != nil {
				b.Error(err)
				return
			}
			arr = append(arr, ent)
		}
		rows.Close()
	}
	b.StopTimer()
}

func BenchmarkModelUpdate(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var user = User{Age: 31, DB_name: "demo9"}
		_, err := Model(&user).ID(9).Update()
		if err != nil{
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

/*
go test -v -run=none -bench="BenchmarkModelUpdate" -benchmem
goos: windows
BenchmarkModelUpdate-6                46          27765943 ns/op            20 allocs/op
*/

func BenchmarkModelGet(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var ent User
		_, err := Model(&ent).Where("id=? or age>?", 9, 0).Select("id", "name", "age").Get()
		if err != nil{
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

func BenchmarkModelWithCache(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	UsePrepare(true)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var ent User
		_, err := Model(&ent).Where("id=? or age>?", 9, 0).Select("id", "name", "age").Get()
		if err != nil{
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

/*
go test -v -run=none -bench="BenchmarkModelGet" -benchmem
orm4go: windows
BenchmarkModelGet
BenchmarkModelGet-6					25195874 ns/op		31 allocs/op
BenchmarkModelGetWithCache
BenchmarkModelGetWithCache-6		11722227 ns/op		27 allocs/op

go test -v -run=none -bench="BenchmarkModelGet" -benchmem
orm4go: linux
BenchmarkModelGet
BenchmarkModelGet-2					6749190 ns/op	    34 allocs/op
BenchmarkModelGetWithCache
BenchmarkModelGetWithCache-2		3373236 ns/op	    31 allocs/op
*/

/*
go test -v -run=none -bench="BenchmarkModelGet" -benchmem
db_raw: windows
BenchmarkDbQuery2-6           		23373926 ns/op		22 allocs/op
orm4go: windows
BenchmarkModelGet-6           		25541364 ns/op      30 allocs/op
xorm: windows
BenchmarkModelGet-6           		27501973 ns/op		173 allocs/op
gorm: windows
BenchmarkModelGet-6           		25939504 ns/op		79 allocs/op
*/

/*
go test -v -run=none -bench="BenchmarkModelGet" -benchmem
db_raw: linux
BenchmarkDbQuery2-2   	     		6799005 ns/op		24 allocs/op
orm4go: linux
BenchmarkModelGet-2          		6749190 ns/op		34 allocs/op
xorm: linux
BenchmarkModelGet-2   	     		8010966 ns/op		172 allocs/op
gorm: linux
BenchmarkModelGet-2   	     		6754915 ns/op		73 allocs/op
*/

func BenchmarkModelFind(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var arr []User
		err := Model(&User{}).Where("id>? and name is not null", 0).Select("id", "name", "age").Find(&arr)
		if err != nil{
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}
//go test -v -run=none -bench="BenchmarkModelFind" -benchmem
