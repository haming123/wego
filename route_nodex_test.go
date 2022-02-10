package wego

import (
	"fmt"
	"testing"
)

func printGetPathPart(strs ...string) {
	str1 := strs[0]
	name := strs[1]
	str2 := strs[2]
	fmt.Printf("%s,%s,%s\n", str1, name, str2)
}

func TestGetPathPart4Add(t *testing.T) {
	printGetPathPart(splitPath("/hello"))
	printGetPathPart(splitPath("/hello/aaa"))
	printGetPathPart(splitPath("/hello/:name"))
	printGetPathPart(splitPath("/hello/:name/aaa"))
	printGetPathPart(splitPath("/:name"))
	printGetPathPart(splitPath(":name"))
}

func TestGetPathPart4Get(t *testing.T) {
	t.Log(trimValue("aaa"))
	t.Log(trimValue("aaa/bbb"))
	t.Log(trimValue("/"))
}

func TestAddRouteLine1(t *testing.T) {
	var tree TreeNodeX
	var rval RouteInfo
	tree.AddRoute("GET", "/hello/:name", &rval)
}

func TestAddRouteLine2(t *testing.T) {
	var tree TreeNodeX
	var rval RouteInfo
	tree.AddRoute("GET", "/hello/*name/aaa", &rval)
}

func web_home(c *WebContext)  {
	c.Path = "_"
}
func web_index(c *WebContext)  {
	c.Path = "/"
}
func web_hello(c *WebContext)  {
	c.Path = "/hello"
}
func web_hello_dir(c *WebContext)  {
	c.Path = "/hello/"
}
func web_hello_b(c *WebContext)  {
	c.Path = "/hello/b"
}
func web_hello_p(c *WebContext)  {
	c.Path = "/hello/:name"
}
func web_hello_p_p(c *WebContext)  {
	c.Path = "/hello/:name/:age"
}
func web_assets_x(c *WebContext)  {
	c.Path = "/assets/*filepath"
}

func AddRoute4Test(r *TreeNodeX){
	r.AddRoute("GET", "", &RouteInfo{handler_ctx: web_home})
	r.AddRoute("GET", "/", &RouteInfo{handler_ctx: web_index})
	r.AddRoute("GET", "/hello", &RouteInfo{handler_ctx: web_hello})
	r.AddRoute("GET", "/hello/", &RouteInfo{handler_ctx: web_hello_dir})
	r.AddRoute("GET", "/hello/b", &RouteInfo{handler_ctx: web_hello_b})
	r.AddRoute("GET", "/hello/:name/:age", &RouteInfo{handler_ctx: web_hello_p_p})
	r.AddRoute("GET", "/hello/:name", &RouteInfo{handler_ctx: web_hello_p})
	r.AddRoute("GET", "/assets/*filepath", &RouteInfo{handler_ctx: web_assets_x})
}

func TestRouteAddAndGet(t *testing.T) {
	var r TreeNodeX
	AddRoute4Test(&r)
	ctx := new(WebContext)

	path := ""
	ctx.RouteParam.Reset()
	hd := r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		path = "_"
		t.Logf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Errorf("%s: not find\n", path)
	}

	path = "/"
	ctx.RouteParam.Reset()
	hd = r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		t.Logf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Errorf("%s: not find\n", path)
	}

	path = "/hello"
	ctx.RouteParam.Reset()
	hd = r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		t.Logf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Errorf("%s: not find\n", path)
	}

	path = "/hello/"
	ctx.RouteParam.Reset()
	hd = r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		t.Logf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Errorf("%s: not find\n", path)
	}

	path = "/hello/b"
	ctx.RouteParam.Reset()
	hd = r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		t.Logf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Errorf("%s: not find\n", path)
	}

	path = "/hello/lisi"
	ctx.RouteParam.Reset()
	hd = r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		t.Logf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Errorf("%s: not find\n", path)
	}

	path = "/hello/lisi/12"
	ctx.RouteParam.Reset()
	hd = r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		t.Logf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Errorf("%s: not find\n", path)
	}

	path = "/assets/999"
	ctx.RouteParam.Reset()
	hd = r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		t.Logf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Errorf("%s: not find\n", path)
	}

	path = "/assets/999/888"
	ctx.RouteParam.Reset()
	hd = r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		t.Logf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Errorf("%s: not find\n", path)
	}
}

func AddRoute4Test2(r *TreeNodeX) {
	r.AddRoute("GET", "/hello", &RouteInfo{handler_ctx: web_hello})
	r.AddRoute("GET", "/world/", &RouteInfo{handler_ctx: web_hello_dir})
}

func TestGetRouteNotMatch2(t *testing.T) {
	var r TreeNodeX
	AddRoute4Test2(&r)
	ctx := new(WebContext)

	path := "/hello/"
	hd := r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		t.Errorf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Logf("%s: not find\n", path)
	}

	path = "/world"
	hd = r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		t.Errorf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Logf("%s: not find\n", path)
	}

	path = "/world/aaa"
	hd = r.GetRoute("GET", path, &ctx.RouteParam)
	if hd != nil {
		hd.handler_ctx(ctx)
		t.Errorf("%s => %s : %v\n", path, ctx.Path, ctx.RouteParam)
	} else {
		t.Logf("%s: not find\n", path)
	}
}
