package worm

import "database/sql"

type DbRows struct {
	*sql.Rows
}

func (rows *DbRows) Scan(dest ...interface{}) error {
	return rows_scan(rows.Rows, dest...)
}

func Scan(rows DbRows, dest ...interface{}) error {
	return rows_scan(rows.Rows, dest...)
}

func ScanModel(rows DbRows, ent_ptr interface{}) error {
	return scanModel(rows.Rows, ent_ptr)
}

func ScanStringRow(rows DbRows) (StringRow, error) {
	return scanStringRow(rows.Rows)
}
