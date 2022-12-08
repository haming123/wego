package worm

import (
	"bytes"
	"database/sql"
	"fmt"
)

/*
不同的数据库中，SQL语句使用的占位符语法不尽相同。
MySQL	?
SQLServer	?
PostgreSQL	$1, $2等
SQLite	? 和$1
Oracle	:name
*/

type ColumnInfo struct {
	Name            string
	SQLType         string
	DbType          string
	Comment         string
	Length          int
	Length2         int
	Nullable        bool
	IsPrimaryKey    bool
	IsAutoIncrement bool
}

type Dialect interface {
	GetName() string
	DbType2GoType(colType string) string
	GetColumns(db_raw *sql.DB, tableName string) ([]ColumnInfo, error)
	LimitSql(offset int64, limit int64) string
	ParsePlaceholder(sql_tpl string) string

	ModelInsertHasOutput(md *DbModel) bool
	GenModelInsertSql(md *DbModel) string
	GenModelUpdateSql(md *DbModel) string
	GenModelDeleteSql(md *DbModel) string
	GenModelGetSql(md *DbModel) string
	GenModelFindSql(md *DbModel) string

	TableInsertHasOutput(tb *DbTable) bool
	GenTableInsertSql(tb *DbTable) (string, []interface{})
	GenTableUpdateSql(tb *DbTable) (string, []interface{})
	GenTableDeleteSql(tb *DbTable) string
	GenTableGetSql(tb *DbTable) string
	GenTableFindSql(tb *DbTable) string

	GenJointGetSql(lk *DbJoint) string
	GenJointFindSql(lk *DbJoint) string
}

type DialectBase struct {
}

func (db *DialectBase) GetName() string {
	return ""
}

func (db *DialectBase) DbType2GoType(colType string) string {
	return colType
}

func (db *DialectBase) GetColumns(db_raw *sql.DB, tableName string) ([]ColumnInfo, error) {
	return nil, nil
}

func (db *DialectBase) LimitSql(offset int64, limit int64) string {
	return ""
}

func (db *DialectBase) ParsePlaceholder(sql_tpl string) string {
	tpl_str := sql_tpl
	return tpl_str
}

func (db *DialectBase) ModelInsertHasOutput(md *DbModel) bool {
	return false
}

func (db *DialectBase) TableInsertHasOutput(tb *DbTable) bool {
	return false
}

