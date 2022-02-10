package wego

import (
	"net/http"
	"reflect"
)

type BeforeExecer interface {
	BeforeExec(ctx *WebContext)
}

type AfterExecer interface {
	AfterExec(ctx *WebContext)
}

type RouteInfo struct {
	group        *RouteGroup
	route_type 	 nodeType
	pattern      string
	filters      []FilterInfo
	before_func  HandlerFunc
	before_mthd	 BeforeExecer
	after_func   HandlerFunc
	after_mthd   AfterExecer
	handler_type HandlerType
	handler_raw  http.HandlerFunc
	handler_ctx  HandlerFunc
	handler_ctl  reflect.Value
	handler_name string
	func_name    string
}

func (r *RouteInfo) BeforExec(handler HandlerFunc) *RouteInfo {
	r.before_func = handler
	return r
}

func (r *RouteInfo) AfterExec(handler HandlerFunc) *RouteInfo {
	r.after_func = handler
	return r
}

func (r *RouteInfo) GetHandlerName() string {
	return r.handler_name
}
