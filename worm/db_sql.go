package worm

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
)

type DbSQL struct {
	SqlContex
	db_ptr  *DbSession
	sql_tpl string
	values  []interface{}
	ctx     context.Context
}

func NewDbSQL(dbs *DbSession, sql_str string, args ...interface{}) *DbSQL {
	tb := &DbSQL{}
	tb.db_ptr = dbs
	tb.sql_tpl = sql_str
	tb.values = args
	return tb
}

func (tb *DbSQL) GetContext() context.Context {
	return tb.ctx
}

func (tb *DbSQL) Context(ctx context.Context) *DbSQL {
	tb.ctx = ctx
	return tb
}

func (tb *DbSQL) UsePrepare(val bool) *DbSQL {
	tb.use_prepare.Valid = true
	tb.use_prepare.Bool = val
	return tb
}

func (tb *DbSQL) ShowLog(val bool) *DbSQL {
	tb.show_log.Valid = true
	tb.show_log.Bool = val
	return tb
}

func (tb *DbSQL) UseMaster(val bool) *DbSQL {
	tb.use_master.Valid = true
	tb.use_master.Bool = val
	return tb
}

func (tb *DbSQL) Exec() (sql.Result, error) {
	return tb.db_ptr.ExecSQL(&tb.SqlContex, tb.sql_tpl, tb.values...)
}

func (tb *DbSQL) Rows() (*sql.Rows, error) {
	return tb.db_ptr.ExecQuery(&tb.SqlContex, tb.sql_tpl, tb.values...)
}

func (tb *DbSQL) Get(arg ...interface{}) (bool, error) {
	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, tb.sql_tpl, tb.values...)
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		rows.Close()
		return false, nil
	}

	//err = rows.Scan(arg...)
	err = Scan(rows, arg...)
	if err != nil {
		rows.Close()
		return false, err
	}

	rows.Close()
	return true, nil
}

func (tb *DbSQL) GetValues(arg ...interface{}) (bool, error) {
	return tb.Get(arg...)
}

func (tb *DbSQL) GetInt() (sql.NullInt64, error) {
	var val sql.NullInt64
	fld := FieldValue{"", &val, false}
	has, err := tb.Get(&fld)
	val.Valid = has
	if err != nil {
		return val, err
	}
	return val, nil
}

func (tb *DbSQL) GetFloat() (sql.NullFloat64, error) {
	var val sql.NullFloat64
	fld := FieldValue{"", &val, false}
	has, err := tb.Get(&fld)
	val.Valid = has
	if err != nil {
		return val, err
	}
	return val, nil
}

func (tb *DbSQL) GetString() (sql.NullString, error) {
	var val sql.NullString
	fld := FieldValue{"", &val, false}
	has, err := tb.Get(&fld)
	val.Valid = has
	if err != nil {
		return val, err
	}
	return val, nil
}

func (tb *DbSQL) GetTime() (sql.NullTime, error) {
	var val sql.NullTime
	fld := FieldValue{"", &val, false}
	has, err := tb.Get(&fld)
	val.Valid = has
	if err != nil {
		return val, err
	}
	return val, nil
}

func (tb *DbSQL) GetModel(ent_ptr interface{}) (bool, error) {
	v_ent := reflect.ValueOf(ent_ptr)
	if v_ent.Kind() != reflect.Ptr {
		return false, errors.New("ent_ptr must be Pointer")
	}
	if v_ent.IsNil() {
		return false, errors.New("ent_ptr is nil")
	}
	t_ent_base := GetDirectType(v_ent.Type())
	if t_ent_base.Kind() != reflect.Struct {
		return false, errors.New("ent_ptr muse be Struct")
	}

	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, tb.sql_tpl, tb.values...)
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		rows.Close()
		return false, nil
	}

	err = ScanModel(rows, ent_ptr)
	if err != nil {
		rows.Close()
		return false, err
	}

	rows.Close()
	return true, nil
}

func (tb *DbSQL) GetRow() (StringRow, error) {
	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, tb.sql_tpl, tb.values...)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		rows.Close()
		return nil, nil
	}

	ret, err := ScanStringRow(rows)
	if err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()
	return ret, nil
}

func (tb *DbSQL) FindValues(arr_ptr_arr ...interface{}) (int, error) {
	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, tb.sql_tpl, tb.values...)
	if err != nil {
		return 0, err
	}

	num, err := findValues(rows, arr_ptr_arr...)
	if err != nil {
		return 0, err
	}
	rows.Close()
	return num, err
}

func (tb *DbSQL) FindInt() ([]int64, error) {
	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, tb.sql_tpl, tb.values...)
	if err != nil {
		return nil, err
	}

	var arr []int64
	var val int64 = 0
	fld := FieldValue{"", &val, false}
	for rows.Next() {
		err = rows.Scan(&fld)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	rows.Close()
	return arr, nil
}

func (tb *DbSQL) FindFloat() ([]float64, error) {
	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, tb.sql_tpl, tb.values...)
	if err != nil {
		return nil, err
	}

	var arr []float64
	var val float64 = 0.0
	fld := FieldValue{"", &val, false}
	for rows.Next() {
		err = rows.Scan(&fld)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	rows.Close()
	return arr, nil
}

func (tb *DbSQL) FindString() ([]string, error) {
	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, tb.sql_tpl, tb.values...)
	if err != nil {
		return nil, err
	}

	var arr []string
	var val string = ""
	fld := FieldValue{"", &val, false}
	for rows.Next() {
		err = rows.Scan(&fld)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	rows.Close()
	return arr, nil
}

func (tb *DbSQL) FindModel(arr_ptr interface{}) error {
	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, tb.sql_tpl, tb.values...)
	if err != nil {
		return err
	}

	err = ScanModelArray(rows, arr_ptr)
	rows.Close()
	return err
}

func (tb *DbSQL) FindRow() (*StringTable, error) {
	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, tb.sql_tpl, tb.values...)
	if err != nil {
		return nil, err
	}

	ret, err := ScanStringTable(rows)
	rows.Close()
	return ret, err
}
