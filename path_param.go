package wego

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

type PathParam struct {
	items	[]ParamItem
}

func (this *PathParam)Init() {
	if this.items == nil {
		this.items = make([]ParamItem, 5)
		this.items = this.items[:0]
	}
	this.items = this.items[:0]
}

func (this *PathParam)Reset() {
	if this.items != nil {
		this.items = this.items[:0]
	}
}

func (this *PathParam)SetValue(key string, val string)  {
	this.items = append(this.items, ParamItem{key, val})
}

func (this *PathParam) GetValues(key string) ([]string, bool) {
	if this.items == nil {
		return nil, false
	}
	val, ok := this.GetValue(key)
	if !ok {
		return nil, false
	}
	return []string{val}, true
}

func (this *PathParam) GetValue(key string) (string, bool) {
	if this.items == nil {
		return "", false
	}
	for i:=0; i < len(this.items); i++{
		if this.items[i].Key == key {
			return this.items[i].Val, true
		}
	}
	return "", false
}

func (this *PathParam) GetString(key string, defaultValue ...string) ValidString {
	val, has := this.GetValue(key)
	if has == false && len(defaultValue) > 0 {
		return ValidString{defaultValue[0], errNotFind}
	} else if has == false {
		return ValidString{"", errNotFind}
	}
	return ValidString{val, nil}
}

func (this *PathParam) MustString(key string, defaultValue ...string) string {
	return this.GetString(key, defaultValue...).Value
}

func (this *PathParam) GetBool(key string, defaultValue ...bool) ValidBool {
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

func (this *PathParam) MustBool(key string, defaultValue ...bool) bool {
	return this.GetBool(key, defaultValue...).Value
}

func (this *PathParam) GetInt(key string, defaultValue ...int) ValidInt {
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

func (this *PathParam) MustInt(key string, defaultValue ...int) int {
	return this.GetInt(key, defaultValue...).Value
}

func (this *PathParam) GetInt32(key string, defaultValue ...int32) ValidInt32 {
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

func (this *PathParam) MustInt32(key string, defaultValue ...int32) int32 {
	return this.GetInt32(key, defaultValue...).Value
}

func (this *PathParam) GetInt64(key string, defaultValue ...int64) ValidInt64 {
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

func (this *PathParam) MustInt64(key string, defaultValue ...int64) int64 {
	return this.GetInt64(key, defaultValue...).Value
}

func (this *PathParam) GetFloat(key string, defaultValue ...float64) ValidFloat {
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

func (this *PathParam) MustFloat(key string, defaultValue ...float64) float64 {
	return this.GetFloat(key, defaultValue...).Value
}

func (this *PathParam) GetTime(key string, format string, defaultValue ...time.Time) ValidTime {
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

func (this *PathParam) MustTime(key string, format string, defaultValue ...time.Time) time.Time {
	return this.GetTime(key, format, defaultValue...).Value
}

func (this *PathParam) GetStruct(ptr interface{}) error {
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
