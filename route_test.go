package wego

import (
	"net/http"
	"testing"
)

func TestRoutePattern(t *testing.T) {
	web, err := NewWeb()
	if err != nil{
		t.Error(err)
		return
	}

	web.PATH("/static", func(c *WebContext) {
		c.WriteText(200, "this is a static route")
	})
	web.PATH("/user/:id", func(c *WebContext) {
		c.WriteTextF(200, "param id=%s", c.RouteParam.GetString("id").Value)
	})
	web.PATH("/files/*name", func(c *WebContext) {
		c.WriteTextF(200, "param name=%s", c.RouteParam.GetString("name").Value)
	})
}

func TestRouteRestful(t *testing.T) {
	web, err := NewWeb()
	if err != nil{
		t.Error(err)
		return
	}

	web.GET("/users/:id", func(c *WebContext) {
		//查询一个用户
	})
	web.POST("/users/:id", func(c *WebContext) {
		//创建一个用户
	})
	web.PUT("/users/:id", func(c *WebContext) {
		//更新用户信息
	})
	web.PATCH("/users/:id", func(c *WebContext) {
		//更新用户的部分信息
	})
	web.DELETE("/user/666", func(c *WebContext) {
		//删除用户
	})
}

func handlerWegoFunc(c *WebContext)  {
	c.WriteText(200, "hello world")
}
func handlerGoHandler(w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte("hello world"))
}
type DemoController struct {
}
func (this *DemoController) WriteText(ctx *WebContext) {
	ctx.WriteText(200, "hello world")
}
func TestRouteHandler(t *testing.T) {
	web, err := NewWeb()
	if err != nil{
		t.Error(err)
		return
	}
	web.GET("/wego_func", handlerWegoFunc)
	web.GET("/wego_method", (*DemoController).WriteText)
	web.GET("/go_hander", handlerGoHandler)
}

func TestRouteGroup(t *testing.T) {
	web, err := NewWeb()
	SetDebugLogLevel(5)
	if err != nil{
		t.Error(err)
		return
	}

	// 创建v1组
	v1 := web.NewGroup("/v1")
	{
		// 在v1这个分组下，注册路由
		v1.POST("/login", func(c *WebContext){})
		v1.POST("/list", func(c *WebContext){})
		v1.POST("/info", func(c *WebContext){})
	}

	// 创建v2组
	v2 := web.NewGroup("/v2")
	{
		// 在v2这个分组下，注册路由
		v2.POST("/login", func(c *WebContext){})
		v2.POST("/list", func(c *WebContext){})
		v2.POST("/info", func(c *WebContext){})
	}
}

