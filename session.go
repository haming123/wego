package wego

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

type SessionData map[string]string
type SessionInfo struct {
	ctx  *WebContext
	Sid  string
	Data SessionData
	Edit bool
}

func (this *SessionInfo) Reset() {
	this.Sid = ""
	this.Edit = false
	this.Data = nil
}

func (this *SessionInfo) Set(key string, value interface{}) {
	strval, _ := json.Marshal(value)
	if this.Data == nil {
		this.Data = map[string]string{key:string(strval)}
		this.Edit = true
		return
	}
	this.Data[key] = string(strval)
	this.Edit = true
}

func (this *SessionInfo) Get (key string) (string, bool) {
	if this.Data == nil {
		return "", false
	}

	if v, ok := this.Data[key]; ok {
		return v, true
	}

	return "", false
}

func (this *SessionInfo) GetString (key string) (string, error) {
	str_data, has := this.Get(key)
	if has == false {
		return "", errNotFind
	}

	var val string
	err := json.Unmarshal([]byte(str_data), &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func (this *SessionInfo) GetBool (key string) (bool, error) {
	str_data, has := this.Get(key)
	if has == false {
		return false, errNotFind
	}

	var val bool
	err := json.Unmarshal([]byte(str_data), &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func (this *SessionInfo) GetInt (key string) (int64, error) {
	str_data, has := this.Get(key)
	if has == false {
		return 0, errNotFind
	}

	var val int64
	err := json.Unmarshal([]byte(str_data), &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func (this *SessionInfo) GetFloat (key string) (float64, error) {
	str_data, has := this.Get(key)
	if has == false {
		return 0, errNotFind
	}
	return  strconv.ParseFloat(str_data,  64)
}

func (this *SessionInfo) GetTime(key string) (time.Time, error) {
	str_data, has := this.Get(key)
	if has == false {
		return time.Time{}, errNotFind
	}

	var val time.Time
	err := json.Unmarshal([]byte(str_data), &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func (this *SessionInfo) GetStuct (key string, ptr interface{}) error {
	if ptr == nil {
		return errors.New("ptr must be *Struct")
	}

	v_ent := reflect.ValueOf(ptr)
	if v_ent.Kind() != reflect.Ptr {
		return errors.New("ptr must be *Struct")
	}

	str_data, has := this.Get(key)
	if has == false {
		return errNotFind
	}

	err := json.Unmarshal([]byte(str_data), ptr)
	if err != nil {
		return err
	}
	return nil
}

func (this *SessionInfo) Delete(key string) {
	if this.Data == nil {
		return
	}

	delete(this.Data, key)
	this.Edit = true
}

func (this *SessionInfo) Clear() {
	this.Data = nil
	this.Edit = true
}

func (this *SessionInfo) Save() error {
	out := this.ctx.Output.ResponseWriter
	ses := &this.ctx.Session
	if ses.Sid == "" {
		return nil
	}

	//json编码
	data, err := json.Marshal(ses.Data)
	if err != nil {
		return err
	}

	sid := ses.Sid
	http_ctx := context.WithValue(context.Background(), "http.ResponseWriter", out)
	err = this.ctx.engine.Session.SaveData(http_ctx, sid, data)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     this.ctx.engine.Session.cookieName,
		Value:    url.QueryEscape(sid),
		Path:     "/",
		HttpOnly: true,
		MaxAge: int(this.ctx.engine.Session.maxAge),
	}
	http.SetCookie(out, cookie)
	return nil
}

func (this *SessionInfo)Read() error {
	cookieValue := ""
	req := this.ctx.Input.Request
	cookie, _ := req.Cookie(this.ctx.engine.Session.cookieName)
	if cookie != nil {
		cookieValue = cookie.Value
	}

	sid, err := url.QueryUnescape(cookieValue)
	if err != nil {
		return  err
	}

	if sid == "" {
		sid, _ = this.ctx.engine.Session.CreateSid()
		this.ctx.Session.Sid = sid
		this.ctx.Session.Edit = true
		return nil
	}

	this.ctx.Session.Sid = sid
	http_ctx := context.WithValue(context.Background(), "http.Request", req)
	data, err := this.ctx.engine.Session.ReadData(http_ctx, sid)
	if err != nil {
		return err
	}

	ses := &this.ctx.Session
	err = json.Unmarshal(data, &ses.Data)
	if err != nil {
		return err
	}

	return nil
}
