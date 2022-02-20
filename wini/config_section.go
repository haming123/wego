package wini

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var errNotFind = errors.New("not find")
type ConfigSection map[string]string

func (this ConfigSection) GetString (key string, defaultValue ...string) ValidString {
	val, has := this[key]
	if has == false && len(defaultValue) > 0 {
		return ValidString{defaultValue[0], errNotFind}
	} else if has == false {
		return ValidString{"", errNotFind}
	}
	return ValidString{val, nil}
}

func (this ConfigSection) MustString (key string, defaultValue ...string) string {
	return this.GetString(key, defaultValue...).Value
}

func (this ConfigSection) GetBool(key string, defaultValue ...bool) ValidBool {
	val_str, has := this[key]
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

func (this ConfigSection) MustBool(key string, defaultValue ...bool) bool {
	return this.GetBool(key, defaultValue...).Value
}

func (this ConfigSection) GetInt(key string, defaultValue ...int) ValidInt {
	val_str, has := this[key]
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

func (this ConfigSection) MustInt(key string, defaultValue ...int) int {
	return this.GetInt(key, defaultValue...).Value
}

func (this ConfigSection) GetInt32(key string, defaultValue ...int32) ValidInt32 {
	val_str, has := this[key]
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

func (this ConfigSection) MustInt32(key string, defaultValue ...int32) int32 {
	return this.GetInt32(key, defaultValue...).Value
}

func (this ConfigSection) GetInt64(key string, defaultValue ...int64) ValidInt64 {
	val_str, has := this[key]
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

func (this ConfigSection) MustInt64(key string, defaultValue ...int64) int64 {
	return this.GetInt64(key, defaultValue...).Value
}

func (this ConfigSection) GetFloat(key string, defaultValue ...float64) ValidFloat {
	val_str, has := this[key]
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

func (this ConfigSection) MustFloat(key string, defaultValue ...float64) float64 {
	return this.GetFloat(key, defaultValue...).Value
}

func (this ConfigSection) GetTime(key string, format string, defaultValue ...time.Time) ValidTime {
	val_str, has := this[key]
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

func (this ConfigSection) MustTime(key string, format string, defaultValue ...time.Time) time.Time {
	return this.GetTime(key, format, defaultValue...).Value
}

func (this ConfigSection) GetStrings (key string, delim string) ([]string, error) {
	str_data, has := this[key]
	if has == false {
		return []string{}, errNotFind
	}
	if str_data == "" {
		return []string{str_data}, nil
	}
	strs := strings.Split(str_data, delim)
	return strs, nil
}

func (this ConfigSection) GetBools (key string, delim string) ([]bool, error) {
	strs, err := this.GetStrings(key, delim)
	if err != nil {
		return []bool{}, err
	}
	vals := make([]bool, len(strs))
	for i, str := range strs {
		str = strings.Trim(str, " ")
		val, err := strconv.ParseBool(str)
		if err != nil {
			return []bool{}, err
		}
		vals[i] = val
	}
	return vals, nil
}

func (this ConfigSection) GetInts (key string, delim string) ([]int, error) {
	strs, err := this.GetStrings(key, delim)
	if err != nil {
		return []int{}, err
	}
	vals := make([]int, len(strs))
	for i, str := range strs {
		str = strings.Trim(str, " ")
		val, err := strconv.Atoi(str)
		if err != nil {
			return []int{}, err
		}
		vals[i] = val
	}
	return vals, nil
}

func (this ConfigSection) GetInt64s (key string, delim string) ([]int64, error) {
	strs, err := this.GetStrings(key, delim)
	if err != nil {
		return []int64{}, err
	}
	vals := make([]int64, len(strs))
	for i, str := range strs {
		str = strings.Trim(str, " ")
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return []int64{}, err
		}
		vals[i] = val
	}
	return vals, nil
}

func (this ConfigSection) GetFloats (key string, delim string) ([]float64, error) {
	strs, err := this.GetStrings(key, delim)
	if err != nil {
		return []float64{}, err
	}
	vals := make([]float64, len(strs))
	for i, str := range strs {
		str = strings.Trim(str, " ")
		val, err := strconv.ParseFloat(str,  64)
		if err != nil {
			return []float64{}, err
		}
		vals[i] = val
	}
	return vals, nil
}

func (this ConfigSection)getStringByTag (tag TagInfo) (string, error) {
	str_data, has := this[tag.FieldName]
	if has == false && tag.HasValue {
		str_data = tag.DefValue
	} else if has == false {
		return "", errNotFind
	}
	return str_data, nil
}

func (this ConfigSection)getBoolByTag (tag TagInfo) (bool, error) {
	str_data, err := this.getStringByTag(tag)
	if err != nil  {
		return false, err
	}
	return strconv.ParseBool(str_data)
}

func (this ConfigSection)getIntByTag (tag TagInfo) (int, error) {
	str_data, err := this.getStringByTag(tag)
	if err != nil  {
		return 0, err
	}
	return strconv.Atoi(str_data)
}

func (this ConfigSection)getInt64ByTag (tag TagInfo) (int64, error) {
	str_data, err := this.getStringByTag(tag)
	if err != nil  {
		return 0, err
	}
	return strconv.ParseInt(str_data, 10, 64)
}

func (this ConfigSection)getFloatByTag (tag TagInfo) (float64, error) {
	str_data, err := this.getStringByTag(tag)
	if err != nil  {
		return 0, err
	}
	return strconv.ParseFloat(str_data,  64)
}

func (this ConfigSection) GetStruct(ptr interface{}) error {
	//ptr必须是一个非空指针
	if ptr == nil {
		return errors.New("ptr must be *Struct")
	}
	//ptr必须是一个结构体指针
	v_ent := reflect.ValueOf(ptr)
	if v_ent.Kind() != reflect.Ptr {
		return errors.New("ptr must be *Struct")
	}
	//ptr指向的对象必须是一个结构体
	v_ent = reflect.Indirect(v_ent)
	if v_ent.Kind() != reflect.Struct {
		return errors.New("ptr must be *Struct")
	}

	t_ent := v_ent.Type()
	f_num := t_ent.NumField()
	for i:=0; i < f_num; i++{
		ff := t_ent.Field(i)
		tag_info := GetTagInfo(ff, "ini")
		if tag_info.FieldName == "" {
			continue
		}

		fv := v_ent.Field(i)
		if !fv.CanSet() {
			continue
		}

		switch fv.Kind() {
		case reflect.String:
			val, _ := this.getStringByTag(tag_info)
			fv.SetString(val)
		case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
			val, _ := this.getInt64ByTag(tag_info)
			fv.SetInt(val)
		case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
			val, _ := this.getInt64ByTag(tag_info)
			fv.SetUint(uint64(val))
		case reflect.Bool:
			val, _ := this.getBoolByTag(tag_info)
			fv.SetBool(val)
		case reflect.Float32, reflect.Float64:
			val, _ := this.getFloatByTag(tag_info)
			fv.SetFloat(val)
		}
	}
	return nil
}