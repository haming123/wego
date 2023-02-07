package worm

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

type dialectMssql struct {
	DialectBase
}

func (db *dialectMssql) GetName() string {
	return "mssql"
}

func (db *dialectMssql) LimitSql(offset int64, limit int64) string {
	return ""
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
	sql_str := "select a.name as name, b.name as ctype, a.is_nullable as nullable, a.is_identity as is_identity "
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
		var nullable, isIncrement bool
		err = rows.Scan(&col_name, &ctype, &nullable, &isIncrement)
		if err != nil {
			return nil, err
		}

		var col ColumnInfo
		col.Name = strings.Trim(col_name, "` ")
		col.Nullable = nullable
		col.IsAutoIncrement = isIncrement
		col.SQLType = strings.ToUpper(ctype)
		col.DbType = col.SQLType

		cols = append(cols, col)
	}
	return cols, nil
}

func (db *dialectMssql) ModelInsertHasOutput(md *DbModel) bool {
	return true
}

func (db *dialectMssql) GenModelInsertSql(md *DbModel) string {
	var buffer bytes.Buffer
	index := 0
	buffer.WriteString(fmt.Sprintf("insert into %s (", md.table_name))
	for i, item := range md.flds_addr {
		if md.GetFieldFlag4Insert(i) == false {
			continue
		}
		if index > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(item.FName)
		index += 1
	}
	buffer.WriteString(")")

	//mssql not support LastInsertId
	//use output and queryrrow to get LastInsertId
	if len(md.field_id) > 0 {
		buffer.WriteString(" OUTPUT Inserted.")
		buffer.WriteString(md.field_id)
		buffer.WriteString(" ")
	} else {
		buffer.WriteString(" OUTPUT 0 ")
	}

	index = 0
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

	return buffer.String()
}

func (db *dialectMssql) GenModelGetSql(md *DbModel) string {
	var buffer bytes.Buffer

	buffer.WriteString("select top 1 ")
	buffer.WriteString(md.gen_select_fields())
	buffer.WriteString(" from ")
	buffer.WriteString(md.table_name)
	if len(md.table_alias) > 0 {
		buffer.WriteString(" ")
		buffer.WriteString(md.table_alias)
	}

	if len(md.db_where.Tpl_sql) > 0 {
		buffer.WriteString(" where ")
		buffer.WriteString(md.db_where.Tpl_sql)
	}

	if len(md.group_by) > 0 {
		buffer.WriteString(" group by ")
		buffer.WriteString(md.group_by)
	}

	if len(md.order_by) > 0 {
		buffer.WriteString(" order by ")
		buffer.WriteString(md.order_by)
	}

	return buffer.String()
}

// select top @pageSize id from tablename
// where id not in (
// select top @offset id from tablename
// )
func (db *dialectMssql) GenModelFindSql(md *DbModel) string {
	var buffer bytes.Buffer

	buffer.WriteString("select ")
	if md.db_limit > 0 {
		buffer.WriteString(fmt.Sprintf("top %d ", md.db_limit))
	}
	buffer.WriteString(md.gen_select_fields())
	buffer.WriteString(" from ")
	buffer.WriteString(md.table_name)
	if len(md.table_alias) > 0 {
		buffer.WriteString(" ")
		buffer.WriteString(md.table_alias)
	}

	if len(md.db_where.Tpl_sql) > 0 {
		buffer.WriteString(" where ")
		buffer.WriteString(md.db_where.Tpl_sql)
	}

	if len(md.group_by) > 0 {
		buffer.WriteString(" group by ")
		buffer.WriteString(md.group_by)
	}

	if len(md.order_by) > 0 {
		buffer.WriteString(" order by ")
		buffer.WriteString(md.order_by)
	}

	return buffer.String()
}

func (db *dialectMssql) TableInsertHasOutput(tb *DbTable) bool {
	return len(tb.output_str) > 0
}

