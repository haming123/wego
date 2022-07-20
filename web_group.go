package wego

import (
	"net/http"
	"path"
	"reflect"
	"runtime"
)

type FilterInfo struct {
	name   string
	filter HandlerFunc
}

type RouteGroup struct {
	engine      *WebEngine
	parent      *RouteGroup
	path        string
	filters     []FilterInfo
	before_exec HandlerFunc
	after_exec  HandlerFunc
}

func (group *RouteGroup) InitRoot(engine *WebEngine) {
	group.engine = engine
	group.parent = nil
	group.path = ""
}

func (group *RouteGroup) NewGroup(path string) *RouteGroup {
	child := &RouteGroup{}
	child.engine = group.engine
	child.parent = group
	child.path = group.path + path
	return child
}

func (group *RouteGroup) BeforExec(handler HandlerFunc) {
	group.before_exec = handler
}

func (group *RouteGroup) AfterExec(handler HandlerFunc) {
	group.after_exec = handler
}

func (group *RouteGroup) AddFilter(name string, handler HandlerFunc) bool {
	for _, item := range group.filters {
		if item.name == name {
			return false
		}
	}
	item := FilterInfo{name: name, filter: handler}
	group.filters = append(group.filters, item)
	return true
}

//获取Handler相关的过滤器
func (group *RouteGroup) getHandlerFilter() []FilterInfo {
	for g := group; g != nil; g = g.parent {
		if len(g.filters) > 0 {
			return g.filters
		}
	}
	return nil
}

//获取全路径过滤器
func (group *RouteGroup) getFullFilter() []FilterInfo {
	var arr_filter []FilterInfo
	map_filter := make(map[string]bool)
	for g := group; g != nil; g = g.parent {
		for i := len(g.filters) - 1; i >= 0; i-- {
			name := g.filters[i].name
			_, has := map_filter[name]
			if has {
				continue
			}
			map_filter[name] = true
			arr_filter = append(arr_filter, g.filters[i])
		}
	}

	//数组反转
	length := len(arr_filter)
	for i := 0; i < length/2; i++ {
		temp := arr_filter[length-1-i]
		arr_filter[length-1-i] = arr_filter[i]
		arr_filter[i] = temp
	}
	return arr_filter
}

func (group *RouteGroup) getBeforExec() HandlerFunc {
	for g := group; g != nil; g = g.parent {
		if g.before_exec != nil {
			return g.before_exec
		}
	}
	return nil
}

func (group *RouteGroup) getAfterExec() HandlerFunc {
	for g := group; g != nil; g = g.parent {
		if g.after_exec != nil {
			return g.after_exec
		}
	}
	return nil
}

func GetNameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func (group *RouteGroup) createHandlerRoute4Raw(handler http.HandlerFunc) *RouteInfo {
	rinfo := &RouteInfo{}
	rinfo.group = group
	rinfo.handler_type = FT_RAW_HANDLER
	rinfo.handler_raw = handler
	rinfo.handler_name = GetNameOfFunction(handler)
	rinfo.func_name = getFuncShortName(rinfo.handler_name)
	//rinfo.filters = group.getHandlerFilter()
	//rinfo.before_exec = group.getBeforExec()
	//rinfo.after_exec = group.getAfterExec()
	return rinfo
}

func (group *RouteGroup) createHandlerRoute4Ctx(handler HandlerFunc) *RouteInfo {
	rinfo := &RouteInfo{}
	rinfo.group = group
	rinfo.handler_type = FT_CTX_HANDLER
	rinfo.handler_ctx = handler
	rinfo.handler_name = GetNameOfFunction(handler)
	rinfo.func_name = getFuncShortName(rinfo.handler_name)
	//rinfo.filters = group.getHandlerFilter()
	//rinfo.before_exec = group.getBeforExec()
	//rinfo.after_exec = group.getAfterExec()
	return rinfo
}

func (group *RouteGroup) createFuncRoute4Ctl(handler interface{}) *RouteInfo {
	rinfo := &RouteInfo{}
	rinfo.group = group
	rinfo.handler_type = FT_CTL_HANDLER
	ctl_type, ctl_func := getControlFuncInfo(handler)
	rinfo.handler_ctl = reflect.New(ctl_type)
	rinfo.handler_name = ctl_func
	rinfo.func_name = getFuncShortName(rinfo.handler_name)
	if mthd, ok := rinfo.handler_ctl.Interface().(BeforeExecer); ok {
		rinfo.before_mthd = mthd
	}
	if mthd, ok := rinfo.handler_ctl.Interface().(AfterExecer); ok {
		rinfo.after_mthd = mthd
	}
	//rinfo.filters = group.getHandlerFilter()
	//rinfo.before_exec = group.getBeforExec()
	//rinfo.after_exec = group.getAfterExec()
	return rinfo
}

func (group *RouteGroup) addRoute(method string, pattern string, handler interface{}) *RouteInfo {
	var rinfo *RouteInfo
	if v, ok := handler.(HandlerFunc); ok {
		rinfo = group.createHandlerRoute4Ctx(v)
	} else if v, ok := handler.(func(*WebContext)); ok {
		rinfo = group.createHandlerRoute4Ctx(v)
	} else if v, ok := handler.(http.HandlerFunc); ok {
		rinfo = group.createHandlerRoute4Raw(v)
	} else if v, ok := handler.(func(http.ResponseWriter, *http.Request)); ok {
		rinfo = group.createHandlerRoute4Raw(v)
	} else {
		rinfo = group.createFuncRoute4Ctl(handler)
	}
	pattern = group.path + pattern
	rinfo.pattern = pattern
	debug_log.Debugf("%-6s %-20s --> %s\n", method, pattern, rinfo.handler_name)
	group.engine.Route.addRoute(method, pattern, rinfo)
	return rinfo
}

func (group *RouteGroup) PATH(pattern string, handler interface{}) *RouteInfo {
	return group.addRoute(MethodPath, pattern, handler)
}

func (group *RouteGroup) GET(pattern string, handler interface{}) *RouteInfo {
	return group.addRoute(MethodGet, pattern, handler)
}

func (group *RouteGroup) POST(pattern string, handler interface{}) *RouteInfo {
	return group.addRoute(MethodPost, pattern, handler)
}

func (group *RouteGroup) DELETE(pattern string, handler interface{}) *RouteInfo {
	return group.addRoute(MethodDelete, pattern, handler)
}

func (group *RouteGroup) PATCH(pattern string, handler interface{}) *RouteInfo {
	return group.addRoute(MethodPatch, pattern, handler)
}

func (group *RouteGroup) PUT(pattern string, handler interface{}) *RouteInfo {
	return group.addRoute(MethodPut, pattern, handler)
}

func (group *RouteGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.path, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *WebContext) {
		file := c.RouteParam.GetString("filepath").Value
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(&c.Output, c.Input.Request)
	}
}

func (group *RouteGroup) StaticFile(relativePath, filepath string) *RouteInfo {
	handler := func(c *WebContext) {
		c.WriteFile(filepath)
	}
	return group.GET(relativePath, handler)
}

func (group *RouteGroup) StaticPath(relativePath string, root string) *RouteInfo {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	return group.GET(urlPattern, handler)
}
