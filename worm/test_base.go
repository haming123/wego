package worm

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var DbConn *sql.DB = nil
var SlaveDb *sql.DB = nil
var dsn_mysql string = ""
var dsn_slave string = ""

func readTestConfig(fileName string) (string, error) {
	cur_path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	file_path := filepath.Join(cur_path, fileName)
	data, err := os.ReadFile(file_path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func init() {
	dsn_mysql, _ = readTestConfig("test_mysql.conf")
	if dsn_mysql == "" {
		dsn_mysql, _ = readTestConfig("test_mysql.conf_bak")
	}
	dsn_mysql = strings.Trim(dsn_mysql, "\n")
	dsn_mysql = strings.Trim(dsn_mysql, "\r")

	dsn_slave, _ = readTestConfig("test_slave.conf")
	if dsn_slave == "" {
		dsn_slave, _ = readTestConfig("test_slave.conf_bak")
	}
	dsn_slave = strings.Trim(dsn_slave, "\n")
	dsn_slave = strings.Trim(dsn_slave, "\r")
}

func open_db(cnnstr string) (*sql.DB, error) {
	var err error
	dbcnn, err := sql.Open("mysql", cnnstr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = dbcnn.Ping()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//fmt.Println("******************open_db**********")
	return dbcnn, nil
}

func OpenDb() (*sql.DB, error) {
	if DbConn != nil {
		return DbConn, nil
	}
	dbcnn, err := open_db(dsn_mysql)
	DbConn = dbcnn
	return DbConn, err
}

func OpenSalveDb() (*sql.DB, error) {
	if SlaveDb != nil {
		return SlaveDb, nil
	}
	dbcnn, err := open_db(dsn_mysql)
	SlaveDb = dbcnn
	return SlaveDb, err
}
