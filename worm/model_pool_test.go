package worm

import (
	"testing"
)

func TestModelPoolPut (t *testing.T) {
	InitEngine4Test()

	pool := &ModelPool{}
	md := Model(&User{})
	md.md_pool = pool
	pool.Put(md)
	t.Logf("size=%d\n", len(pool.pool))
	pool.Put(md)
	t.Logf("size=%d\n", len(pool.pool))
}

func TestModelPoolGetMoAutoPut (t *testing.T) {
	InitEngine4Test()
	SetDebugLogLevel(LOG_DEBUG)
	dbs := NewSession()

	var ent User
	_, err := NewModel_user(dbs, true).Select("id", "name", "age").ID(1).Get(&ent)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ent)
	_, err = NewModel_user(dbs, true).Select("id", "name", "age").ID(1).Get(&ent)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestModelPoolGetMo (t *testing.T) {
	InitEngine4Test()
	SetDebugLogLevel(LOG_DEBUG)
	dbs := NewSession()

	var ent User
	md := NewModel_user(dbs)
	_, err := md.Select("id", "name", "age").ID(1).Get(&ent)
	if err != nil {
		t.Error(err)
		md.PutToPool()
		return
	}
	t.Log(ent)
	md.PutToPool()
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
	ShowSqlLog(false)
	b.StopTimer()

	dbs := NewSession()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		md := NewModel_user(dbs)
		md.put_pool(md.md_pool)
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
