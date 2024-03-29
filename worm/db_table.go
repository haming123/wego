package worm

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	SQL_TYPE_INS int = 1
	SQL_TYPE_UPD int = 2
	SQL_TYPE_DEL int = 3
	SQL_TYPE_SEL int = 4
)

type KeyVal struct {
	Key string
	Val interface{}
}
type ValueMap map[string]interface{}
type ValueArray []KeyVal

type SqlExp struct {
	Tpl_sql string
	Values  []interface{}
}

type DbTable struct {
	SqlContex
	sql_type    int
	db_ptr      *DbSession
	fld_values  ValueArray
	table_name  string
	table_alias string
	select_str  string
	output_str  string
	return_str  string
	join_str    string
	db_where    DbWhere
	group_by    string
	db_having   DbWhere
	order_by    string
	db_limit    int64
	db_offset   int64
	lock_str    string
	ctx         context.Context
	sql_err     error
}

func NewDbTable(dbs *DbSession, table_name string) *DbTable {
	tb := &DbTable{}
	tb.db_ptr = dbs
	tb.table_name = table_name
	return tb
}

func (tb *DbTable) GetContext() context.Context {
	return tb.ctx
}

func (tb *DbTable) Context(ctx context.Context) *DbTable {
	tb.ctx = ctx
	return tb
}

func (tb *DbTable) UsePrepare(val bool) *DbTable {
	tb.use_prepare.Valid = true
	tb.use_prepare.Bool = val
	return tb
}

func (tb *DbTable) ShowLog(val bool) *DbTable {
	tb.show_log.Valid = true
	tb.show_log.Bool = val
	return tb
}

func (tb *DbTable) UseMaster(val bool) *DbTable {
	tb.use_master.Valid = true
	tb.use_master.Bool = val
	return tb
}

func (tb *DbTable) Select(fields ...string) *DbTable {
	tb.select_str = strings.Join(fields, ",")
	return tb
}

func (tb *DbTable) Output(output string) *DbTable {
	tb.output_str = output
	return tb
}

func (tb *DbTable) ReturnIng(returning string) *DbTable {
	tb.return_str = returning
	return tb
}

func (tb *DbTable) Alias(alias string) *DbTable {
	tb.table_alias = alias
	return tb
}

func (tb *DbTable) Value(col_name string, val interface{}) *DbTable {
	var item KeyVal
	item.Key = col_name
	item.Val = val

	if tb.fld_values == nil {
		tb.fld_values = make([]KeyVal, 0, 10)
	}
	tb.fld_values = append(tb.fld_values, item)
	return tb
}

func (tb *DbTable) Values(map_data ValueMap) *DbTable {
	if tb.fld_values == nil {
		tb.fld_values = make([]KeyVal, 0, len(map_data))
	}
	for k, v := range map_data {
		var item KeyVal
		item.Key = k
		item.Val = v
		tb.fld_values = append(tb.fld_values, item)
	}
	return tb
}

