package worm

import (
	"testing"
)

func BenchmarkGetVo1(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var ent UserVo
		_, err := Model(&User{}).Where("id=?", 1).Get(&ent)
		if err != nil{
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

func BenchmarkGetVo2(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var ent UserVo2
		_, err := Model(&User{}).Where("id=?", 1).Get(&ent)
		if err != nil{
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

func BenchmarkFindVo1(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var arr []UserVo
		err := Model(&User{}).Where("id>?", 0).Find(&arr)
		if err != nil{
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

func BenchmarkFindVo2(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var arr []UserVo2
		err := Model(&User{}).Where("id>?", 0).Find(&arr)
		if err != nil{
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

func BenchmarkJoinGetVo(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var vo UserBookVo
		tb := Model(&User{}).Select("id","name","age").TableAlias("u")
		tb.Join(&DB_Book{}, "b", "b.author=u.id", "name")
		_, err := tb.Where("u.id=?", 1).Get(&vo)
		if err != nil{
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

func BenchmarkJoinFindVo(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var datas []UserBookVo
		tb := Model(&User{}).Select("id","name","age").TableAlias("u")
		tb.LeftJoin(&DB_Book{}, "b", "b.author=u.id", "name")
		err := tb.Where("u.id>?", 0).Find(&datas)
		if err != nil{
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}
