package wego

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"unsafe"
)

const ContentTypeText = "text/plain; charset=utf-8"
const ContentTypeHTML = "text/html; charset=utf-8"
const ContentTypeJSON = "application/json; charset=utf-8"
const ContentTypeXML = "application/xml; charset=utf-8"

type WebStatus int

const (
	STATUS_BEG 	WebStatus = iota
	STATUS_END
)

type WebState struct {
	Status 		int
	Error  		error
}

func (this *WebState) Reset() {
	this.Status = 0
	this.Error = nil
}

func (this *WebState) Set(code int, err error) {
	this.Status = code
	this.Error = err
}

type WebContext struct {
	Config     	*WebConfig
	engine     	*WebEngine
	Input      	WebRequest
	Output     	WebResponse
	Route      	*RouteInfo
	Path       	string
	Param      	WebParam
	RouteParam 	PathParam
	QueryParam 	FormParam
	FormParam  	FormParam
	Data       	ContextData
	Session    	SessionInfo
	Start      	time.Time
	filters    	[]FilterInfo
	state 		WebState
}

func (c *WebContext) reset() {
	c.Input.Request = nil
	c.Route = nil
	c.Path = ""
	c.Output.reset()
	c.RouteParam.Reset()
	c.QueryParam.Reset()
	c.FormParam.Reset()
	c.Data.reset()
	c.Session.Reset()
	c.filters = nil
	c.state.Reset()
}

func newContext() *WebContext {
	ctx := &WebContext{ }
	ctx.Input.ctx = ctx
	ctx.Param.ctx = ctx
	ctx.Session.ctx = ctx
	ctx.RouteParam.Init()
	ctx.QueryParam.Init(ctx)
	ctx.FormParam.Init(ctx)
	return ctx
}

func (c *WebContext) UseGzip(flag bool, min_size ...int64) *WebContext {
	c.Output.gzip_flag = flag
	if len(min_size) > 0 {
		c.Output.gzip_size = min_size[0]
	}
	return c
}

func (c *WebContext) Ended() bool {
	return c.state.Status > 0
}

func (c *WebContext) Next() {
	if len(c.filters) < 1 {
		return
	}
	debug_log.Debug("call filter: " + c.filters[0].name)
	handler := c.filters[0].filter
	c.filters = c.filters[1:]
	handler(c)
}

func (c *WebContext) Abort(code int) {
	c.state.Set(code, errors.New("Abort"))
	c.WriteNoContent(code, "")
}

func (c *WebContext) Abort401() {
	c.state.Set(401, errors.New("Abort401"))
	if c.engine.hanlder_401 != nil {
		c.engine.hanlder_401(c)
	} else {
		c.Abort(401)
	}
}

func (c *WebContext) Abort500() {
	c.state.Set(500, errors.New("Abort500"))
	if c.engine.hanlder_500 != nil {
		c.engine.hanlder_500(c)
	} else {
		c.Abort(500)
	}
}

func (c *WebContext) AbortWithText(code int, value string) {
	c.state.Set(code, errors.New("AbortWithText"))
	c.WriteText(code, value)
}

func (c *WebContext) AbortWithTextF(code int, format string, values ...interface{}) {
	c.state.Set(code, errors.New("AbortWithText"))
	c.WriteTextF(code, format, values...)
}

func (c *WebContext) AbortWithError(code int, err error) {
	c.state.Set(code, errors.New("AbortWithError"))
	c.WriteText(code, err.Error())
}

func (c *WebContext) AbortWithJson(code int, obj interface{}) {
	c.state.Set(code, errors.New("AbortWithJson"))
	c.WriteJSON(code, obj)
}

func (c *WebContext) AbortWithXml(code int, obj interface{}) {
	c.state.Set(code, errors.New("AbortWithXml"))
	c.WriteXML(code, obj)
}

func (c *WebContext) AbortWithHtml(code int, filenames string, data interface{}) {
	c.state.Set(code, errors.New("AbortWithHtml"))
	c.WriteHTML(code, filenames, data)
}