func (tb *DbTable) SetWhere(sqlw *DbWhere) *DbTable {
	tb.db_where.Init(sqlw.Tpl_sql, sqlw.Values...)
	return tb
}
func (tb *DbTable) Where(sql string, vals ...interface{}) *DbTable {
	tb.db_where.Init(sql, vals...)
	return tb
}
func (tb *DbTable) WhereIn(sql string, vals ...interface{}) *DbTable {
	tb.db_where.Reset()
	tb.db_where.AndIn(sql, vals...)
	return tb
}
func (tb *DbTable) WhereNotIn(sql string, vals ...interface{}) *DbTable {
	tb.db_where.Reset()
	tb.db_where.AndNotIn(sql, vals...)
	return tb
}
func (tb *DbTable) WhereIf(cond bool, sql string, vals ...interface{}) *DbTable {
	if cond {
		tb.db_where.Init(sql, vals...)
	}
	return tb
}
func (tb *DbTable) ID(val int64) *DbTable {
	tb.db_where.Init("id=?", val)
	return tb
}
func (tb *DbTable) And(sql string, vals ...interface{}) *DbTable {
	tb.db_where.And(sql, vals...)
	return tb
}
func (tb *DbTable) Or(sql string, vals ...interface{}) *DbTable {
	tb.db_where.Or(sql, vals...)
	return tb
}
func (tb *DbTable) AndIf(cond bool, sql string, vals ...interface{}) *DbTable {
	tb.db_where.AndIf(cond, sql, vals...)
	return tb
}
func (tb *DbTable) OrIf(cond bool, sql string, vals ...interface{}) *DbTable {
	tb.db_where.OrIf(cond, sql, vals...)
	return tb
}
func (tb *DbTable) AndExp(sqlw_sub *DbWhere) *DbTable {
	tb.db_where.AndExp(sqlw_sub)
	return tb
}
func (tb *DbTable) OrExp(sqlw_sub *DbWhere) *DbTable {
	tb.db_where.OrExp(sqlw_sub)
	return tb
}
func (tb *DbTable) AndIn(sql string, vals ...interface{}) *DbTable {
	tb.db_where.AndIn(sql, vals...)
	return tb
}
func (tb *DbTable) AndNotIn(sql string, vals ...interface{}) *DbTable {
	tb.db_where.AndNotIn(sql, vals...)
	return tb
}
func (tb *DbTable) OrIn(sql string, vals ...interface{}) *DbTable {
	tb.db_where.AndIn(sql, vals...)
	return tb
}
func (tb *DbTable) OrNotIn(sql string, vals ...interface{}) *DbTable {
	tb.db_where.OrNotIn(sql, vals...)
	return tb
}

func (tb *DbTable) Join(table string, alias string, join_on string) *DbTable {
	str := fmt.Sprintf("%s join %s %s on %s", tb.join_str, table, alias, join_on)
	if len(alias) < 1 {
		str = fmt.Sprintf("%s join %s on %s", tb.join_str, table, join_on)
	}
	tb.join_str = str
	return tb
}

func (tb *DbTable) LeftJoin(table string, alias string, join_on string) *DbTable {
	str := fmt.Sprintf("%s left join %s %s on %s", tb.join_str, table, alias, join_on)
	if len(alias) < 1 {
		str = fmt.Sprintf("%s left join %s on %s", tb.join_str, table, join_on)
	}
	tb.join_str = str
	return tb
}

func (tb *DbTable) RightJoin(table string, alias string, join_on string) *DbTable {
	str := fmt.Sprintf("%s right join %s %s on %s", tb.join_str, table, alias, join_on)
	if len(alias) < 1 {
		str = fmt.Sprintf("%s right join %s on %s", tb.join_str, table, join_on)
	}
	tb.join_str = str
	return tb
}

func (tb *DbTable) Top(rows int64) *DbTable {
	tb.db_limit = rows
	return tb
}

func (tb *DbTable) Limit(rows int64) *DbTable {
	tb.db_limit = rows
	return tb
}

func (tb *DbTable) Offset(offset int64) *DbTable {
	tb.db_offset = offset
	return tb
}

func (tb *DbTable) Page(rows int64, page_no int64) *DbTable {
	tb.db_offset = page_no * rows
	tb.db_limit = rows
	return tb
}

func (tb *DbTable) Lock(val string) *DbTable {
	tb.lock_str = val
	return tb
}

func (tb *DbTable) GroupBy(val string) *DbTable {
	tb.group_by = val
	return tb
}

func (tb *DbTable) Having(sql string, vals ...interface{}) *DbTable {
	tb.db_having.Init(sql, vals...)
	return tb
}

func (tb *DbTable) OrderBy(val string) *DbTable {
	tb.order_by = val
	return tb
}

