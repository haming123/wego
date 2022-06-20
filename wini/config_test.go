package wini

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGet(t *testing.T) {
	var cfg ConfigData
	err := ParseFile("./app.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	//通过GetXXX获取配置项的值
	val := cfg.GetString("app_name")
	if val.Error != nil {
		t.Error(err)
		return
	}
	t.Log(val.Value)
}

func TestSectionGet(t *testing.T) {
	var cfg ConfigData
	err := ParseFile("./app.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	//通过GetXXX获取配置项的值
	val := cfg.Section("mysql").GetString("db_name")
	if val.Error != nil {
		t.Error(err)
		return
	}
	t.Log(val.Value)
}

func TestIniGetXXX(t *testing.T) {
	var cfg ConfigData
	err := ParseFile("./app2.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	val := cfg.GetString("str_value")
	if val.Error != nil {
		t.Error(val.Error)
		return
	}
	t.Log(val.Value)

	val_bool := cfg.GetBool("bool_value")
	if val.Error != nil {
		t.Error(val_bool.Error)
		return
	}
	t.Log(val_bool.Value)

	val_int := cfg.GetInt("int_value")
	if val.Error != nil {
		t.Error(val_int.Error)
		return
	}
	t.Log(val_int.Value)

	val_float := cfg.GetFloat("float_value")
	if val.Error != nil {
		t.Error(val_float.Error)
		return
	}
	t.Log(val_float.Value)
}

func TestIniMustXXX(t *testing.T) {
	var cfg ConfigData
	err := ParseFile("./app2.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(cfg.MustString("str_value"))
	t.Log(cfg.MustBool("bool_value"))
	t.Log(cfg.MustInt("int_value", 222))
	t.Log(cfg.MustFloat("float_value"))
}

func TestGetArray(t *testing.T) {
	var cfg ConfigData
	err := ParseFile("./app2.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	arr, err := cfg.GetInts("ints_value", ",")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(arr)
}

func GetIndentJson(ent interface{}) string {
	data, _ := json.MarshalIndent(ent, "", "    ")
	return string(data)
}

type DbConfig struct {
	MysqlHost string `ini:"db_host"`
	MysqlUser string `ini:"db_user"`
	MysqlPwd  string `ini:"db_pwd"`
	MysqlDb   string `ini:"db_name"`
}

type AppConfig struct {
	AppName  string   `ini:"app_name"`
	HttpPort uint     `ini:"http_port;default=8080"`
	GoPath   string   `ini:"go_path"`
	DbParam  DbConfig `ini:"mysql"`
}

func TestIniGetStruct(t *testing.T) {
	var cfg ConfigData
	err := ParseFile("./app.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	var data AppConfig
	err = cfg.GetStruct(&data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestIniGetSectionStruct(t *testing.T) {
	var cfg ConfigData
	err := ParseFile("./app.conf", &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	var data DbConfig
	err = cfg.Section("mysql").GetStruct(&data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestSaveConfig(t *testing.T) {
	cur_path, err := os.Getwd()
	if err != nil {
		t.Error(err)
		return
	}

	var cfg ConfigData
	file_path := filepath.Join(cur_path, "./app.conf")
	err = ParseFile(file_path, &cfg)
	if err != nil {
		t.Error(err)
		return
	}

	cfg.SaveConfig("./app2.conf")
}
