package worm

import (
	"reflect"
	"testing"
)

func BenchmarkModelEoMatch(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var vo UserEo3
		md := Model(&User{})
		t_vo := GetDirectType(reflect.TypeOf(&vo))
		t_mo := GetDirectType(reflect.TypeOf(md.ent_ptr))
		pflds := NewPublicFields(t_mo.NumField())
		genPubField4VoMo(pflds, t_vo, t_mo)
		//b.Log(pflds)
	}
	b.StopTimer()
}

func BenchmarkModelEoMatch2(b *testing.B) {
	InitEngine4Test()
	ShowSqlLog(false)
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var vo UserEo3
		md := Model(&User{})
		t_vo := GetDirectType(reflect.TypeOf(&vo))
		pflds := NewPublicFields(md.ent_type.NumField())
		var pos FieldPos
		genPubField4VoMoNest(pflds, md, t_vo, pos, 0)
		//b.Log(pflds)
	}
	b.StopTimer()
}
