package worm

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type dialectMysql struct {
	DialectBase
}

func (db *dialectMysql) GetName() string {
	return "mysql"
}

func (db *dialectMysql) LimitSql(offset int64, limit int64) string  {
	return fmt.Sprintf(" limit %d, %d ", offset, limit)
}

func (db *dialectMysql) DbType2GoType(colType string) string {
	switch colType {
	case "CHAR", "VARCHAR", "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT", "ENUM", "SET":
		return "string"
	case "BIGINT":
		return "int64"
	case "TINYINT", "SMALLINT", "MEDIUMINT", "INT":
		return "int32"
	case "FLOAT", "REAL", "DOUBLE PRECISION", "DOUBLE":
		return "float64"
	case "DECIMAL", "NUMERIC":
		return "float64"
	case "DATETIME", "TIMESTAMP":
		return "time.Time"
	case "BIT":
		return "[]byte"
	case "BINARY", "VARBINARY", "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB":
		return "[]byte"
	default:
		return ""
	}
}

func (db *dialectMysql) GetColumns(db_raw *sql.DB, tableName string) ([]ColumnInfo, error) {
	strs :=strings.Split(tableName, ".")
	if len(strs) != 2 {
		return nil, errors.New("table name must be dbname.tablename")
	}
	db_name := strs[0]
	tb_name := strs[1]

	sql_str := fmt.Sprintf("SELECT COLUMN_NAME, COLUMN_TYPE, COLUMN_KEY, EXTRA, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA='%s' AND TABLE_NAME='%s'", db_name, tb_name)
	rows, err := db_raw.Query(sql_str)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]ColumnInfo, 0)
	for rows.Next() {
		var col_name, col_type, col_key, extra, comment sql.NullString
		err = rows.Scan(&col_name, &col_type, &col_key, &extra, &comment)
		if err != nil {
			return nil, err
		}

		var col ColumnInfo
		col.Name = col_name.String
		col.Comment = comment.String

		fields := strings.Fields(col_type.String)
		db_type := fields[0]
		cts := strings.Split(db_type, "(")
		db_type = cts[0]
		db_type = strings.ToUpper(db_type)
		col.SQLType = col_type.String
		col.DbType = db_type

		var len1, len2 int
		if len(cts) == 2 && db_type != "ENUM" && db_type != "SET"  {
			idx := strings.Index(cts[1], ")")
			lens := strings.Split(cts[1][0:idx], ",")
			len1, err = strconv.Atoi(strings.TrimSpace(lens[0]))
			if err != nil {
				return nil, err
			}
			if len(lens) == 2 {
				len2, err = strconv.Atoi(lens[1])
				if err != nil {
					return nil, err
				}
			}
		}
		col.Length = len1
		col.Length2 = len2

		if col_key.String == "PRI" {
			col.IsPrimaryKey = true
		}
		if extra.String == "auto_increment" {
			col.IsAutoIncrement = true
		}

		cols = append(cols, col)
	}

	return  cols, nil
}

