package main

import (
	"github.com/haming123/wego/config"
)

type AppConfig struct {
	DbHost string 		`ini:"db_host"`
	DbPort string 		`ini:"db_port"`
	DbName string 		`ini:"db_db"`
	DbUser string 		`ini:"db_user"`
	DbPwd  string 		`ini:"db_pwd"`
	PkgName string 		`ini:"pkg_name"`
	CreateTime string 	`ini:"create_time"`
	UseFieldTag bool	`ini:"use_field_tag"`
	UseModelPool bool	`ini:"use_model_pol"`
}

func ReadAppConfig(conf_file ...string) (AppConfig, error) {
	var conf AppConfig
	data, err := config.InitConfigData(conf_file...)
	if err != nil {
		return conf, err
	}

	err = data.GetStruct(&conf)
	if err != nil {
		return conf, err
	}

	if conf.PkgName == "" {
		conf.PkgName = "model"
	}

	return conf, nil
}
