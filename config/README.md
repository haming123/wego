# wego/config

### 介绍
wego/config是一款GO语言版本的ini配置文件解析工具，wego/config具有以下特征：
首先准备app.conf配置文件:
```
1）提供了GetString、GetInt...以及MustrString、MustInt...函数，方便配置数据的获取。
2）支持通过struct的tag来自动将配置数据赋值给struct的字段。
3）支持使用环境变量配置项的值。
```
### 安装
go get github.com/haming123/wego/config

### 快速上手
首先准备app.conf配置文件:
```
#应用名称
app_name = demo
#日志级别: 0 OFF 1 FATAL 2 ERROR 3 WARN 4 INFO 5 DEBUG
level = 5
#获取环境变量
go_path = ${GOPATH}

[mysql]
db_name = demodb
db_host = 127.0.0.1:3306
db_user = root
db_pwd = demopwd
```
接下看看如何读取配置文件，并获取配置项的值：
```go
package main
import (
	"fmt"
	"wego/config"
)
func main()  {
	var cfg config.ConfigData
	err := config.ParseFile("./app.conf", &cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	//通过GetXXX获取配置项的值
	val := cfg.GetString("app_name")
	if val.Error != nil {
		fmt.Println( val.Error)
		return
	}
	fmt.Println(val.Value)
}
```

### 读取Section中的配置项
ini 文件是以分区（section）组织的。分区以[name]开始，在下一个分区前结束。所有分区前的内容属于默认分区（[root]）。以下代码是section配置项的读取示例：
```go
func TestSectionGet(t *testing.T) {
	var cfg config.ConfigData
	err := config.ParseFile("./app.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	//通过GetXXX获取配置项的值
	val := cfg.Section("mysql").GetString("db_name")
	if val.Error != nil {
		t.Error(err)
		return
	}
	t.Log(val.Value)
}
```

### 各种类型的数据的读取
为了方便各种类型的配置数据的获取, wego/config提供了GetString、GetInt...等函数，例如以下配置文件的读取：
```
#各种数据类型的配置项
str_value = hello
bool_value = true
int_value = 99
float_value = 123.45
```
```go
func TestIniGetXXX(t *testing.T) {
	var cfg config.ConfigData
	err := config.ParseFile("./app2.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	val := cfg.GetString("str_value")
	if val.Error != nil {
		t.Error(val.Error)
		return
	}
	t.Log(val.Value)

	val_bool := cfg.GetString("bool_value")
	if val.Error != nil {
		t.Error(val_bool.Error)
		return
	}
	t.Log(val_bool.Value)

	val_int := cfg.GetString("int_value")
	if val.Error != nil {
		t.Error(val_int.Error)
		return
	}
	t.Log(val_int.Value)

	val_float := cfg.GetString("float_value")
	if val.Error != nil {
		t.Error(val_float.Error)
		return
	}
	t.Log(val_float.Value)
}
```

### 数据的快捷读取
使用GetXXX函数读取配置项需要进行错误判断，这样的代码写起来会非常繁琐。为此，wego/config提供对应的MustXXX方法，这个方法只返回一个值， 
同时它可接受缺省参数，如果没有配置对应的配置项或配置内容无法转换，则使用缺省值作为返回值。
```go
func TestIniMustXXX(t *testing.T) {
	var cfg config.ConfigData
	err := config.ParseFile("./app2.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(cfg.MustString("str_value"))
	t.Log(cfg.MustBool("bool_value"))
	t.Log(cfg.MustInt("int_value"))
	t.Log(cfg.MustFloat("float_value"))
}
```

#### 数组类型数据的读取
wego/config也支持数组类型数据的读取，要求：数组要作为一个配置项添加到ini文件中，并且数组成员之间用指定的分隔符（例如“,”）分隔：
```
#数组配置项
ints_value = 1,2,3,4,5
```
```go
func TestGetArray(t *testing.T) {
	var cfg config.ConfigData
	err := config.ParseFile("./app2.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	arr, err := cfg.GetInts("ints_value", ",")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(arr)
}
```
#### 结构体字段映射与数据的读取
wego/config持通过struct的tag来获取struct字段与配置项的映射关系，并可以通过映射关系自动给struct字段赋值，首先需要在struct定义中指定映射关系：
```go
type DbConfig struct {
	MysqlHost 	string 		`ini:"db_host"`
	MysqlUser 	string 		`ini:"db_user"`
	MysqlPwd  	string 		`ini:"db_pwd"`
	MysqlDb   	string 		`ini:"db_name"`
}

type AppConfig struct {
	AppName  	string   	`ini:"app_name"`
	HttpPort	uint     	`ini:"http_port;default=8080"`
	GoPath   	string   	`ini:"go_path"`
	DbParam  	DbConfig 	`ini:"mysql"`
}
```
说明：
wego/config的定义映射关系时支持配置缺省值，再进行数据解析时若没有配置内容，则使用缺省值作为字段的值。
wego/config使用GetStruct来struct字段赋值，例如：
```go
func TestIniGetStruct(t *testing.T) {
	var cfg config.ConfigData
	err := config.ParseFile("./app.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	var data AppConfig
	err = cfg.GetStruct(&data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}
```

也可以调用Section的GetStruct函数直接从section中获取配置内容，例如：
```go
func TestIniGetSectionStruct(t *testing.T) {
	var cfg config.ConfigData
	err := config.ParseFile("./app.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	var data DbConfig
	err = cfg.Section("mysql").GetStruct(&data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}
```