func (db *dialectMssql) GenTableInsertSql(tb *DbTable) (string, []interface{}) {
	index := 0
	vals := []interface{}{}

	var buffer1 bytes.Buffer
	buffer1.WriteString(fmt.Sprintf("insert into %s (", tb.table_name))

	var buffer2 bytes.Buffer
	buffer2.WriteString(" values (")
	for _, fld := range tb.fld_values {
		if index > 0 {
			buffer1.WriteString(",")
			buffer2.WriteString(",")
		}

		buffer1.WriteString(fld.Key)
		if fld.Val == nil {
			buffer2.WriteString("null")
		} else if exp, ok := fld.Val.(SqlExp); ok {
			buffer2.WriteString(exp.Tpl_sql)
			if exp.Values != nil {
				vals = append(vals, exp.Values...)
			}
		} else {
			buffer2.WriteString("?")
			vals = append(vals, fld.Val)
		}

		index += 1
	}
	buffer1.WriteString(")")
	buffer2.WriteString(")")

	//mssql not support LastInsertId
	//use output and queryrrow to get LastInsertId
	if len(tb.output_str) > 0 {
		buffer1.WriteString(" OUTPUT ")
		buffer1.WriteString(tb.output_str)
		buffer1.WriteString(" ")
	} else {
		buffer1.WriteString(" OUTPUT 0 ")
	}

	buffer1.Write(buffer2.Bytes())
	return buffer1.String(), vals
}

func (db *dialectMssql) GenTableGetSql(tb *DbTable) string {
	var buffer bytes.Buffer

	buffer.WriteString("select top 1  ")
	buffer.WriteString(tb.select_str)
	buffer.WriteString(" from ")
	buffer.WriteString(tb.table_name)
	if len(tb.table_alias) > 0 {
		buffer.WriteString(" ")
		buffer.WriteString(tb.table_alias)
	}

	if len(tb.join_str) > 0 {
		buffer.WriteString(tb.join_str)
	}

	if len(tb.db_where.Tpl_sql) > 0 {
		buffer.WriteString(" where ")
		buffer.WriteString(tb.db_where.Tpl_sql)
	}

	if len(tb.group_by) > 0 {
		buffer.WriteString(" group by ")
		buffer.WriteString(tb.group_by)
	}

	if len(tb.db_having.Tpl_sql) > 0 {
		buffer.WriteString(" having ")
		buffer.WriteString(tb.db_having.Tpl_sql)
	}

	if len(tb.order_by) > 0 {
		buffer.WriteString(" order by ")
		buffer.WriteString(tb.order_by)
	}

	return buffer.String()
}

func (db *dialectMssql) GenTableFindSql(tb *DbTable) string {
	var buffer bytes.Buffer

	buffer.WriteString("select ")
	if tb.db_limit > 0 {
		buffer.WriteString(fmt.Sprintf("top %d ", tb.db_limit))
	}
	buffer.WriteString(tb.select_str)
	buffer.WriteString(" from ")
	buffer.WriteString(tb.table_name)
	if len(tb.table_alias) > 0 {
		buffer.WriteString(" ")
		buffer.WriteString(tb.table_alias)
	}

	if len(tb.join_str) > 0 {
		buffer.WriteString(tb.join_str)
	}

	if len(tb.db_where.Tpl_sql) > 0 {
		buffer.WriteString(" where ")
		buffer.WriteString(tb.db_where.Tpl_sql)
	}

	if len(tb.group_by) > 0 {
		buffer.WriteString(" group by ")
		buffer.WriteString(tb.group_by)
	}

	if len(tb.db_having.Tpl_sql) > 0 {
		buffer.WriteString(" having ")
		buffer.WriteString(tb.db_having.Tpl_sql)
	}

	if len(tb.order_by) > 0 {
		buffer.WriteString(" order by ")
		buffer.WriteString(tb.order_by)
	}

	return buffer.String()
}