func (c *WebContext) AbortWithTemplate(code int, templ *template.Template, data interface{}) {
	c.state.Set(code, errors.New("AbortWithTemplate"))
	c.WriteTemplate(code, templ, data)
}

func (c *WebContext) ReadBody() ([]byte, error) {
	return ioutil.ReadAll(c.Input.Request.Body)
}

func (c *WebContext) ReadJSON(obj interface{}) error {
	decoder := json.NewDecoder(c.Input.Request.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}

func (c *WebContext) ReadXML(obj interface{}) error {
	decoder := xml.NewDecoder(c.Input.Request.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}

func (c *WebContext) GetFile(name string) (*multipart.FileHeader, error) {
	f, fh, err := c.Input.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

func (c *WebContext) GetFiles(key string) ([]*multipart.FileHeader, error) {
	if files, ok := c.Input.Request.MultipartForm.File[key]; ok {
		return files, nil
	}
	return nil, http.ErrMissingFile
}

func (c *WebContext) SaveToFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (c *WebContext) Status(code int) {
	c.Output.SetStatus(code)
}

func (c *WebContext) SetHeader(key string, value string) {
	c.Output.ResponseWriter.Header().Set(key, value)
}

/*
Domain：定义Cookie的生效作用域，只有当域名和路径同时满足的时候，浏览器才会将Cookie发送给Server。
Expires：过期时间，绝对时间点。
Max-Age：收到报文后多久的过期时间，单位是秒。设置为0时立刻失效。
HttpOnly：如设为True，JavaScript这些脚本无法访问，即有效减少 XSS 攻击
Secure：如设置secure为true，浏览器只会在HTTPS和SSL等安全协议中传输Cookie。
SameSite：可以防范 CSRF 攻击。
	SameSite=None：无论是否跨站都会发送 Cookie，
	SameSite=Lax：允许部分第三方请求携带 Cookie，
	SameSite=Strict：仅允许同站请求携带 Cookie，即当前网页 URL 与请求目标 URL 完全一致
 */
func (c *WebContext) SetCookie(ck *http.Cookie) {
	http.SetCookie(&c.Output, ck)
}

func (c *WebContext) Redirect(code int, localurl string) {
	c.state.Set(code, nil)
	http.Redirect(&c.Output, c.Input.Request, localurl, code)
}

func (c *WebContext) fail(code int, err error) {
	if err == nil {
		err = errors.New("Fail")
	}
	c.state.Set(code, err)
	c.SetHeader("Content-Type", ContentTypeText)
	c.Output.SetStatus(code)
	c.Output.Write(StringToBytes(err.Error()))
}

func (c *WebContext) Write(code int, contentType string, data []byte) {
	if len(contentType) > 0 {
		c.SetHeader("Content-Type", contentType)
	}
	c.Status(code)
	c.Output.Write(data)
}

func (c *WebContext) WriteNoContent(code int, contentType string)  {
	if len(contentType) > 0 {
		c.SetHeader("Content-Type", contentType)
	}
	c.Output.WriteHeader(code)
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func (c *WebContext) WriteText(code int, data string) {
	c.SetHeader("Content-Type", ContentTypeText)
	c.Status(code)
	c.Output.Write(StringToBytes(data))
}

func (c *WebContext) WriteTextF(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", ContentTypeText)
	c.Status(code)
	fmt.Fprintf(&c.Output, format, values...)
}

func (c *WebContext) WriteTextBytes(code int, data []byte) {
	c.SetHeader("Content-Type", ContentTypeText)
	c.Status(code)
	c.Output.Write(data)
}

func (c *WebContext) WriteJSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", ContentTypeJSON)
	c.Status(code)

	err := json.NewEncoder(&c.Output).Encode(obj)
	if err != nil {
		c.fail(500, err)
	}
}

//美化JSON的输出，可以使用json.MarshalIndent方法对obj进行编码
func (c *WebContext) WriteIndentedJSON(code int, obj interface{}) {
	data, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		c.fail(500, err)
	}
	c.Write(code, ContentTypeJSON, data)
}

//输出JSONP，使用callback参数名作为接收回调函数的名称
func (c *WebContext) WriteJSONP(code int, obj interface{}) {
	callback := c.QueryParam.GetString("callback").Value
	if callback == "" {
		c.WriteJSON(code, obj)
		return
	}

	data, err := json.Marshal(obj)
	if err != nil {
		c.fail(500, err)
	}

	callback = template.JSEscapeString(callback)
	buff := bytes.NewBufferString(callback)
	buff.WriteString("(")
	buff.Write(data)
	buff.WriteString(");\r\n")

	data = buff.Bytes()
	c.Write(code, ContentTypeJSON, data)
}

//防JSON劫持缺省的前缀是while(1);
func (c *WebContext) WriteSecureJSON(code int, obj interface{}) {
	data, err := json.Marshal(obj)
	if err != nil {
		c.fail(500, err)
	}

	if bytes.HasPrefix(data, []byte("[")) == false {
		c.Write(code, ContentTypeJSON,  data)
		return
	}

	if bytes.HasSuffix(data, []byte("]")) == false {
		c.Write(code, ContentTypeJSON,  data)
		return
	}

	prefix := c.engine.Config.JsonPrefix
	buff := bytes.NewBufferString(prefix)
	buff.Write(data)

	data = buff.Bytes()
	c.Write(code, ContentTypeJSON, data)
}

func (c *WebContext) WriteXML(code int, obj interface{}) {
	c.SetHeader("Content-Type", ContentTypeXML)
	c.Status(code)

	err := xml.NewEncoder(&c.Output).Encode(obj)
	if err != nil {
		c.fail(500, err)
	}
}

func (c *WebContext) WriteIndentedXML(code int, obj interface{}) {
	data, err := xml.MarshalIndent(obj, "", "    ")
	if err != nil {
		c.fail(500, err)
	}
	c.Write(code, ContentTypeXML, data)
}

func (c *WebContext) WriteHTML(code int, filename string, data interface{}) {
	c.SetHeader("Content-Type", ContentTypeHTML)
	c.Status(code)

	t, err := c.engine.Templates.getTemplate(filename)
	if err != nil {
		c.fail(500, err)
		return
	}

	err = t.Execute(&c.Output, data)
	if err != nil {
		c.fail(500, err)
		return
	}
}

func (c *WebContext) WriteLayoutHTML(code int, layout_file string, content_file string, data interface{}) {
	c.SetHeader("Content-Type", ContentTypeHTML)
	c.Status(code)

	t, err := c.engine.Templates.getTemplate(layout_file, content_file)
	if err != nil {
		c.fail(500, err)
		return
	}

	err = t.Execute(&c.Output, data)
	if err != nil {
		c.fail(500, err)
		return
	}
}

func (c *WebContext) WriteHTMLS(code int, filenames []string, data interface{}) {
	c.SetHeader("Content-Type", ContentTypeHTML)
	c.Status(code)

	t, err := c.engine.Templates.getTemplate(filenames...)
	if err != nil {
		c.fail(500, err)
		return
	}

	err = t.Execute(&c.Output, data)
	if err != nil {
		c.fail(500, err)
		return
	}
}

func (c *WebContext) WriteTemplate(code int, templ *template.Template, data interface{}) {
	c.SetHeader("Content-Type", ContentTypeHTML)
	c.Status(code)

	err := templ.Execute(&c.Output, data)
	if err != nil {
		c.fail(500, err)
		return
	}
}

func (c *WebContext) WriteHtmlBytes(code int, data []byte) {
	c.SetHeader("Content-Type", ContentTypeHTML)
	c.Status(code)
	c.Output.Write(data)
}

func (c *WebContext) WriteFile(filePath string, fileName ...string) {
	var save_name string
	if len(fileName) > 0 && fileName[0] != "" {
		save_name = fileName[0]
	} else {
		save_name = filepath.Base(filePath)
	}
	c.Output.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", save_name))
	http.ServeFile(&c.Output, c.Input.Request, filePath)
}
