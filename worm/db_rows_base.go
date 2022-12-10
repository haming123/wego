package worm

import "database/sql"

type DbRows struct {
	*sql.Rows
}

func (rows *DbRows) Scan(dest ...interface{}) error {
	return rows_scan(rows.Rows, dest...)
}

//执行行数据的scan
//将数据库查询的结果拷贝到desc对应的指针变量中
//在scan前将变量的指针包装为&FieldValue
//FieldValue实现了scanner接口用于接收数据库数据
//FieldValue能够处理字段为null的情况
func rows_scan(rows *sql.Rows, dest ...interface{}) error {
	values := make([]interface{}, len(dest))
	for i := 0; i < len(dest); i++ {
		fld := &FieldValue{"", dest[i], false}
		values[i] = fld
	}
	return rows.Scan(values...)
}

func Scan(rows DbRows, dest ...interface{}) error {
	return rows_scan(rows.Rows, dest...)
}
