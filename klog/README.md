## 介绍
klog用GO语言实现的一个高性能的结构化日志包，klog用于生成统计型日志。主要特征：
1) klog支持按照天或小时进行日志文件的滚动输出。
2) 采用基于buffer的文件输出方式，满足高性能输出日志。
3) 采用链式api来添加字段。

## 安装
go get github.com/haming123/wego/klog

## 快速上手
```go
package main
import "github.com/haming123/wego/klog"
func main() {
	klog.InitKlog("./logs", klog.ROTATE_DAY);
	defer klog.Close()

	klog.NewL("login").UserId("user_1").Add("login_type", "weixin").Output()
	klog.NewL("index").Client("127.0.0.1:666").Add("exe", 12).Output()
}
```
输出内容为：
```
ctm=2021/12/02 11:44:02 `class=login `userid=user_1 `login_type=weixin ` 
tm=2021/12/02 11:44:02 `class=index `client=127.0.0.1:666 `exe=12 ` 
```

## 日志格式以及常用字段
日志格式采用key=value的格式，字段之间采用" \`"来分隔，日志行之间采用" \\n" 来分隔，例如：
```
    k1=v1 `k2=v2 `...... \n
    k1=v1 `k2=v2 `...... \n
```
若字段的内容有"\`"符号，则编码为"\\\`"。同样若字段的内容有"\\n"符号，则编码为"\\\n"。例如若userid=hello\`, 则日志输出为：
```
ctm=2021/12/02 11:44:02 `class=demo `userid=hello\` `login_type=test `.
```
klog在输出日志是自动添加了日志的创建时间ctm，例如：
```
ctm=2021/12/02 11:44:02 `class=login `userid=user_1 `login_type=weixin ` 
```
klog设置了一下常用的字段：
1) class(ClassName):日志类型，必填字段，相当于与数据库的表明
2) func(FuncName):可选字段，通常为功能或页面的名称
3) client(Client):可选字段，客户端的IP
4) userid(UserId):可选字段，用户ID

## 代码的执行时间
klog支持在日志中输出代码的执行时间:
```go
package main
import (
	"github.com/haming123/wego/klog"
	"time"
)
func main() {
	var klog = klog.NewKlog("./logs", klog.ROTATE_DAY)
	defer klog.Close()
	bet_time := time.Now()

	//代码逻辑开始
	//...
	//代码逻辑结束

	klog.NewL("class_name").BeginTime(bet_time).UserId("user_name").Output()
}
```
输出内容为：
```
ctm=2021/12/02 12:32:29 `class=class_name `userid=user_name `etm=2 `  
```
其中etm便是代码的执行时间（单位:毫秒）

## 性能指标
```go
BenchmarkWriteLog-6      1223031		991.0 ns/op			5 allocs/op
```
