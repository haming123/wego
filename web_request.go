package wego

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

type WebRequest struct {
	*http.Request
	ctx		*WebContext
}

func (this *WebRequest) Method() string {
	request := this.Request
	return request.Method
}

func (this *WebRequest) URI() string {
	request := this.Request
	return request.RequestURI
}

func (this *WebRequest) URL() string {
	request := this.Request
	return request.URL.Path
}

func (this *WebRequest) GetHeader(key string) string {
	request := this.Request
	return request.Header.Get(key)
}

func (this *WebRequest) Referer() string {
	return this.GetHeader("Referer")
}

func (this *WebRequest) ContentType() string {
	return this.GetHeader("Content-Type")
}

func (this *WebRequest) UserAgent() string {
	return this.GetHeader("User-Agent")
}

func (this *WebRequest) Cookie(name string) (string, error) {
	request := this.Request
	cookie, err := request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (this *WebRequest) Host() string {
	request := this.ctx.Input.Request
	host, _, err := net.SplitHostPort(request.Host)
	if err == nil {
		return host
	}
	return "localhost"
}

func (this *WebRequest) RemoteIP() string {
	request := this.Request
	remote_addr := strings.TrimSpace(request.RemoteAddr)
	ip, _, err := net.SplitHostPort(remote_addr)
	if err != nil {
		return ""
	}
	return ip
}

func (this *WebRequest) getIpFromHeader(header string) string {
	str := this.GetHeader(header)
	if str == "" {
		return ""
	}

	items := strings.Split(str, ",")
	ip_str := strings.TrimSpace(items[0])

	ip, _, err := net.SplitHostPort(ip_str)
	if err != nil {
		return ""
	}

	return ip
}

func (this *WebRequest) ClientIP() string {
	for _, header := range this.ctx.engine.IPHeaders {
		ip := this.getIpFromHeader(header)
		if ip != "" {
			return ip
		}
	}
	return this.RemoteIP()
}
