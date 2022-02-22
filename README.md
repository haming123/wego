## 简介
wego是一个Go语言编写的高性能的Web框架，可以用来快速开发RESTful服务以及后端服务等各种应用。
wego框架是一个完整的MVC框架，包括路由模块、数据库ORM模块、view模板处理以及Session模块。
wego具有性能高、方便易用，兼容性好，扩展性强等特点，具体特征如下：
1. 基于Radix树开发的路由模块，路由查询性能高。
2. 支持路由组。
3. 为路由参数、Query参数、Form参数的访问提供率方便易于使用的API，并可以将参数映射到Struct。
4. 为JSON、XML和HTML渲染提供了易于使用的API。
5. 支持过滤器中间件，方便您对Web框架进行扩展。
6. 支持BeforeRoute、BeforeExec、AfterExec拦截器，方便您进行身份验证、日志输出。
7. 内置Crash处理机制，wego可以recover一个HTTP请求中的panic，这样可确保您的服务器始终可用。
8. 内置Config模块，方便对应用的参数进行管理。
9. 内置Session模块，您可以选择cookie、redis、memcache、memory缓存引擎存储Session数据。
10. 支持混合路由，固定路由、参数路由、通配符路由可以混合，不冲突。
11. 内置log模块，用于生成应用日志。
12. 采用缓存来管理HTML的Template，既方便输出Html页面，又可以使用缓存提升系统性能。
13. 良好的兼容性，wego支持go原生的func(http.ResponseWriter, *http.Request)路由处理函数，这样您的代码少量修改就可以使用wego了。
14. wego兼容两种编码习惯，可以使用普通函数作为路由处理函数，也可以使用strcut的成员函数作为路由处理函数。

### 安装
go get github.com/haming123/wego

