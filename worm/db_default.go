package worm

import (
	"database/sql"
	"errors"
)

var db_default *DbEngine

func InitEngine(driverName string, db *sql.DB) error {
	var err error
	db_default, err = NewEngine(driverName, db)
	return err
}

func InitMysql(db *sql.DB) error {
	var err error
	db_default, err = NewMysql(db)
	return err
}

func InitPostgres(db *sql.DB) error {
	var err error
	db_default, err = NewPostgres(db)
	return err
}

func InitSqlServer(db *sql.DB) error {
	var err error
	db_default, err = NewSqlServer(db)
	return err
}

func InitSqlite3(db *sql.DB) error {
	var err error
	db_default, err = NewSqlite3(db)
	return err
}

func AddSlave(db *sql.DB, db_name string, weight int) error {
	engine := db_default
	if engine == nil {
		return errors.New("engine is nil")
	}
	return engine.AddSlave(db, db_name, weight)
}

func UsePrepare(flag bool)  {
	engine := db_default
	if engine == nil {
		return
	}
	engine.UsePrepare(flag)
	if engine.def_session != nil {
		engine.def_session.UsePrepare(flag)
	}
}

func ShowSqlLog(flag bool)  {
	engine := db_default
	if engine == nil {
		return
	}
	engine.ShowSqlLog(flag)
	if engine.def_session != nil {
		engine.def_session.ShowLog(flag)
	}
}

func SetSqlLogCB(cb SqlPrintCB)  {
	engine := db_default
	if engine == nil {
		return
	}
	engine.SetSqlLogCB(cb)
}

func SetMaxStmtCacheNum(max_len int)  {
	engine := db_default
	if engine == nil {
		return
	}
	engine.SetMaxStmtCacheNum(max_len)
}

func Expr(expr string, args ...interface{}) SqlExp {
	return SqlExp{Tpl_sql: expr, Values: args}
}

func NewSession() *DbSession {
	engine := db_default
	if engine == nil {
		return nil
	}
	return engine.NewSession()
}

func Model(ent_ptr interface{}) *DbModel {
	engine := db_default
	if engine == nil {
		return nil
	}
	dbs := engine.def_session
	return dbs.Model(ent_ptr)
}

func Joint(ent_ptr interface{}, alias string, fields ...string) *DbJoint {
	engine := db_default
	if engine == nil {
		return nil
	}
	dbs := engine.def_session
	return dbs.Joint(ent_ptr, alias, fields...)
}

func SQL(sql_str string, args ...interface{}) *DbSQL {
	engine := db_default
	if engine == nil {
		return nil
	}
	dbs := engine.def_session
	return dbs.SQL(sql_str, args...)
}

func Table(table_name string) *DbTable {
	engine := db_default
	if engine == nil {
		return nil
	}
	dbs := engine.def_session
	return dbs.Table(table_name)
}
