package wego

import (
	"errors"
	"io"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FormParam struct {
	ctx 	*WebContext
	param 	[]ParamItem
	mdata	url.Values
	loaded 	bool
}

func (this *FormParam)Init(ctx *WebContext) {
	this.ctx = ctx
}

func (this *FormParam)Reset() {
	this.loaded = false
	if this.param != nil {
		this.param = this.param[:0]
	}
	this.mdata = nil
}

func (this *FormParam)parseQuery(query string) (err error) {
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}

		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}

		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}

		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		this.param = append(this.param, ParamItem{key, value})
	}

	if len(this.param) > 50 {
		if this.mdata == nil {
			this.mdata = make(url.Values, len(this.param))
		}
		for _, item := range this.param {
			vals, ok := this.mdata[item.Key]
			if ok {
				vals = append(vals, item.Val)
				this.mdata[item.Key] = vals
			} else {
				this.mdata[item.Key] = []string{item.Val}
			}
		}
	}

	return err
}

func (this *FormParam)parsePostForm() error {
	r := this.ctx.Input.Request
	var reader io.Reader = r.Body
	maxFormSize := int64(1<<63 - 1)
	if _, ok := r.Body.(*maxBytesReader); !ok {
		maxFormSize = int64(10 << 20) // 10 MB is a lot of text.
		reader = io.LimitReader(r.Body, maxFormSize+1)
	}
	buff, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if int64(len(buff)) > maxFormSize {
		return errors.New("http: POST too large")
	}
	return this.parseQuery(string(buff))
}

func (this *FormParam)parsePostMultipart(max_size int64) error {
	err := this.ctx.Input.Request.ParseMultipartForm(max_size)
	if err != nil {
		return err
	}
	/*
	for key, vals := range this.ctx.Input.Request.MultipartForm.Value {
		for _, value := range vals {
			this.param = append(this.param, ParamItem{key, value})
		}
	}*/
	this.mdata = this.ctx.Input.Request.MultipartForm.Value
	return nil
}

func (this *FormParam)parsePost() error {
	r := this.ctx.Input.Request
	if r.Method != "POST" && r.Method != "PUT" && r.Method != "PATCH" {
		return nil
	}

	if r.Body == nil {
		return errors.New("missing form body")
	}

	max_size := int64(32 << 20) //32 MB
	if this.ctx.engine.Config.ServerParam.MaxBody > 0 {
		//Config.ServerParam.MaxBody, unit:M
		max_size = this.ctx.engine.Config.ServerParam.MaxBody*1024*1024
	}

	ct := r.Header.Get("Content-Type")
	media, _ := parseMediaType(ct)
	if media == "application/x-www-form-urlencoded" {
		//r.ParseForm(); this.mdata = r.PostForm
		return this.parsePostForm()
	} else if media == "multipart/form-data" {
		this.parsePostMultipart(max_size)
	}

	return nil
}

func (this *FormParam)loadData() {
	if this.loaded == true {
		return
	}
	this.loaded = true

	if this == &this.ctx.QueryParam {
		//this.mdata = this.ctx.Input.Request.URL.Query()
		this.parseQuery(this.ctx.Input.Request.URL.RawQuery)
	} else {
		this.parsePost()
	}
}

func (this *FormParam)SetValue(key string, val string)  {
	this.param = append(this.param, ParamItem{key, val})
}

func (this *FormParam) GetValues(key string) ([]string, bool) {
	if this.loaded == false {
		this.loadData()
	}

	if this.mdata != nil {
		vals, has := this.mdata[key]
		return vals, has
	}

	if this.param == nil {
		return nil, false
	}

	var arr []string
	for i:=0; i < len(this.param); i++{
		if this.param[i].Key == key {
			arr = append(arr, this.param[i].Val)
		}
	}

	return arr, len(arr)>0
}

func (this *FormParam) GetValue(key string) (string, bool) {
	if this.loaded == false {
		this.loadData()
	}

	if this.mdata != nil {
		vals, has := this.mdata[key]
		if has == false {
			return "", false
		}
		return vals[0], true
	}

	if this.param == nil {
		return "", false
	}

	for i:=0; i < len(this.param); i++{
		if this.param[i].Key == key {
			return this.param[i].Val, true
		}
	}
	return "", false
}

func (this *FormParam) GetString(key string, defaultValue ...string) ValidString {
	val, has := this.GetValue(key)
	if has == false && len(defaultValue) > 0 {
		return ValidString{defaultValue[0], errNotFind}
	} else if has == false {
		return ValidString{"", errNotFind}
	}
	return ValidString{val, nil}
}

