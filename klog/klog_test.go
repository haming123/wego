package klog

import (
	"testing"
)

func BenchmarkWriteLog(b *testing.B) {
	InitKlog("./main/logs", ROTATE_HOUR);
	defer Close()
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		NewL("login").UserId("user_1").Add("login_type", "phone number").Add("sex", "男").Add("age", 12).Add("area", "guangdong").Add("float", 12.34).Add("bool", true).Add("int", 999).Output()
	}
	b.StopTimer()
}
/*
go test -v -run=none -bench="BenchmarkWriteLog" -benchmem
goos: windows
goarch: amd64
pkg: klog4go
BenchmarkWriteLog
BenchmarkWriteLog-6      1223031		991.0 ns/op			5 allocs/op

go test -v -run=none -bench="BenchmarkWriteLog" -benchmem
goos: linux
goarch: amd64
pkg: klog4go
BenchmarkWriteLog
BenchmarkWriteLog-2   	  715176		1572 ns/op			11 allocs/op
*/
