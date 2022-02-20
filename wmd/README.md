### 介绍
wego/wmd是一款Go语言Markdown文档解析器，突出特点是解析性能高，支持基本的Markdown语法，并对图片以及表格进行了扩展。具体特征如下：
* 解析性能高。
* 支持图片块，可定义图片的aling属性以及图片的高度、宽度。
* 支持表格块的显示，可定义表格的aling属性、宽度属性、是否显示标题，以及单元格的aling属性。

### 安装
go get github.com/haming123/wego/wmd

### 快速上手
```go
func HandlerShowMd(w http.ResponseWriter, r *http.Request)  {
	file_name := "./demo.md"
	input, err := os.ReadFile(file_name)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}
	w.Write(wmd.MarshalHtml(input))
}
```

### 性能测试
#### demo.md解析测试
对于demo.md的内容（去掉扩展部分）进行解析，分别使用wmd以及blackfriday进行测试：
```go
func BenchmarkMarshalHtml(b *testing.B) {
	file_name := "../demo.md"
	input, err := os.ReadFile(file_name)
	if err != nil {
		b.Error(err)
		return
	}

	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MarshalHtml(input)
	}
	b.StopTimer()
}
```

```go
func BenchmarkMarshalHtml2(b *testing.B) {
	file_name := "./demo.md"
	input, err := os.ReadFile(file_name)
	if err != nil {
		b.Error(err)
		return
	}

	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		blackfriday.MarkdownCommon(input)
	}
	b.StopTimer()
}
```
测试的结果如下：
```
go test -v -run=none -bench="BenchmarkMarshalHtml" -benchmem
pkg: wmd
BenchmarkMarshalHtml-6             98636             11932 ns/op            6921 B/op         59 allocs/op
pkg: blackfriday
BenchmarkMarshalHtml2-6            42728             27885 ns/op           18900 B/op        242 allocs/op
```
#### README.md解析测试
README.md文件主要是代码和文字，测试结果如下：
```
pkg: wmd
BenchmarkMarshalHtml-6            126990              9226 ns/op            6113 B/op         33 allocs/op
pkg: blackfriday
BenchmarkMarshalHtml2-6            51812             23099 ns/op           22182 B/op        122 allocs/op
```