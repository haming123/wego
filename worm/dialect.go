package worm

import (
	"database/sql"
	"errors"
)

/*
不同的数据库中，SQL语句使用的占位符语法不尽相同。
MySQL	?
PostgreSQL	$1, $2等
SQLite	? 和$1
Oracle	:name
*/

type ColumnInfo struct {
	Name            string
	SQLType         string
	DbType        	string
	Comment         string
	Length 			int
	Length2 		int
	Nullable        bool
	IsPrimaryKey    bool
	IsAutoIncrement bool
}

type Dialect interface {
	GetName() string
	Quote(key string) string
	LimitSql(offset int64, limit int64) string
	ParsePlaceholder(sql_tpl string) string
	DbType2GoType(colType string) string
	GetColumns(db_raw *sql.DB, tableName string) ([]ColumnInfo, error)
}

var g_dialect_map = map[string]Dialect{}

func RegDialect(name string, dialect Dialect) {
	g_dialect_map[name] = dialect
}

func GetDialect(name string) (Dialect, error) {
	dialect, ok := g_dialect_map[name]
	if ok == false {
		return nil, errors.New("incorrect driver name")
	}
	return dialect, nil
}
