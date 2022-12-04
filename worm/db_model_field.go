package worm

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func set_value(fld_ptr interface{}, val interface{}) error {
	if fld_ptr == nil {
		return errors.New("fld_ptr must be Pointer")
	}
	v_fld := reflect.ValueOf(fld_ptr)
	if v_fld.Kind() != reflect.Ptr {
		return errors.New("fld_ptr must be Pointer")
	}

	v_fld = reflect.Indirect(v_fld)
	v_val := reflect.ValueOf(val)
	t_fld := v_fld.Type()
	t_val := v_val.Type()
	k_fld := t_fld.Kind()
	k_val := t_val.Kind()
	if k_fld == k_val {
		v_fld.Set(v_val)
		return nil
	}

	//类型不相同，但可以进行转换，转换后赋值
	if k_val == reflect.String {
		switch k_fld {
		case reflect.Int, reflect.Int32, reflect.Int64:
			tmp, err := strconv.ParseInt(val.(string), 10, 64)
			if err != nil {
				return err
			}
			v_fld.SetInt(tmp)
			return nil
		case reflect.Float32, reflect.Float64:
			tmp, err := strconv.ParseFloat(val.(string), 64)
			if err != nil {
				return err
			}
			v_fld.SetFloat(tmp)
			return nil
		}
	} else if k_fld == reflect.Int64 {
		switch k_val {
		case reflect.Int, reflect.Int32:
			v_fld.SetInt(v_val.Int())
			return nil
		}
	} else if k_fld == reflect.Int32 || k_fld == reflect.Int {
		switch k_val {
		case reflect.Int, reflect.Int32:
			v_fld.SetInt(v_val.Int())
			return nil
		}
	} else if k_fld == reflect.Float64 {
		switch k_val {
		case reflect.Float32:
			v_fld.SetFloat(v_val.Float())
			return nil
		}
	}

	return errors.New(fmt.Sprintf("incorrect data type: %v != %v", t_fld, t_val))
}

func SetValue(md *DbModel, fld_ptr interface{}, val interface{}) error {
	if md != nil {
		return md.SetValue(fld_ptr, val)
	} else {
		return set_value(fld_ptr, val)
	}
}

func GetPointer(md *DbModel, fld_ptr interface{}) interface{} {
	if md == nil {
		return fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return fld_ptr
	}
}

func GetIndirect(md *DbModel, fld_ptr interface{}) interface{} {
	if fld_ptr == nil {
		return nil
	}

	v_fld := reflect.Indirect(reflect.ValueOf(fld_ptr))
	if md == nil {
		return v_fld.Interface()
	} else {
		md.SelectX(fld_ptr)
		return v_fld.Interface()
	}
}

func GetBool(md *DbModel, fld_ptr *bool) bool {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetInt(md *DbModel, fld_ptr *int) int {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetInt32(md *DbModel, fld_ptr *int32) int32 {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetInt64(md *DbModel, fld_ptr *int64) int64 {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetFloat32(md *DbModel, fld_ptr *float32) float32 {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetFloat64(md *DbModel, fld_ptr *float64) float64 {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetString(md *DbModel, fld_ptr *string) string {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}

func GetTime(md *DbModel, fld_ptr *time.Time) time.Time {
	if md == nil {
		return *fld_ptr
	} else {
		md.SelectX(fld_ptr)
		return *fld_ptr
	}
}
