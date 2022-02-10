package wego

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ParamValue interface {
	GetValue(key string) (string, bool)
	GetValues(key string) ([]string, bool)
}

//将urql参数解析到struct
//只进行struct的的一层解析
type FuncGetValues func(string, bool) ([]string, bool)
func bindStruct(v_ent reflect.Value, pv ParamValue) error {
	t_ent := v_ent.Type()
	f_num := t_ent.NumField()
	for i:=0; i < f_num; i++{
		ff := t_ent.Field(i)
		tag_info := GetTagInfo(ff, "form")
		if tag_info.FieldName == "" {
			continue
		}

		fv := v_ent.Field(i)
		if !fv.CanSet() {
			continue
		}

		if fv.Kind() == reflect.Ptr {
			v_ptr_new := reflect.New(fv.Type().Elem())
			fv.Set(v_ptr_new)
			fv = v_ptr_new.Elem()
		}

		switch fv.Kind() {
		case reflect.Slice:
			datas, ok := pv.GetValues(tag_info.FieldName)
			if ok == false && tag_info.HasValue {
				datas = []string{ tag_info.DefValue }
			}
			if len(datas) < 1 {
				continue
			}
			err := setSlice(fv, datas)
			if err != nil {
				return err
			}
		case reflect.Array:
			datas, ok := pv.GetValues(tag_info.FieldName)
			if ok == false && tag_info.HasValue {
				datas = []string{ tag_info.DefValue }
			}
			if len(datas) < 1 {
				continue
			}
			if len(datas) == fv.Len() {
				err := setArray(fv, datas)
				if err != nil {
					return err
				}
			}
		default:
			val_str, ok := pv.GetValue(tag_info.FieldName)
			if ok == false && tag_info.HasValue {
				val_str = tag_info.DefValue
				ok = true
			}
			if ok == false {
				continue
			}
			err := setFieldValue(fv, val_str)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func setSlice(fv reflect.Value, datas []string) error {
	slice := reflect.MakeSlice(fv.Type(), len(datas), len(datas))
	err := setArray(slice, datas)
	if err == nil {
		fv.Set(slice)
	}
	return err
}

func setArray(fv reflect.Value, datas []string) error {
	for i, data := range datas {
		err := setFieldValue(fv.Index(i), data)
		if err != nil {
			return err
		}
	}
	return nil
}

var errUnsupportedType = errors.New("unsupported type")
func setFieldValue(fv reflect.Value, data string) error {
	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setIntXField(fv, data)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUintXField(fv, data)
	case reflect.Float32, reflect.Float64:
		return setFloatXField(fv, data)
	case reflect.Bool:
		return setBoolField(fv, data)
	case reflect.String:
		fv.SetString(data)
	case reflect.Struct:
		if _, ok := fv.Interface().(time.Time); ok {
			return setTimeField(fv, data)
		} else {
			return errUnsupportedType
		}
	default:
		return errUnsupportedType
	}
	return nil
}

func setIntXField(fv reflect.Value, data string) error {
	val, err := strconv.ParseInt(data, 10, 64)
	if err == nil {
		fv.SetInt(val)
	}
	return err
}

func setUintXField(fv reflect.Value, data string) error {
	val, err := strconv.ParseUint(data, 10, 64)
	if err == nil {
		fv.SetUint(val)
	}
	return err
}

func setFloatXField(fv reflect.Value, data string) error {
	val, err := strconv.ParseFloat(data, 64)
	if err == nil {
		fv.SetFloat(val)
	}
	return err
}

func setBoolField(fv reflect.Value, data string) error {
	if data == "" { data = "false" }
	val, err := strconv.ParseBool(data)
	if err == nil {
		fv.SetBool(val)
	}
	return err
}

func setTimeField(fv reflect.Value, data string) error {
	val, err := ParseTime(data)
	if err == nil {
		fv.Set(reflect.ValueOf(val))
	}
	return err
}

const (
	formatTime      = "15:04:05"
	formatDate      = "2006-01-02"
	formatDateTime  = "2006-01-02 15:04:05"
	formatDateTimeT = "2006-01-02T15:04:05"
)
func ParseTime(value string) (time.Time, error) {
	var pattern string
	if len(value) >= 25 {
		value = value[:25]
		pattern = time.RFC3339
	} else if strings.HasSuffix(strings.ToUpper(value), "Z") {
		pattern = time.RFC3339
	} else if len(value) >= 19 {
		if strings.Contains(value, "T") {
			pattern = formatDateTimeT
		} else {
			pattern = formatDateTime
		}
		value = value[:19]
	} else if len(value) >= 10 {
		if len(value) > 10 {
			value = value[:10]
		}
		pattern = formatDate
	} else if len(value) >= 8 {
		if len(value) > 8 {
			value = value[:8]
		}
		pattern = formatTime
	}
	return time.ParseInLocation(pattern, value, time.Local)
}