func (this *FormParam) MustString(key string, defaultValue ...string) string {
	return this.GetString(key, defaultValue...).Value
}

func (this *FormParam) GetBool(key string, defaultValue ...bool) ValidBool {
	val_str, has := this.GetValue(key)
	if has == false && len(defaultValue) > 0 {
		return ValidBool{defaultValue[0], errNotFind}
	} else if has == false {
		return ValidBool{false, errNotFind}
	}

	val, err := strconv.ParseBool(val_str)
	if err != nil {
		return ValidBool{defaultValue[0], err}
	} else {
		return ValidBool{val, nil}
	}
}

func (this *FormParam) MustBool(key string, defaultValue ...bool) bool {
	return this.GetBool(key, defaultValue...).Value
}

func (this *FormParam) GetInt(key string, defaultValue ...int) ValidInt {
	val_str, has := this.GetValue(key)
	if has == false && len(defaultValue) > 0 {
		return ValidInt{defaultValue[0], errNotFind}
	} else if has == false {
		return ValidInt{0, errNotFind}
	}

	val, err := strconv.Atoi(val_str)
	if err != nil {
		return ValidInt{defaultValue[0], err}
	} else {
		return ValidInt{val, nil}
	}
}

func (this *FormParam) MustInt(key string, defaultValue ...int) int {
	return this.GetInt(key, defaultValue...).Value
}

func (this *FormParam) GetInt32(key string, defaultValue ...int32) ValidInt32 {
	val_str, has := this.GetValue(key)
	if has == false && len(defaultValue) > 0 {
		return ValidInt32{defaultValue[0], errNotFind}
	} else if has == false {
		return ValidInt32{0, errNotFind}
	}

	val, err := strconv.ParseInt(val_str, 10, 32)
	if err != nil {
		return ValidInt32{defaultValue[0], err}
	} else {
		return ValidInt32{int32(val), nil}
	}
}

func (this *FormParam) MustInt32(key string, defaultValue ...int32) int32 {
	return this.GetInt32(key, defaultValue...).Value
}

func (this *FormParam) GetInt64(key string, defaultValue ...int64) ValidInt64 {
	val_str, has := this.GetValue(key)
	if has == false && len(defaultValue) > 0 {
		return ValidInt64{defaultValue[0], errNotFind}
	} else if has == false {
		return ValidInt64{0, errNotFind}
	}

	val, err := strconv.ParseInt(val_str, 10, 64)
	if err != nil {
		return ValidInt64{defaultValue[0], err}
	} else {
		return ValidInt64{val, nil}
	}
}

func (this *FormParam) MustInt64(key string, defaultValue ...int64) int64 {
	return this.GetInt64(key, defaultValue...).Value
}

func (this *FormParam) GetFloat(key string, defaultValue ...float64) ValidFloat {
	val_str, has := this.GetValue(key)
	if has == false && len(defaultValue) > 0 {
		return ValidFloat{defaultValue[0], errNotFind}
	} else if has == false {
		return ValidFloat{0, errNotFind}
	}

	val, err := strconv.ParseFloat(val_str, 64)
	if err != nil {
		return ValidFloat{defaultValue[0], err}
	} else {
		return ValidFloat{val, nil}
	}
}

func (this *FormParam) MustFloat(key string, defaultValue ...float64) float64 {
	return this.GetFloat(key, defaultValue...).Value
}

func (this *FormParam) GetTime(key string, format string, defaultValue ...time.Time) ValidTime {
	val_str, has := this.GetValue(key)
	if has == false && len(defaultValue) > 0 {
		return ValidTime{defaultValue[0], errNotFind}
	} else if has == false {
		return ValidTime{time.Time{}, errNotFind}
	}

	val, err := time.Parse(format, val_str)
	if err != nil {
		return ValidTime{defaultValue[0], err}
	} else {
		return ValidTime{val, nil}
	}
}

func (this *FormParam) MustTime(key string, format string, defaultValue ...time.Time) time.Time {
	return this.GetTime(key, format, defaultValue...).Value
}

func (this *FormParam) GetStruct(ptr interface{}) error {
	if ptr == nil {
		return errors.New("ptr must be *Struct")
	}
	//ent_ptr必须是一个指针
	v_ent := reflect.ValueOf(ptr)
	if v_ent.Kind() != reflect.Ptr {
		return errors.New("ptr must be *Struct")
	}
	//ent_ptr必须是一个结构体指针
	v_ent = reflect.Indirect(v_ent)
	if v_ent.Kind() != reflect.Struct {
		return errors.New("ptr must be *Struct")
	}

	err := bindStruct(v_ent, this)
	if err != nil {
		return err
	}
	return nil
}