func (db *DialectBase) GenModelInsertSql(md *DbModel) string {
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

func (db *DialectBase) GenModelUpdateSql(md *DbModel) string {
	var buffer bytes.Buffer
	buffer.WriteString("update ")
	buffer.WriteString(md.table_name)
	buffer.WriteString(" set ")
	index := 0
	for i, item := range md.flds_addr {
		if md.GetFieldFlag4Update(i) == false {
			continue
		}
		if index > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(item.FName)
		buffer.WriteString("=?")
		index += 1
	}

	if len(md.db_where.Tpl_sql) > 0 {
		buffer.WriteString(" where ")
		buffer.WriteString(md.db_where.Tpl_sql)
	}

	return buffer.String()
}

func (db *DialectBase) GenModelDeleteSql(md *DbModel) string {
	var buffer bytes.Buffer
	buffer.WriteString("delete from ")
	buffer.WriteString(md.table_name)
	if len(md.db_where.Tpl_sql) > 0 {
		buffer.WriteString(" where ")
		buffer.WriteString(md.db_where.Tpl_sql)
	}
	return buffer.String()
}

func (db *DialectBase) GenModelGetSql(md *DbModel) string {
	var buffer bytes.Buffer

	buffer.WriteString("select ")
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

	buffer.WriteString(" limit 1")
	return buffer.String()
}

func (db *DialectBase) GenModelFindSql(md *DbModel) string {
	var buffer bytes.Buffer

	buffer.WriteString("select ")
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

	if md.db_limit > 0 {
		dialect := md.db_ptr.engine.db_dialect
		str_val := dialect.LimitSql(md.db_offset, md.db_limit)
		buffer.WriteString(str_val)
	}

	return buffer.String()
}

func (db *DialectBase) GenJointGetSql(lk *DbJoint) string {
	var buffer bytes.Buffer

	index := 0
	buffer.WriteString("select ")
	for _, table := range lk.tables {
		str := table.gen_select_fields()
		if len(str) < 1 {
			continue
		}
		if index > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(table.gen_select_fields())
		index += 1
	}

	buffer.WriteString(" from ")
	for index, table := range lk.tables {
		if index > 0 {
			buffer.WriteString(" ")
			buffer.WriteString(get_join_type_str(table.join_type))
			buffer.WriteString(" ")
		}
		buffer.WriteString(table.table_name)
		if len(table.table_alias) > 0 {
			buffer.WriteString(" ")
			buffer.WriteString(table.table_alias)
		}
		if len(table.join_on) > 0 {
			buffer.WriteString(" on ")
			buffer.WriteString(table.join_on)
		}
	}

	if len(lk.db_where.Tpl_sql) > 0 {
		buffer.WriteString(" where ")
		buffer.WriteString(lk.db_where.Tpl_sql)
	}

	if len(lk.order_by) > 0 {
		buffer.WriteString(" order by ")
		buffer.WriteString(lk.order_by)
	}

	buffer.WriteString(" limit 1")
	return buffer.String()
}

func (db *DialectBase) GenJointFindSql(lk *DbJoint) string {
	var buffer bytes.Buffer

	index := 0
	buffer.WriteString("select ")
	for _, table := range lk.tables {
		str := table.gen_select_fields()
		if len(str) < 1 {
			continue
		}
		if index > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(table.gen_select_fields())
		index += 1
	}

	buffer.WriteString(" from ")
	for index, table := range lk.tables {
		if index > 0 {
			buffer.WriteString(" ")
			buffer.WriteString(get_join_type_str(table.join_type))
			buffer.WriteString(" ")
		}
		buffer.WriteString(table.table_name)
		if len(table.table_alias) > 0 {
			buffer.WriteString(" ")
			buffer.WriteString(table.table_alias)
		}
		if len(table.join_on) > 0 {
			buffer.WriteString(" on ")
			buffer.WriteString(table.join_on)
		}
	}

	if len(lk.db_where.Tpl_sql) > 0 {
		buffer.WriteString(" where ")
		buffer.WriteString(lk.db_where.Tpl_sql)
	}

	if len(lk.order_by) > 0 {
		buffer.WriteString(" order by ")
		buffer.WriteString(lk.order_by)
	}

	if lk.db_limit > 0 {
		dialect := lk.db_ptr.engine.db_dialect
		str_val := dialect.LimitSql(lk.db_offset, lk.db_limit)
		buffer.WriteString(str_val)
	}

	return buffer.String()
}

//生成insert sql语句
//sql语句与values数组必须在一个循环中生成，因为map多次遍历时次序可能不同
func (db *DialectBase) GenTableInsertSql(tb *DbTable) (string, []interface{}) {
	index := 0
	vals := []interface{}{}

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
	return buffer1.String(), vals
}

//生成 update sql语句
//sql语句与values数组必须在一个循环中生成，因为map多次遍历时次序可能不同
func (db *DialectBase) GenTableUpdateSql(tb *DbTable) (string, []interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString("update ")
	buffer.WriteString(tb.table_name)
	buffer.WriteString(" set ")

	index := 0
	vals := []interface{}{}
	for name, val := range tb.fld_values {
		if index > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(name)
		if val == nil {
			buffer.WriteString("=null")
		} else if exp, ok := val.(SqlExp); ok {
			buffer.WriteString("=")
			buffer.WriteString(exp.Tpl_sql)
			if exp.Values != nil {
				vals = append(vals, exp.Values...)
			}
		} else {
			buffer.WriteString("=?")
			vals = append(vals, val)
		}
		index += 1
	}

	if len(tb.db_where.Tpl_sql) > 0 {
		buffer.WriteString(" where ")
		buffer.WriteString(tb.db_where.Tpl_sql)
	}

	return buffer.String(), vals
}

func (db *DialectBase) GenTableDeleteSql(tb *DbTable) string {
	sql_str := fmt.Sprintf("delete from %s", tb.table_name)
	if len(tb.db_where.Tpl_sql) > 0 {
		sql_str += " where " + tb.db_where.Tpl_sql
	}
	return sql_str
}

func (db *DialectBase) GenTableGetSql(tb *DbTable) string {
	var buffer bytes.Buffer

	buffer.WriteString("select ")
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

	buffer.WriteString(" limit 1")
	return buffer.String()
}

func (db *DialectBase) GenTableFindSql(tb *DbTable) string {
	var buffer bytes.Buffer

	buffer.WriteString("select ")
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

	if tb.db_limit > 0 {
		dialect := tb.db_ptr.engine.db_dialect
		str_val := dialect.LimitSql(tb.db_offset, tb.db_limit)
		buffer.WriteString(str_val)
	}

	return buffer.String()
}
