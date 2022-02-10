package worm

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

type dialectMssql struct {
}

func init() {
	RegDialect("mssql", &dialectMssql{})
}

func (db *dialectMssql) GetName() string {
	return "mssql"
}

func (db *dialectMssql) Quote(key string) string {
	return fmt.Sprintf("[%s]", key)
}

func (db *dialectMssql) LimitSql(offset int64, limit int64) string  {
	return fmt.Sprintf(" limit %d, %d ", offset, limit)
}

func (db *dialectMssql) ParsePlaceholder(sql_tpl string) string {
	tpl_str := sql_tpl
	var buffer bytes.Buffer
	for i:=0; i < len(sql_tpl); i++ {
		index := strings.Index(tpl_str, "?")
		if index < 0 {
			break;
		}
		txt_str := tpl_str[0:index]
		tpl_str = tpl_str[index+1:]
		bindvar := fmt.Sprintf("@p%d", i+1)
		buffer.WriteString(txt_str)
		buffer.WriteString(bindvar)
	}
	if len(tpl_str) > 0 {
		buffer.WriteString(tpl_str)
	}
	return buffer.String()
}

func (p *dialectMssql) DbType2GoType(colType string) string {
	switch colType {
	case "VARCHAR", "TEXT", "CHAR", "NVARCHAR", "NCHAR", "NTEXT":
		return "string"
	case "DATE", "DATETIME", "DATETIME2", "TIME":
		return "string"
	case "FLOAT", "REAL":
		return "float64"
	case "BIGINT", "DATETIMEOFFSET":
		return "int64"
	case "TINYINT", "SMALLINT", "INT":
		return "int32"
	default:
		return "[]byte"
	}
}

func (db *dialectMssql) GetColumns(db_raw *sql.DB, table_name string) ([]ColumnInfo, error) {
	sql_str := "select a.name as name, b.name as ctype, a.is_nullable as nullable, ISNULL(p.is_primary_key, 0), a.is_identity as is_identity "
	sql_str += "from sys.columns a left join sys.types b on a.user_type_id=b.user_type_id "
	sql_str += "where a.object_id=object_id('" + table_name + "')"
	rows, err := db_raw.Query(sql_str)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]ColumnInfo, 0)
	for rows.Next() {
		var col_name, ctype string
		var nullable, isPK, isIncrement bool
		err = rows.Scan(&col_name, &ctype, &nullable, &isPK, &isIncrement)
		if err != nil {
			return nil, err
		}

		var col ColumnInfo
		col.Name = strings.Trim(col_name, "` ")
		col.Nullable = nullable
		col.IsPrimaryKey = isPK
		col.IsAutoIncrement = isIncrement
		col.SQLType = strings.ToUpper(ctype)
		col.DbType = col.SQLType

		cols = append(cols, col)
	}
	return  cols, nil
}