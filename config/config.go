package config

import (
	"os"
	"path/filepath"
	"time"
)

var config_data ConfigData
func InitConfigData(file_name ...string) (*ConfigData, error) {
	fileName := "app.conf"
	if len(file_name) > 0 && file_name[0] != "" {
		fileName = file_name[0]
	}

	cur_path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	file_path := filepath.Join(cur_path, fileName)
	err = ParseFile(file_path, &config_data)
	if err != nil {
		return nil, err
	}

	return &config_data, nil
}

func Section(section ...string) ConfigSection {
	this := &config_data
	return this.Section(section...)
}

func GetString (key string, defaultValue ...string) ValidString {
	this := &config_data
	return this.GetString(key, defaultValue...)
}

func MustString (key string, defaultValue ...string) string {
	this := &config_data
	return this.MustString(key, defaultValue...)
}

func GetBool (key string, defaultValue ...bool) ValidBool {
	this := &config_data
	return this.GetBool(key, defaultValue...)
}

func MustBool (key string, defaultValue ...bool) bool {
	this := &config_data
	return this.MustBool(key, defaultValue...)
}

func GetInt (key string, defaultValue ...int) ValidInt {
	this := &config_data
	return this.GetInt(key, defaultValue...)
}

func MustInt (key string, defaultValue ...int) int {
	this := &config_data
	return this.MustInt(key, defaultValue...)
}

func GetInt32 (key string, defaultValue ...int32) ValidInt32 {
	this := &config_data
	return this.GetInt32(key, defaultValue...)
}

func MustInt32 (key string, defaultValue ...int32) int32 {
	this := &config_data
	return this.MustInt32(key, defaultValue...)
}

func GetInt64 (key string, defaultValue ...int64) ValidInt64 {
	this := &config_data
	return this.GetInt64(key, defaultValue...)
}

func MustInt64 (key string, defaultValue ...int64) int64 {
	this := &config_data
	return this.MustInt64(key, defaultValue...)
}

func GetFloat (key string, defaultValue ...float64) ValidFloat {
	this := &config_data
	return this.GetFloat(key, defaultValue...)
}

func MustFloat (key string, defaultValue ...float64) float64 {
	this := &config_data
	return this.MustFloat(key, defaultValue...)
}

func GetTime (key string, format string, defaultValue ...time.Time) ValidTime {
	this := &config_data
	return this.GetTime(key, format, defaultValue...)
}

func MustTime (key string, format string, defaultValue ...time.Time) time.Time {
	this := &config_data
	return this.MustTime(key, format, defaultValue...)
}

func GetStrings (key string, delim string) ([]string, error) {
	this := &config_data
	return this.GetStrings(key, delim)
}

func GetBools (key string, delim string) ([]bool, error) {
	this := &config_data
	return this.GetBools(key, delim)
}

func GetInts (key string, delim string) ([]int, error) {
	this := &config_data
	return this.GetInts(key, delim)
}

func GetInt64s (key string, delim string) ([]int64, error) {
	this := &config_data
	return this.GetInt64s(key, delim)
}

func GetFloats (key string, delim string) ([]float64, error) {
	this := &config_data
	return this.GetFloats(key, delim)
}

func GetStruct (ptr interface{}) error {
	this := &config_data
	return this.GetStruct(ptr)
}
