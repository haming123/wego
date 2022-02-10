package log

import (
	"testing"
)

func TestLogger(t *testing.T) {
	Debug("This is Debug Message")
	Info("This is Info Message")
	Warn("This is Warn Message")
	Error("This is Error Message")
	Fatal("This is Fatal Message")
}

func TestLoggerFormt(t *testing.T) {
	Debugf("This is a %s Message", "Debug")
	Infof("This is a %s Message", "Info")
}

func TestLoggerJSON(t *testing.T) {
	type User struct {
		Name 	string
		Age 	int
	}
	user := User{Name:"lisi", Age:12}
	DebugJSON(user)
}

func TestLoggerJSON2(t *testing.T) {
	type User struct {
		Name 	string
		Age 	int
	}
	ShowIndent(true)
	user := User{Name:"lisi", Age:12}
	DebugJSON(user)
}

func TestLoggerXML(t *testing.T) {
	type User struct {
		Name 	string
		Age 	int
	}
	user := User{Name:"lisi", Age:12}
	DebugXML(user)
}

func TestLogger_Output(t *testing.T) {
	Output("[SQL]", "This is SQL Message")
	Output("[DEBUG]", "This is Debug Message")
}

func TestLoggerFile(t *testing.T) {
	InitFileLogger("../demo/logs", LOG_DEBUG)
	defer Close()

	Debug("This is a Debug Message")
	Info("This is a Debug Info")
}

func TestLoggerShowCaller(t *testing.T) {
	InitFileLogger("../demo/logs", LOG_DEBUG)
	ShowCaller(false)
	defer Close()

	Debug("This is a Debug Message")
	Info("This is a Debug Info")
}

func BenchmarkTermLog(b *testing.B) {
	log := NewTermLogger(LOG_DEBUG)
	for i := 0; i < b.N; i++ {
		log.Info("This is a log message")
	}
}

func BenchmarkFileLog(b *testing.B) {
	log := InitFileLogger("./", LOG_DEBUG)
	defer log.Close()
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Info("This is a log message")
	}
	b.StopTimer()
}

//go test -v -run=none -bench="BenchmarkFileLog" -benchmem
/*
goos: windows
goarch: amd64
pkg: log4go
BenchmarkFileLog
BenchmarkFileLog-6        983985		1079 ns/op			256 B/op		4 allocs/op

goos: linux
goarch: amd64
pkg: log4go
BenchmarkFileLog
BenchmarkFileLog-2   	  407358		2546 ns/op			456 B/op		10 allocs/op
*/

func BenchmarkFileLogNoCaller(b *testing.B) {
	log := InitFileLogger("./", LOG_DEBUG)
	log.ShowCaller(false)
	defer log.Close()
	b.StopTimer()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Info("This is a log message")
	}
	b.StopTimer()
}
//go test -v -run=none -bench="BenchmarkFileLogNoCaller" -benchmem
/*
goos: windows
goarch: amd64
pkg: log4go
BenchmarkFileLogNoCaller
BenchmarkFileLogNoCaller-6       6514368		172.3 ns/op			0 B/op			0 allocs/op

goos: linux
goarch: amd64
pkg: log4go
BenchmarkFileLogNoCaller
BenchmarkFileLogNoCaller-2   	 1996180		600 ns/op			192 B/op		6 allocs/op
*/