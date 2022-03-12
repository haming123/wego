package worm

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

type postgresDialect struct {
	dialectBase
}

func init() {
	RegDialect("postgres", &postgresDialect{})
}

func (db *postgresDialect) GetName() string {
	return "postgres"
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

func (db *postgresDialect)ModelInsertHasOutput(md *DbModel) bool {
	return true
}

func (db *postgresDialect)GenModelInsert(md *DbModel) string {
	var buffer bytes.Buffer
	index := 0;
	buffer.WriteString(fmt.Sprintf("insert into %s (", md.table_name))
	for i, item := range md.flds_addr {
		if md.GetFieldFlag4Insert(i) == false {
			continue
		}
		if index > 0{
			buffer.WriteString(",")
		}
		buffer.WriteString(item.FName)
		index += 1
	}
	buffer.WriteString(")")

	index = 0;
	buffer.WriteString(" values (")
	for i, _ := range md.flds_addr {
		if md.GetFieldFlag4Insert(i) == false {
			continue
		}
		if index > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString("?")
		index += 1
	}
	buffer.WriteString(")")

	//mssql not support LastInsertId
	//use RETURNING and queryrrow to get LastInsertId
	if len(md.field_id) > 0 {
		buffer.WriteString(" RETURNING ")
		buffer.WriteString(md.field_id)
		buffer.WriteString(" ")
	} else {
		buffer.WriteString(" RETURNING 0 ")
	}

	return buffer.String()
}

func (db *postgresDialect)TableInsertHasOutput(tb *DbTable) bool {
	return len(tb.return_str) > 0
}

func (db *postgresDialect)GenTableInsertSql(tb *DbTable) (string, []interface{}) {
	index := 0;
	vals:= []interface{}{}

	var buffer1 bytes.Buffer
	buffer1.WriteString(fmt.Sprintf("insert into %s (", tb.table_name))

	var buffer2 bytes.Buffer
	buffer2.WriteString(" values (")
	for name, val := range tb.fld_values {
		if index > 0 {
			buffer1.WriteString(",")
			buffer2.WriteString(",")
		}

		buffer1.WriteString(name)
		if val == nil {
			buffer2.WriteString("null")
		} else if exp, ok := val.(SqlExp); ok {
			buffer2.WriteString(exp.Tpl_sql)
			if exp.Values != nil {
				vals = append(vals, exp.Values...)
			}
		} else {
			buffer2.WriteString("?")
			vals = append(vals, val)
		}

		index += 1
	}
	buffer1.WriteString(")")
	buffer2.WriteString(")")
	buffer1.Write(buffer2.Bytes())

	//mssql not support LastInsertId
	//use RETURNING and queryrrow to get LastInsertId
	if len(tb.return_str) > 0 {
		buffer1.WriteString(" RETURNING ")
		buffer1.WriteString(tb.return_str)
		buffer1.WriteString(" ")
	} else {
		buffer1.WriteString(" RETURNING 0 ")
	}

	return buffer1.String(), vals
}
