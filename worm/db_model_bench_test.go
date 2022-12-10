package worm

import (
	"testing"
)

//go test -v -run=none -bench="BenchmarkDbQueryRow" -benchmem
func BenchmarkDbQueryRow(b *testing.B) {
	dbcnn, _ := OpenDb()
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var ent User
		rows, err := dbcnn.Query("select id,name,age from user where id=? limit 1", 1)
		if err != nil {
			b.Error(err)
			return
		}

		if !rows.Next() {
			b.Error(err)
			rows.Close()
			return
		}

		err = rows.Scan(&ent.DB_id, &ent.DB_name, &ent.Age)
		if err != nil {
			b.Error(err)
			rows.Close()
			return
		}
		rows.Close()
	}
	b.StopTimer()
}

//go test -v -run=none -bench="BenchmarkDbQueryRows" -benchmem
func BenchmarkDbQueryRows(b *testing.B) {
	dbcnn, _ := OpenDb()
	b.StopTimer()

	b.StartTimer()
	var arr []User
	for i := 0; i < b.N; i++ {
		var ent User
		rows, err := dbcnn.Query("select id,name,age from user where id>? and name is not null", 0)
		if err != nil {
			b.Error(err)
			return
		}

		for rows.Next() {
			err = rows.Scan(&ent.DB_id, &ent.DB_name, &ent.Age)
			if err != nil {
				b.Error(err)
				rows.Close()
				return
			}
			arr = append(arr, ent)
		}
		rows.Close()
	}
	b.StopTimer()
}

func BenchmarkModelUpdate(b *testing.B) {
	dbcnn, _ := OpenDb()
	InitEngine(&dialectMysql{}, dbcnn)
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var user = User{Age: 31, DB_name: "demo9"}
		_, err := Model(&user).ID(1).Update()
		if err != nil {
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

/*
go test -v -run=none -bench="BenchmarkModelGet" -benchmem
goos: windows
BenchmarkModelUpdate-6                46          27765943 ns/op            20 allocs/op
*/

//go test -v -run=none -bench="BenchmarkModelGet" -benchmem
func BenchmarkModelGet(b *testing.B) {
	dbcnn, _ := OpenDb()
	InitEngine(&dialectMysql{}, dbcnn)
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var ent User
		_, err := Model(&ent).Where("id=?", 1).Select("id", "name", "age").Get()
		if err != nil {
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

//go test -v -run=none -bench="BenchmarkModelGet" -benchmem
func BenchmarkModelGetWithPool(b *testing.B) {
	dbcnn, _ := OpenDb()
	InitEngine(&dialectMysql{}, dbcnn)
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var ent User
		md := pool_user.Get()
		_, err := md.Where("id=?", 1).Select("id", "name", "age").Get(&ent)
		if err != nil {
			b.Error(err)
			return
		}
		pool_user.Put(md)
	}
	b.StopTimer()
}

//go test -v -run=none -bench="BenchmarkModelWithCache" -benchmem
func BenchmarkModelGetWithCache(b *testing.B) {
	dbcnn, _ := OpenDb()
	InitEngine(&dialectMysql{}, dbcnn)
	ShowSqlLog(false)
	UsePrepare(true)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var ent User
		_, err := Model(&ent).Where("id=?", 1).Select("id", "name", "age").Get()
		if err != nil {
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

//go test -v -run=none -bench="BenchmarkModelFind" -benchmem
func BenchmarkModelFind(b *testing.B) {
	dbcnn, _ := OpenDb()
	InitEngine(&dialectMysql{}, dbcnn)
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var arr []User
		err := Model(&User{}).Where("id>?", 0).Select("id", "name", "age").Limit(10).Find(&arr)
		if err != nil {
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

func BenchmarkModelFindWithPool(b *testing.B) {
	dbcnn, _ := OpenDb()
	InitEngine(&dialectMysql{}, dbcnn)
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var arr []User
		md := pool_user.Get()
		err := md.Where("id>?", 0).Select("id", "name", "age").Limit(10).Find(&arr)
		if err != nil {
			b.Error(err)
			return
		}
		pool_user.Put(md)
	}
	b.StopTimer()
}

//go test -v -run=none -bench="BenchmarkModelWithCache" -benchmem
func BenchmarkModelFindWithCache(b *testing.B) {
	dbcnn, _ := OpenDb()
	InitEngine(&dialectMysql{}, dbcnn)
	ShowSqlLog(false)
	UsePrepare(true)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var arr []User
		err := Model(&User{}).Where("id>?", 0).Select("id", "name", "age").Limit(10).Find(&arr)
		if err != nil {
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
