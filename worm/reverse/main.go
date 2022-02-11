package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/haming123/wego/worm"
)

//./reverse.exe -c ./appconfig.yaml -t user
var g_cfg AppConfig
func main() {
	conf_file := ""
	table_name := ""
	model_file := ""
	flag.StringVar(&conf_file, "c", "", "config file")
	flag.StringVar(&table_name, "t", "", "table name for generate model code")
	flag.StringVar(&model_file, "s", "", "model file")
	flag.Parse()

	if len(table_name) < 1 {
		fmt.Println("please input table name (usage: -c config_file -t table name -s model_file)")
		return
	}
	fmt.Printf("table_name:%s\n", table_name)

	var err error
	g_cfg, err = ReadAppConfig(conf_file)
	if err != nil {
		fmt.Println(err)
		return
	}

	cnnstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		g_cfg.DbUser, g_cfg.DbPwd, g_cfg.DbHost, g_cfg.DbPort, g_cfg.DbName)
	fmt.Println(cnnstr)
	db_drv, err := sql.Open("mysql", cnnstr)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db_drv.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	db, err := worm.NewEngine("mysql", db_drv)
	if err != nil {
		fmt.Println( err)
		return
	}

	table_name2 := fmt.Sprintf("%s.%s", g_cfg.DbName, table_name)
	CodeGen4Table(db, table_name2, model_file)
}

