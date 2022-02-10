package wego

import (
	"context"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type HandlerType int
const (
	FT_CTL_HANDLER HandlerType = iota
	FT_CTX_HANDLER
	FT_RAW_HANDLER
)

type HttpHandler func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*WebContext)

var RouteCtxKey string = "RouteContext"
func GetWebContext(r *http.Request) *WebContext {
	rctx, _ := r.Context().Value(RouteCtxKey).(*WebContext)
	return rctx
}

func callHandler4Raw(c *WebContext) {
	r := c.Input.Request
	r = r.WithContext(context.WithValue(r.Context(), RouteCtxKey, c))
	c.Route.handler_raw(&c.Output, r)
}

func callHandler4Ctx(c *WebContext) {
	c.Route.handler_ctx(c)
}

func callHandler4Ctl(c *WebContext) {
	v_ctl := c.Route.handler_ctl
	method := v_ctl.MethodByName(c.Route.func_name)
	param := []reflect.Value{reflect.ValueOf(c)}
	method.Call(param)
}

func callHandler(c *WebContext) {
	if c.Route.handler_type == FT_RAW_HANDLER {
		callHandler4Raw(c)
	} else if c.Route.handler_type == FT_CTX_HANDLER {
		callHandler4Ctx(c)
	} else {
		callHandler4Ctl(c)
	}
}

func getFuncShortName(full_name string) string {
	arr_name := strings.Split(full_name, ".")
	if len(arr_name) < 1 {
		return ""
	}
	path_level := len(arr_name)
	return arr_name[path_level-1]
}

func getControlFuncInfo(ctl_func interface{}) (reflect.Type, string) {
	v_func := reflect.ValueOf(ctl_func)
	t_func := v_func.Type()
	if t_func.Kind() != reflect.Func {
		panic("ctl_func muse be a contol's method")
	}

	func_info := runtime.FuncForPC(v_func.Pointer())
	if func_info == nil {
		panic("call runtime.FuncForPC failed")
	}
	func_name := func_info.Name()

	num_param := t_func.NumIn();
	if num_param != 2 {
		panic("invalid number of param in")
	}

	ctl_type := t_func.In(0)
	if ctl_type.Kind() == reflect.Ptr {
		ctl_type = ctl_type.Elem()
	}

	return ctl_type, func_name
}

