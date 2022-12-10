package worm

import (
	"testing"
)

func TestModelPool(t *testing.T) {
	InitEngine4Test()
	SetDebugLogLevel(LOG_DEBUG)
	dbs := NewSession()

	md := pool_user.Get(dbs)
	defer pool_user.Put(md)

	var ent User
	_, err := md.Select("id", "name", "age").ID(1).Get(&ent)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ent)
}

func TestModelPool2(t *testing.T) {
	InitEngine4Test()
	SetDebugLogLevel(LOG_DEBUG)
	dbs := NewSession()

	md := pool_user.Get(dbs)

	var ent User
	_, err := md.Select("id", "name", "age").ID(1).Get(&ent)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ent)

	pool_user.Put(md)
	pool_user.Put(md)
	_, err = md.Select("id", "name", "age").ID(1).Get(&ent)
	if err == nil {
		t.Error("model in pool")
		return
	}
}

func BenchmarkCreateModel(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Model(&User{})
	}
	b.StopTimer()
}

func BenchmarkCreateModelByPool(b *testing.B) {
	InitEngine4Test()
	//ShowSqlLog(false)
	b.StopTimer()

	dbs := NewSession()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		md := pool_user.Get(dbs)
		pool_user.Put(md)
	}
	b.StopTimer()
}

/*
go test -v -run=none -bench="BenchmarkCreateModel" -benchmem
BenchmarkCreateModel
BenchmarkCreateModel-6           2565315               426.2 ns/op           3 allocs/op
BenchmarkCreateModelByPool
BenchmarkCreateModelByPool-6    31093563                37.91 ns/op          0 allocs/op
*/
