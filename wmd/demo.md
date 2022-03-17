## 介绍
wego/wmd是一款Go语言Markdown文档解析器，突出特点是解析性能高，支持基本的Markdown语法，并对图片以及表格进行了扩展。具体特征如下：
* 解析性能高。
* 支持图片块，可定义图片的aling属性以及图片的高度、宽度。
* 支持表格块的显示，可定义表格的aling属性、宽度属性、是否显示标题，以及单元格的aling属性。

## 安装
go get github.com/haming123/wego/wmd

## 解析md文件
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

## 基础语法
* 标题
使用"#"，可表示1-6级标题。例如：
```
# 标题1 
## 标题2 
### 标题3 
#### 标题4
##### 标题5
###### 标题6
```

* 水平线分割线
使用三个"-"开始，必须单独一行。例如：
```
---
```

* 无序列表
以"*"符号开始，"*"符号之后的空格不能少。例如：
```
* 列表项1
* 列表项2
    * 列表项2-1
    * 列表项2-2
```

* 引用
在引用的文字前加">"即可。引用也可以嵌套，例如加两个">"：
```
> 引用内容1
> 引用内容2
>> 嵌套引用2-1
>> 嵌套引用2-2
```

* 代码引用
需要引用代码时，如果引用的语句只有一段，不分行，可以用"\`"将语句包起来。例如：
```
这是一段代码`hello world`
```

* 代码块   
在代码块的前一行及后一行使用三个反引号 "`"，第一行反引号后面，可以添加代码块所使用的语言。例如：
```
` ``go
//go代码块
var aaa int = 1
fmt.Println(aaa)
` ``
```

* 粗体文字
在要显示为粗体的文字使用前后两个"*"符号，例如：
```
这是一段**粗体**文字
```

* 文字删除线
在要显示删除线的文字使用前后两个"~"符号，例如：
```
这是一段~~有删除线的~~文字
```

* 插入链接
```
语法：[超链接名](超链接地址 "超链接title")
title可加可不加， 例如：
[请点击](http://link.addr "我是标题")
```

* 插入图片
```
语法：![图片alt](图片地址 ''图片title'')
图片alt就是显示在图片下面的文字，相当于对图片内容的解释。图片title是图片的标题，当鼠标移到图片上时显示的内容。title可加可不加。
使用“#”号分隔图片地址与图片的属性。图片的属性格式为：width=xxx&height=xxx。其中国width是图片的宽度，height是图片的高度。
例如：
![图片下面的文字](http://image-path.png#width=300px&height=200px)
```

## 扩展语法
* 有序列表
以"+"符号开始，"+"符号之后的空格不能少。例如：
```
+ 列表项1
+ 列表项2
```

* 倾斜文字
在要倾斜显示的文字使用前后两个"/"符号，例如：
```
这是一段//倾斜的//文字
```

* 文字下划线
在要显示下划线的文字使用前后两个"_"符号，例如：
```
这是一段__有下划线的__文字
```

* 续行
有时不希望一行文字形成一个段落，则可以使用符号："<"来续行。例如：
```
这是一段文字
< 这些文字将与上面的文字显示为一行
```

* 表格块
```

` ``table
表头1|表头2
:----:|:----
单元格11|单元格12
单元格21|单元格22

` ``
说明：
1）表格属性：
表格居中：表头|表头
表格靠左：<表头|表头
表格靠右：表头|表头>
表格拉伸：<表头|表头>
2）单元格属性：
单元格靠左 :----
单元格靠右 ----:
单元格居中 ----
```

* 图片块
```

` ``image
    http://image-path.png#align=center&width=300px&height=200px

` ``
说明：
使用“#”号分隔图片地址，后面的属性为：
1）aling
align=center：图片居中
align=left：图片靠左
align=right：图片靠右
2）width
图片的宽度
3）height
图片的高度
```

## 性能测试
### demo.md解析测试
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

### README.md解析测试
README.md文件主要是代码和文字，测试结果如下：
```
pkg: wmd
BenchmarkMarshalHtml-6            126990              9226 ns/op            6113 B/op         33 allocs/op
pkg: blackfriday
BenchmarkMarshalHtml2-6            51812             23099 ns/op           22182 B/op        122 allocs/op
```