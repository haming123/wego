package wini

import (
	"bytes"
	"errors"
	"os"
	"reflect"
	"time"
)

type ConfigData struct {
	sections	map[string]ConfigSection
}

func (this *ConfigData)ParseFile(file_path string) error {
	err := parseFile(file_path, this)
	if err != nil {
		return err
	}
	return nil
}

func (this *ConfigData) Section(section ...string) ConfigSection {
	sec_name := "root"
	if len(section) > 0 {
		sec_name = section[0]
	}
	sec_map, ok := this.sections[sec_name]
	if !ok {
		return make(ConfigSection)
	}
	return sec_map
}

func (this *ConfigData) SetData(section string, name string, value string) {
	if section == "" {
		section = "root"
	}

	if this.sections == nil {
		this.sections = make(map[string]ConfigSection)
	}

	sec_map, ok := this.sections[section]
	if !ok {
		sec_map = make(ConfigSection)
	}
	sec_map[name] = value
	this.sections[section] = sec_map
}

func (this *ConfigData) GetString (key string, defaultValue ...string) ValidString {
	return this.Section().GetString(key, defaultValue...)
}

func (this *ConfigData) MustString (key string, default_value ...string) string {
	return this.Section().MustString(key, default_value...)
}

func (this *ConfigData) GetBool (key string, defaultValue ...bool) ValidBool {
	return this.Section().GetBool(key, defaultValue...)
}

func (this *ConfigData) MustBool (key string, defaultValue ...bool) bool {
	return this.Section().MustBool(key, defaultValue...)
}

func (this *ConfigData) GetInt (key string, defaultValue ...int) ValidInt {
	return this.Section().GetInt(key, defaultValue...)
}

func (this *ConfigData) MustInt (key string, defaultValue ...int) int {
	return this.Section().MustInt(key, defaultValue...)
}

func (this *ConfigData) GetInt32 (key string, defaultValue ...int32) ValidInt32 {
	return this.Section().GetInt32(key, defaultValue...)
}

func (this *ConfigData) MustInt32 (key string, defaultValue ...int32) int32 {
	return this.Section().MustInt32(key, defaultValue...)
}

func (this *ConfigData) GetInt64 (key string, defaultValue ...int64) ValidInt64 {
	return this.Section().GetInt64(key, defaultValue...)
}

func (this *ConfigData) MustInt64 (key string, defaultValue ...int64) int64 {
	return this.Section().MustInt64(key, defaultValue...)
}

func (this *ConfigData) GetFloat (key string, defaultValue ...float64) ValidFloat {
	return this.Section().GetFloat(key, defaultValue...)
}

func (this *ConfigData) MustFloat (key string, defaultValue ...float64) float64 {
	return this.Section().MustFloat(key, defaultValue...)
}

func (this ConfigData) GetTime (key string, format string, defaultValue ...time.Time) ValidTime {
	return this.Section().GetTime(key, format, defaultValue...)
}

func (this *ConfigData) MustTime (key string, format string, defaultValue ...time.Time) time.Time {
	return this.Section().MustTime(key, format, defaultValue...)
}

func (this ConfigData) GetStrings (key string, delim string) ([]string, error) {
	return this.Section().GetStrings(key, delim)
}

func (this ConfigData) GetBools (key string, delim string) ([]bool, error) {
	return this.Section().GetBools(key, delim)
}

func (this ConfigData) GetInts (key string, delim string) ([]int, error) {
	return this.Section().GetInts(key, delim)
}

func (this ConfigData) GetInt64s (key string, delim string) ([]int64, error) {
	return this.Section().GetInt64s(key, delim)
}

func (this ConfigData) GetFloats (key string, delim string) ([]float64, error) {
	return this.Section().GetFloats(key, delim)
}

func (this *ConfigData) GetStruct(ptr interface{}) error {
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

	root := this.Section()
	t_ent := v_ent.Type()
	f_num := t_ent.NumField()
	for i:=0; i < f_num; i++{
		ff := t_ent.Field(i)
		tag_info := GetTagInfo(ff, "ini")
		if tag_info.FieldName == "" {
			continue
		}
		if ff.Type.Name() == "ConfigData" {
			continue
		}
		fv := v_ent.Field(i)
		if !fv.CanSet() {
			continue
		}

		switch fv.Kind() {
		case reflect.String:
			val, _ := root.getStringByTag(tag_info)
			fv.SetString(val)
		case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
			val, _ := root.getInt64ByTag(tag_info)
			fv.SetInt(val)
		case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
			val, _ := root.getInt64ByTag(tag_info)
			fv.SetUint(uint64(val))
		case reflect.Bool:
			val, _ := root.getBoolByTag(tag_info)
			fv.SetBool(val)
		case reflect.Float32, reflect.Float64:
			val, _ := root.getFloatByTag(tag_info)
			fv.SetFloat(val)
		case reflect.Struct:
			v_ptr := fv.Addr().Interface()
			this.Section(tag_info.FieldName).GetStruct(v_ptr)
		}
	}
	return nil
}

func (this *ConfigData) SaveConfig(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	buff := bytes.NewBuffer(nil)
	root := this.Section("root")
	for key, val := range root {
		buff.WriteString(key)
		buff.WriteString(" = ")
		buff.WriteString(val)
		buff.WriteString("\n")
	}
	buff.WriteString("\n")

	for s_name, section := range this.sections {
		if s_name == "root" {
			continue
		}
		buff.WriteString("[")
		buff.WriteString(s_name)
		buff.WriteString("]\n")
		for key, val := range section {
			buff.WriteString(key)
			buff.WriteString(" = ")
			buff.WriteString(val)
			buff.WriteString("\n")
		}
		buff.WriteString("\n")
	}
	_, err = buff.WriteTo(f)
	return err
}
