package wego

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

type WebParam struct {
	ctx 	*WebContext
}

func (this *WebParam) GetValues(key string) ([]string, bool) {
	var arr []string
	vals, has := this.ctx.RouteParam.GetValues(key)
	if has == true {
		if arr == nil {
			arr = vals
		} else {
			arr = append(arr, vals...)
		}
	}
	vals, has = this.ctx.QueryParam.GetValues(key)
	if has == true {
		if arr == nil {
			arr = vals
		} else {
			arr = append(arr, vals...)
		}
	}
	vals, has = this.ctx.FormParam.GetValues(key)
	if has == true {
		if arr == nil {
			arr = vals
		} else {
			arr = append(arr, vals...)
		}
	}
	return arr, len(arr)>0
}

func (this *WebParam) GetValue(key string) (string, bool) {
	val, has := this.ctx.RouteParam.GetValue(key)
	if has == false {
		val, has = this.ctx.QueryParam.GetValue(key)
		if has == false {
			val, has = this.ctx.FormParam.GetValue(key)
		}
	}
	return val, has
}

func (this *WebParam) GetString(key string, defaultValue ...string) ValidString {
	val, has := this.GetValue(key)
	if has == false && len(defaultValue) > 0 {
		return ValidString{defaultValue[0], errNotFind}
	} else if has == false {
		return ValidString{"", errNotFind}
	}
	return ValidString{val, nil}
}

func (this *WebParam) MustString(key string, defaultValue ...string) string {
	return this.GetString(key, defaultValue...).Value
}

func (this *WebParam) GetBool(key string, defaultValue ...bool) ValidBool {
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

func (this *WebParam) MustBool(key string, defaultValue ...bool) bool {
	return this.GetBool(key, defaultValue...).Value
}

func (this *WebParam) GetInt(key string, defaultValue ...int) ValidInt {
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

func (this *WebParam) MustInt(key string, defaultValue ...int) int {
	return this.GetInt(key, defaultValue...).Value
}

func (this *WebParam) GetInt32(key string, defaultValue ...int32) ValidInt32 {
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

func (this *WebParam) MustInt32(key string, defaultValue ...int32) int32 {
	return this.GetInt32(key, defaultValue...).Value
}

func (this *WebParam) GetInt64(key string, defaultValue ...int64) ValidInt64 {
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

func (this *WebParam) MustInt64(key string, defaultValue ...int64) int64 {
	return this.GetInt64(key, defaultValue...).Value
}

func (this *WebParam) GetFloat(key string, defaultValue ...float64) ValidFloat {
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

func (this *WebParam) MustFloat(key string, defaultValue ...float64) float64 {
	return this.GetFloat(key, defaultValue...).Value
}

func (this *WebParam) GetTime(key string, format string, defaultValue ...time.Time) ValidTime {
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

func (this *WebParam) MustTime(key string, format string, defaultValue ...time.Time) time.Time {
	return this.GetTime(key, format, defaultValue...).Value
}

func (this *WebParam) GetStruct(ptr interface{}) error {
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