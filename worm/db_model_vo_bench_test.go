package worm

import (
	"testing"
)

func BenchmarkGetEo(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var ent UserEo
		_, err := Model(&User{}).Where("id=?", 1).Get(&ent)
		if err != nil {
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

func BenchmarkGetVo(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var ent UserVo
		_, err := Model(&User{}).Where("id=?", 1).Get(&ent)
		if err != nil {
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

func BenchmarkFindEo(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var arr []UserEo
		err := Model(&User{}).Where("id>?", 0).Find(&arr)
		if err != nil {
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}

func BenchmarkFindVo(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var arr []UserVo
		err := Model(&User{}).Where("id>?", 0).Find(&arr)
		if err != nil {
			b.Error(err)
			return
		}
	}
	b.StopTimer()
}
