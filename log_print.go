package wego

import (
	"fmt"
	"net/http"
	"time"
)

const (
	green  = "\033[97;42m"
	white  = "\033[90;47m"
	yellow = "\033[90;43m"
	red    = "\033[97;41m"
	blue   = "\033[97;44m"
	reset  = "\033[0m"
)

func GetStatusPrintColor(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusBadRequest:
		return green
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	case code >= http.StatusInternalServerError && code < http.StatusNetworkAuthenticationRequired:
		return red
	default:
		return reset
	}
}

func printReqInfo(c *WebContext) {
	method_color := " => "
	reset_color := ""
	if debug_log.show_color == true {
		method_color = blue
		reset_color = reset
	}
	method_str := fmt.Sprintf("%s%s%s", method_color, c.Input.Request.Method, reset_color)
	path := c.Input.Request.RequestURI
	debug_log.Output(method_str, path)
}

func printExeInfo(c *WebContext) {
	status_color := ""
	reset_color := ""
	if debug_log.show_color == true {
		status_color = GetStatusPrintColor(c.Output.StatusCode)
		reset_color = reset
	}
	status_str := fmt.Sprintf("%s%d%s", status_color, c.Output.StatusCode, reset_color)
	func_name := "NotFind"
	if c.Route != nil {
		func_name = c.Route.func_name
	}
	exe_time := time.Now().Sub(c.Start)
	msg := fmt.Sprintf("%s %0.3fms", func_name, float64(exe_time.Nanoseconds())/float64(1e6))
	if c.state.Error != nil {
		msg += " err:" + c.state.Error.Error()
	}
	debug_log.Output(status_str, msg)
}