### 使用文档
请点击：[详细文档](http://39.108.252.54:8080)

### 简单http server
创建一个main.go文件，代码如下：
```go
package main
import (
	"wego"
	log "wego/dlog"
)
func main() {
	web, err := wego.NewWeb()
	if err != nil{
		log.Error(err)
		return
	}

	web.GET("/hello", func(c *wego.WebContext) {
		c.WriteText(200, "world")
	})

	err = web.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
```
然后运行它，打开浏览器，输入http://localhost:8080/hello， 就可以看到如下内容：
```
world
```

### 项目结构
wego框没有对项目结构做出限制，这里给出一个MVC框架项目的建议结构：
```
demo
├── app             
│   └── router.go       - 路由配置文件
│   └── templfun.go     - 模板函数文件
├── controllers         - 控制器目录，必要的时候可以继续划分子目录
│   └── controller_xxx.go
├── models              - 模型目录
│   └── model_xxx.go
├── logs                - 日志文件目录，主要保存项目运行过程中产生的日志
│   └── applog_20211203.log
├── static              - 静态资源目录
│   ├── css
│   ├── img
│   └── js
├── utils               - 公共代码目录
│   └── util_xxx.go
├── views               - 视图模板目录
│   └── html_xxx.html
├── app.conf            - 应用配置文件
└── main.go             - 入口文件
```

### 注册参数路由
```go
func main() {
	web, err := wego.NewWeb()
	if err != nil{
		log.Error(err)
		return
	}

	web.PATH("/user/:id", func(c *wego.WebContext) {
		c.WriteTextF(200, "param id=%s", c.RouteParam.GetString("id").Value)
	})

	err = web.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
```

### 注册模糊匹配路由
```go
func main() {
	web, err := wego.NewWeb()
	if err != nil{
		log.Error(err)
		return
	}

	web.PATH("/files/*name", func(c *wego.WebContext) {
		c.WriteTextF(200, "param name=%s", c.RouteParam.GetString("name").Value)
	})

	err = web.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
```

### 注册RESTful路由
```go
func main() {
	web, err := wego.NewWeb()
	if err != nil{
		log.Error(err)
		return
	}

	web.GET("/users/:id", func(c *wego.WebContext) {
		//查询一个用户
	})
	web.POST("/users/:id", func(c *wego.WebContext) {
		//创建一个用户
	})
	web.PUT("/users/:id", func(c *wego.WebContext) {
		//更新用户信息
	})
	web.PATCH("/users/:id", func(c *wego.WebContext) {
		//更新用户的部分信息
	})
	web.DELETE("/users/:id", func(c *wego.WebContext) {
		//删除用户
	})

	err = web.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
```

### 获取参数
在wego中通过c.Param.GetXXX函数来获取请求参数：
```go
func main() {
	web, err := wego.NewWeb()
	if err != nil{
		log.Error(err)
		return
	}

	web.GET("/user", func(c *WebContext) {
		name := c.Param.GetString("name")
		if name.Error != nil {
			t.Error(name.Error)
		}
		age := c.Param.GetInt("age")
		if age.Error != nil {
			t.Error(age.Error)
		}
        c.WriteText(200, name.Value)
	})

	err = web.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
```

### 使用MustXXX便捷函数获取参数
```go
func main() {
	web, err := wego.NewWeb()
	if err != nil{
		log.Error(err)
		return
	}

	web.GET("/user", func(c *WebContext) {
		name := c.Param.MustString("name")
		t.Log(name)
		age := c.Param.MustInt("age")
		t.Log(age)
	})

	err = web.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
```

### 使用ReadJSON获取参数
若POST请求中Body的数据的格式为JSON格式，可以直接使用WebContext的ReadJSON函数来读取：
 ```go
func main() {
	web, err := wego.NewWeb()
	if err != nil{
		log.Error(err)
		return
	}

	web.POST("/user", func(c *WebContext) {
		var user2 User
		err := c.ReadJSON(&user2)
		if err != nil {
			t.Log(err)
		}
		t.Log(user2)
	})

	err = web.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
```

### 输出JSON数据
wego对于JSON的支持非常好，可以让我们非常方便的开发一个基于JSON的API。若要返回JSON请求结果，您可以使用WriteJSON函数：
```go
func writeJson(c *wego.WebContext) {
	var user User
	user.ID = 1
	user.Name = "demo"
	user.Age = 12
	c.WriteJSON(200, user)
}
```

### 输出HTML数据
wego框架的html结果的输出是基于html/template实现的。以下是一个输出html页面的例子：
```go
func writeHtml(c *wego.WebContext) {
	var user User
	user.ID = 1
	user.Name = "demo"
	user.Age = 12
	c.WriteHTML(200, "./views/index.html", user)
}
```

### 使用模板函数
如果您的模板文件中使用了模板函数，需要预先将所需的模板函数进行登记：
```go
func GetUserID(id int64) string {
	return fmt.Sprintf("ID_%d", id)
}

func main() {
	web, err := wego.NewWeb()
	if err != nil{
		log.Error(err)
		return
	}

	wego.AddTemplFunc("GetUserID", GetUserID)
	web.GET("/templfunc", (c *wego.WebContext) {
        var user User
        user.ID = 1
        user.Name = "lisi"
        user.Age = 12
        c.WriteHTML(200, "./views/index.html", user)
     })

	err = web.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
```

### 设置cookie
```go
func setCookie(c *wego.WebContext)  {
	val, err := c.Input.Cookie("demo")
	if err != nil {
		log.Error(err)
	}
	cookie := &http.Cookie{
		Name:     "demo",
		Value:    "test",
		Path:     "/",
		HttpOnly: true,
	}
	c.SetCookie(cookie)
}
```

### 重定向
```go
func main() {
	web, err := wego.NewWeb()
	if err != nil{
		log.Error(err)
		return
	}

	web.GET("/redirect", func(c *wego.WebContext) {
		c.Redirect(302, "/index")
	})

	err = web.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
```

### 错误处理
```go
func main() {
	web, err := wego.NewWeb()
	if err != nil{
		log.Error(err)
		return
	}

	web.GET("/abort", func(c *wego.WebContext) {
		name := c.Param.GetString("name")
		if name.Error != nil {
			c.AbortWithError(500, name.Error)
			return
		}
		c.WriteText(200, "hello " + name.Value)
	})

	err = web.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
```

### 使用配置文件
* 首先定义一个简单的配置文件
```ini
#应用名称
app_name = demo
#mysql数据库的配置参数
mysql = root:rootpwd@tcp(127.0.0.1:3306)/demp?charset=utf8

[server]
#http监听端口
http_port = 8080
```

* 使用InitWeb()函数初始化Web服务器
```go
func main() {
    web, err := wego.InitWeb()
	if err != nil{
		log.Error(err)
		return
	}

	err = web.Run()
	if err != nil {
		log.Error(err)
	}
}
```
说明：调用InitWeb()函数时可以指定一个配置文件，若没有指定配置文件，则使用缺省的配置文件：./app.conf。

### 获取业务参数
调用InitWeb()后wego会自动将系统参数解析到WebEngine.Config中，业务参数则需要用户自己调用配置数据的GetXXX函数来获取。例如：
```go
func main() {
	web, err := wego.InitWeb()
	if err != nil{
		log.Error(err)
		return
	}

	mysql_cnn := web.Config.GetString("mysql")
	if mysql_cnn.Error != nil {
		log.Error(mysql_cnn.Error)
		return
	}
	log.Info(mysql_cnn.Value)

	err = web.Run()
	if err != nil {
		log.Error(err)
	}
}
```

### 使用Session
* 首选定义配置文件：
```ini
#应用名称
app_name = demo

[server]
#http监听端口
http_port = 8080

[session]
#session 是否开启
session_on = true
#session类型：cookie、cache
session_store=cookie
#客户端的cookie的名称
cookie_name = "wego"
#session 过期时间, 单位秒
life_time = 3600
#session数据的hash字符串
hash_key = demohash
```

* 然后在入口函数中加载配置文件
```go
func main() {
	web, err := wego.NewWeb()
	if err != nil{
		log.Error(err)
		return
	}

	web.GET("/login", login)
	web.GET("/index", index)

	err = web.Run(":8080")
	if err != nil {
		log.Error(err)
	}
}
```

* 然后再login处理器函数中保存session数据：
```go
func login(c *wego.WebContext)  {
	c.Session.Set("uid", 1)
	c.Session.Save()
	c.Redirect(302, "/index")
}
```

* 然后index处理器函数中就可以访问session数据了：
```go
func index(c *wego.WebContext)  {
	id , _ := c.Session.GetInt("uid")
	c.WriteTextF(200, "uid=%d", id)
}
```

### 输出日志
```go
package main
import log "wego/dlog"
func main()  {
	log.Debug("This is a Debug Message")
	log.Info("This is a Info Message")
}
//执行后的输出结果为：
//2021/11/30 07:20:06 [D] main.go:31 This is a Debug Message
//2021/11/30 07:20:06 [I] main.go:32 This is a Debug Info
```
