#### 介绍
dlog是用GO语言实现的一个简单高效、支持文件轮换以及日志分级的日志SDK。其特征如下：
+ 采用文件日志类型采用了内存缓存，满足高性能输出日志。
+ 支持日志分级，具体分级如下：
    * fatal (log.LOG_FATAL)
    * error (log.LOG_ERROR)
    * warn  (log.LOG_WARN)
    * info  (log.LOG_INFO)
    * debug (log.LOG_DEBUG)
+ 支持终端日志类型以及可按照时间进行轮换的文件日志类型。
+ 文件日志类型支持按照天或小时进行轮换输出。

#### 安装说明
go get github.com/haming123/wego/dlog

#### 输出日志到终端
dlog的缺省日志类型为：TermLogger（终端日志类型），使用TermLogger时不需要初始化，可直接使用日志输出函数输出日志。
```go
package main
import log "dlog"
func main()  {
	log.Debug("This is a Debug Message")
	log.Info("This is a Info Message")
}
//执行后的输出结果为：
//2021/11/30 07:20:06 [D] main.go:31 This is a Debug Message
//2021/11/30 07:20:06 [I] main.go:32 This is a Debug Info
```

#### 日志级别
dlog支持5个日志级别，分别是：fatal、error、warn、info、debug。
```go
package main
import log "dlog"
func main()  {
	log.Debug("This is Debug Message")
	log.Info("This is Info Message")
	log.Warn("This is Warn Message")
	log.Error("This is Error Message")
	log.Fatal("This is Fatal Message")
}
```

#### 日志的格式化输出
dlog支持为每个级别的日志提供了一个xxxf、xxxJSON、xxxXML的函数，用于输出不同格式日志：
```go
package main
import log "dlog"
func main()  {
	log.Debugf("This is a %s Message", "Debug")
	log.Infof("This is a %s Message", "Info")
}
```
输出JSON格式的日志：
```go
package main
import log "dlog"
func main()  {
	type User struct {
		Name 	string
		Age 	int
	}
	user := User{Name:"lisi", Age:12}
	log.DebugJSON(user)
}
```
输出XML格式的日志：
```go
package main
import log "dlog"
func main()  {
	type User struct {
		Name 	string
		Age 	int
	}
	user := User{Name:"lisi", Age:12}
	log.DebugXML(user)
}
```
设置JSON/XML的显示格式：
```go
package main
import log "dlog"
func main()  {
	type User struct {
		Name 	string
		Age 	int
	}
	user := User{Name:"lisi", Age:12}
    log.ShowIndent(true)
	log.DebugJSON(user)
}
```

#### 使用Output输出日志
log.Output函数可以指定任意日志前缀输出日志：
```go
package main
import log "dlog"
func main()  {
	log.Output("[SQL]", "This is SQL Message")
	log.Output("[DEBUG]", "This is Debug Message")
}
```

#### 输出日志到文件
FileLogger是dlog提供的一种将日志输出到文件的日志类型。使用FileLogger日志类型前需要对FileLogger进行初始化，为FileLogger指定日志文件的存放目录以及日志文件的轮换方式。
日志文件可以按照天（log.ROTATE_DAY）或小时（log.ROTATE_HOUR）进行进行轮换。
```go
package main
import log "dlog"
func main()  {
	log.InitFileLogger("./logs", log.LOG_DEBUG)
	defer log.Close()

	log.Debug("This is a Debug Message")
	log.Info("This is a Debug Info")
}
```
dlog为了提高日志组件的性能，采用了基于buffer的文件输出方式，因此在系统退出前需要调用Close()函数将buffer中的数据刷新到文件中。

#### 关闭源码信息的输出
dlog缺省会在日志中输出源码所在的文件以及源码的行号，若不需要显示源码信息，可以使用函数：ShowCaller(show bool)来关闭源码信息的输出。
```go
package main
import log "dlog"
func main()  {
	log.InitFileLoggerHour("./logs", log.LOG_DEBUG)
    log.ShowCaller(false)
	defer log.Close()

	log.Debug("This is a Debug Message")
	log.Info("This is a Debug Info")
}
```

#### 性能指标
1. 显示caller
```go
goos: windows
goarch: amd64
BenchmarkFileLog
BenchmarkFileLog-6              983985		    1079 ns/op			4 allocs/op
```

2. 不显示caller
```go
goos: windows
goarch: amd64
BenchmarkFileLogNoCaller
BenchmarkFileLogNoCaller-6       6514368		172.3 ns/op			0 allocs/op
```