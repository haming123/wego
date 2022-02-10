package worm

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

type postgresDialect struct {
}

func init() {
	RegDialect("postgres", &postgresDialect{})
}

func (db *postgresDialect) GetName() string {
	return "postgres"
}

func (db *postgresDialect) Quote(key string) string {
	return fmt.Sprintf("\"%s\"", key)
}

func (db *postgresDialect) LimitSql(offset int64, limit int64) string  {
	return fmt.Sprintf(" limit %d offset %d ", limit, offset)
}

func (db *postgresDialect) ParsePlaceholder(sql_tpl string) string {
	tpl_str := sql_tpl
	var buffer bytes.Buffer
	for i:=0; i < len(sql_tpl); i++ {
		index := strings.Index(tpl_str, "?")
		if index < 0 {
			break;
		}
		txt_str := tpl_str[0:index]
		tpl_str = tpl_str[index+1:]
		bindvar := fmt.Sprintf("$%d", i+1)
		buffer.WriteString(txt_str)
		buffer.WriteString(bindvar)
	}
	if len(tpl_str) > 0 {
		buffer.WriteString(tpl_str)
	}
	return buffer.String()
}

func (p *postgresDialect) DbType2GoType(colType string) string {
	switch colType {
	case "VARCHAR", "TEXT":
		return "string"
	case "BIGINT", "BIGSERIAL":
		return "int64"
	case "SMALLINT", "INT", "INT8", "INT4", "INTEGER", "SERIAL":
		return "int32"
	case "FLOAT", "FLOAT4", "REAL", "DOUBLE PRECISION":
		return "float64"
	case "DATETIME", "TIMESTAMP":
		return "time.Time"
	case "BOOL":
		return "bool"
	default:
		return "[]byte"
	}
}

func (db *postgresDialect) GetColumns(db_raw *sql.DB,tableName string) ([]ColumnInfo, error) {
	return  nil, nil
}


