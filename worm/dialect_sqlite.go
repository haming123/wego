package worm

import (
	"database/sql"
	"fmt"
)

type dialectSqlite struct {
	dialectBase
}

func init() {
	RegDialect("sqlite", &dialectSqlite{})
}

func (db *dialectSqlite) GetName() string {
	return "sqlite"
}

func (db *dialectSqlite) LimitSql(offset int64, limit int64) string  {
	return fmt.Sprintf(" limit %d, %d ", offset, limit)
}

func (p *dialectSqlite) DbType2GoType(colType string) string {
	switch colType {
	case "TEXT":
		return "string"
	case "INTEGER":
		return "int64"
	case "DATETIME":
		return "time.Time"
	case "REAL":
		return "float64"
	case "NUMERIC", "DECIMAL":
		return "string"
	case "BLOB":
		return "[]byte"
	default:
		return "string"
	}
}

func (db *dialectSqlite) GetColumns(db_raw *sql.DB, table_name string) ([]ColumnInfo, error) {
	return  nil, nil
}