func (tb *DbTable) insertWithOutput() (int64, error) {
	sql_str, vals := tb.db_ptr.engine.db_dialect.GenTableInsertSql(tb)
	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, sql_str, vals...)
	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		rows.Close()
		return 0, nil
	}

	var id int64 = 0
	err = rows.Scan(&id)
	if err != nil {
		rows.Close()
		return 0, err
	}
	rows.Close()
	return id, nil
}

func (tb *DbTable) Insert() (int64, error) {
	if tb.db_ptr.engine.db_dialect.TableInsertHasOutput(tb) {
		return tb.insertWithOutput()
	}

	sql_str, vals := tb.db_ptr.engine.db_dialect.GenTableInsertSql(tb)
	res, err := tb.db_ptr.ExecSQL(&tb.SqlContex, sql_str, vals...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (tb *DbTable) Update() (int64, error) {
	if len(tb.db_where.Tpl_sql) < 1 {
		return 0, errors.New("no where clause")
	}

	sql_str, vals := tb.db_ptr.engine.db_dialect.GenTableUpdateSql(tb)
	vals = append(vals, tb.db_where.Values...)
	res, err := tb.db_ptr.ExecSQL(&tb.SqlContex, sql_str, vals...)
	if err != nil {
		return 0, err
	}

	num, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (tb *DbTable) Delete() (int64, error) {
	if len(tb.db_where.Tpl_sql) < 1 {
		return 0, errors.New("no where clause")
	}

	sql_str := tb.db_ptr.engine.db_dialect.GenTableDeleteSql(tb)
	vals := append([]interface{}{}, tb.db_where.Values...)
	res, err := tb.db_ptr.ExecSQL(&tb.SqlContex, sql_str, vals...)
	if err != nil {
		return 0, err
	}

	num, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (tb *DbTable) getRows() (*sql.Rows, error) {
	sql_str := tb.db_ptr.engine.db_dialect.GenTableFindSql(tb)
	vals := []interface{}{}
	vals = append(vals, tb.db_where.Values...)
	vals = append(vals, tb.db_having.Values...)
	return tb.db_ptr.ExecQuery(&tb.SqlContex, sql_str, vals...)
}

func (tb *DbTable) SqlRows() (*sql.Rows, error) {
	return tb.getRows()
}

func (tb *DbTable) Rows() (DbRows, error) {
	var ret DbRows
	rows, err := tb.getRows()
	if err != nil {
		return ret, err
	}
	ret.Rows = rows
	return ret, nil
}

func (tb *DbTable) ModelRows() (StructRows, error) {
	var ret StructRows
	rows, err := tb.getRows()
	if err != nil {
		return ret, err
	}
	ret.Rows = rows
	return ret, nil
}

func (tb *DbTable) StringRows() (StringRows, error) {
	var ret StringRows
	rows, err := tb.getRows()
	if err != nil {
		return ret, err
	}
	ret.Rows = rows
	return ret, nil
}

func (tb *DbTable) Exist() (bool, error) {
	rows, err := tb.getRows()
	if err != nil {
		return false, err
	}

	if !rows.Next() {
		rows.Close()
		return false, nil
	}

	rows.Close()
	return true, nil
}

func (tb *DbTable) Get(arg ...interface{}) (bool, error) {
	rows, err := tb.getRows()
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		rows.Close()
		return false, nil
	}

	//err = rows.Scan(arg...)
	err = rows_scan(rows, arg...)
	if err != nil {
		rows.Close()
		return false, err
	}

	rows.Close()
	return true, nil
}

func (tb *DbTable) GetValues(arg ...interface{}) (bool, error) {
	return tb.Get(arg...)
}

func (tb *DbTable) GetInt() (sql.NullInt64, error) {
	var val sql.NullInt64
	fld := FieldValue{"", &val, false}
	has, err := tb.Get(&fld)
	val.Valid = has
	if err != nil {
		return val, err
	}
	return val, nil
}

func (tb *DbTable) GetFlaot() (sql.NullFloat64, error) {
	return tb.GetFloat()
}

func (tb *DbTable) GetFloat() (sql.NullFloat64, error) {
	var val sql.NullFloat64
	fld := FieldValue{"", &val, false}
	has, err := tb.Get(&fld)
	val.Valid = has
	if err != nil {
		return val, err
	}
	return val, nil
}

func (tb *DbTable) GetString() (sql.NullString, error) {
	var val sql.NullString
	fld := FieldValue{"", &val, false}
	has, err := tb.Get(&fld)
	val.Valid = has
	if err != nil {
		return val, err
	}
	return val, nil
}

func (tb *DbTable) GetTime() (sql.NullTime, error) {
	var val sql.NullTime
	fld := FieldValue{"", &val, false}
	has, err := tb.Get(&fld)
	val.Valid = has
	if err != nil {
		return val, err
	}
	return val, nil
}

func (tb *DbTable) GetModel(ent_ptr interface{}) (bool, error) {
	rows, err := tb.getRows()
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		rows.Close()
		return false, nil
	}
	err = scanModel(rows, ent_ptr)
	if err != nil {
		rows.Close()
		return false, err
	}
	rows.Close()
	return true, nil
}

func (tb *DbTable) GetRow() (StringRow, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		rows.Close()
		return nil, nil
	}
	ret, err := scanStringRow(rows)
	if err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()
	return ret, nil
}

func (tb *DbTable) FindValues(arr_ptr_arr ...interface{}) (int, error) {
	rows, err := tb.getRows()
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

func (tb *DbTable) FindInt() ([]int64, error) {
	rows, err := tb.getRows()
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

func (tb *DbTable) FindFloat() ([]float64, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}

	var arr []float64
	var val float64 = 0
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

func (tb *DbTable) FindString() ([]string, error) {
	rows, err := tb.getRows()
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

func (tb *DbTable) FindTime() ([]time.Time, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}

	var arr []time.Time
	val := time.Time{}
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

func (tb *DbTable) FindIntInt() ([]KeyVal4IntInt, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}

	var arr []KeyVal4IntInt
	val := KeyVal4IntInt{0, 0}
	fld1 := FieldValue{"", &val.Key, false}
	fld2 := FieldValue{"", &val.Val, false}
	for rows.Next() {
		err = rows.Scan(&fld1, &fld2)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	rows.Close()
	return arr, nil
}

func (tb *DbTable) FindIntFloat() ([]KeyVal4IntFloat, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}

	var arr []KeyVal4IntFloat
	val := KeyVal4IntFloat{0, 0}
	fld1 := FieldValue{"", &val.Key, false}
	fld2 := FieldValue{"", &val.Val, false}
	for rows.Next() {
		err = rows.Scan(&fld1, &fld2)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	rows.Close()
	return arr, nil
}

func (tb *DbTable) FindIntString() ([]KeyVal4IntString, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}

	var arr []KeyVal4IntString
	val := KeyVal4IntString{0, ""}
	fld1 := FieldValue{"", &val.Key, false}
	fld2 := FieldValue{"", &val.Val, false}
	for rows.Next() {
		err = rows.Scan(&fld1, &fld2)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	rows.Close()
	return arr, nil
}

func (tb *DbTable) FindIntTime() ([]KeyVal4IntTime, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}

	var arr []KeyVal4IntTime
	val := KeyVal4IntTime{0, time.Time{}}
	fld1 := FieldValue{"", &val.Key, false}
	fld2 := FieldValue{"", &val.Val, false}
	for rows.Next() {
		err = rows.Scan(&fld1, &fld2)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	rows.Close()
	return arr, nil
}

func (tb *DbTable) FindStringInt() ([]KeyVal4StringInt, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}

	var arr []KeyVal4StringInt
	val := KeyVal4StringInt{"", 0}
	fld1 := FieldValue{"", &val.Key, false}
	fld2 := FieldValue{"", &val.Val, false}
	for rows.Next() {
		err = rows.Scan(&fld1, &fld2)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	rows.Close()
	return arr, nil
}

func (tb *DbTable) FindStringFloat() ([]KeyVal4StringFloat, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}

	var arr []KeyVal4StringFloat
	val := KeyVal4StringFloat{"", 0}
	fld1 := FieldValue{"", &val.Key, false}
	fld2 := FieldValue{"", &val.Val, false}
	for rows.Next() {
		err = rows.Scan(&fld1, &fld2)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	rows.Close()
	return arr, nil
}

func (tb *DbTable) FindStringString() ([]KeyVal4StringString, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}

	var arr []KeyVal4StringString
	val := KeyVal4StringString{"", ""}
	fld1 := FieldValue{"", &val.Key, false}
	fld2 := FieldValue{"", &val.Val, false}
	for rows.Next() {
		err = rows.Scan(&fld1, &fld2)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	rows.Close()
	return arr, nil
}

func (tb *DbTable) FindStringTime() ([]KeyVal4StringTime, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}

	var arr []KeyVal4StringTime
	val := KeyVal4StringTime{"", time.Time{}}
	fld1 := FieldValue{"", &val.Key, false}
	fld2 := FieldValue{"", &val.Val, false}
	for rows.Next() {
		err = rows.Scan(&fld1, &fld2)
		if err != nil {
			return arr, err
		}
		arr = append(arr, val)
	}

	rows.Close()
	return arr, nil
}

func (tb *DbTable) FindModel(arr_ptr interface{}) error {
	rows, err := tb.getRows()
	if err != nil {
		return err
	}
	err = scanModelArray(rows, arr_ptr)
	rows.Close()
	return err
}

func (tb *DbTable) FindRow() (*StringTable, error) {
	rows, err := tb.getRows()
	if err != nil {
		return nil, err
	}
	ret, err := scanStringTable(rows)
	rows.Close()
	return ret, err
}

func (tb *DbTable) gen_count_sql(count_field string) string {
	var buffer bytes.Buffer

	buffer.WriteString("select ")
	buffer.WriteString(count_field)
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

	return buffer.String()
}

func (tb *DbTable) Count(field ...string) (int64, error) {
	if len(field) > 1 {
		return 0, errors.New("field vumber > 0")
	}

	count_field := "count(1)"
	if len(field) == 1 {
		count_field = fmt.Sprintf("count(%s)", field[0])
	}

	sql_str := tb.gen_count_sql(count_field)
	if len(tb.group_by) > 0 {
		sub_sql := tb.db_ptr.engine.db_dialect.GenTableFindSql(tb)
		sql_str = fmt.Sprintf("select %s from (%s) tmp", count_field, sub_sql)
	}

	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, sql_str, tb.db_where.Values...)
	if err != nil {
		return 0, err
	}
	if !rows.Next() {
		rows.Close()
		return 0, nil
	}

	var total int64
	err = rows_scan(rows, &total)
	rows.Close()
	return total, nil
}

func (tb *DbTable) DistinctCount(field string) (int64, error) {
	count_field := fmt.Sprintf("count(distinct %s)", field)
	sql_str := tb.gen_count_sql(count_field)
	if len(tb.group_by) > 0 {
		sub_sql := tb.db_ptr.engine.db_dialect.GenTableFindSql(tb)
		sql_str = fmt.Sprintf("select %s from (%s) tmp", count_field, sub_sql)
	}

	rows, err := tb.db_ptr.ExecQuery(&tb.SqlContex, sql_str, tb.db_where.Values...)
	if err != nil {
		return 0, err
	}
	if !rows.Next() {
		rows.Close()
		return 0, nil
	}

	var total int64
	err = rows_scan(rows, &total)
	rows.Close()
	return total, nil
}
