package worm

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
)

const (
	JOIN_INNER    	int = 0
	JOIN_LJOIN     	int = 1
	JOIN_RJOIN     	int = 2
)

func get_join_type_str(join_type int) string {
	sql_str := "join"
	if join_type == JOIN_LJOIN {
		sql_str = "left join"
	} else if join_type == JOIN_RJOIN {
		sql_str = "right join"
	}
	return sql_str
}

type DbJoint struct {
	SqlContex
	db_ptr *DbSession
	md_ptr *DbModel
	md_arr []*DbModel
	db_where DbWhere
	order_by string
	db_limit int64
	db_offset int64
	Err error
	ctx context.Context
}

func NewJoint(dbs *DbSession, ent_ptr interface{}, alias string, fields ...string) *DbJoint {
	md := dbs.NewModel(ent_ptr, true)
	md.table_alias = alias
	if len(fields) > 0 {
		fields = parselTableSelect(alias, fields)
		if len(fields) > 0 {
			md.Select(fields...)
		}
	}

	lk := &DbJoint{}
	lk.db_ptr = dbs
	lk.md_ptr = md
	return lk
}

func (lk *DbJoint)GetContext() context.Context {
	return lk.ctx
}

func (lk *DbJoint)Context(ctx context.Context) *DbJoint {
	lk.ctx = ctx
	return lk
}

func (lk *DbJoint)UsePrepare(val bool) *DbJoint {
	lk.use_prepare.Valid = true
	lk.use_prepare.Bool = val
	return lk
}

func (lk *DbJoint)ShowLog(val bool) *DbJoint {
	lk.show_log.Valid = true
	lk.show_log.Bool = val
	return lk
}

func (lk *DbJoint)UseMaster(val bool) *DbJoint {
	lk.use_master.Valid = true
	lk.use_master.Bool = val
	return lk
}

func (lk *DbJoint)get_table_index(ent_type reflect.Type) int {
	index := -1
	num := len(lk.md_arr)
	for i:=0; i < num; i++ {
		t_table_ent := reflect.TypeOf(lk.md_arr[i].ent_ptr).Elem()
		if t_table_ent == ent_type {
			index = i
			break
		}
	}
	return index
}

//如果字段中存在"*"或"alias.*"则认为是选择全部字段，返回空数组
func parselTableSelect(alias string, fields []string) []string {
	if len(fields) < 1 {
		return fields
	}

	table_star := alias + ".*"
	for i:=0; i < len(fields); i++ {
		if fields[i] == "*" || fields[i] == table_star {
			fields = fields[:0]
			return fields
		}
		ind := strings.Index(fields[i], ".")
		if ind > 0 {
			fields[i] = fields[i][ind+1:]
		}
	}
	return fields
}

func(lk *DbJoint)Join(ent_ptr interface{}, alias string, join_on string, fields ...string) *DbJoint {
	md := lk.db_ptr.NewModel(ent_ptr, true)
	md.table_alias = alias
	md.join_type = JOIN_INNER
	md.join_on = join_on
	if len(fields) > 0 {
		fields = parselTableSelect(alias, fields)
		if len(fields) > 0 {
			md.Select(fields...)
		}
	}
	lk.md_arr = append(lk.md_arr, md)
	return lk
}

func(lk *DbJoint)LeftJoin(ent_ptr interface{}, alias string, join_on string, fields ...string) *DbJoint {
	md := lk.db_ptr.NewModel(ent_ptr, true)
	md.table_alias = alias
	md.join_type = JOIN_LJOIN
	md.join_on = join_on
	if len(fields) > 0 {
		fields = parselTableSelect(alias, fields)
		if len(fields) > 0 {
			md.Select(fields...)
		}
	}
	lk.md_arr = append(lk.md_arr, md)
	return lk
}

func(lk *DbJoint)RightJoin(ent_ptr interface{}, alias string, join_on string, fields ...string) *DbJoint {
	md := lk.db_ptr.NewModel(ent_ptr, true)
	md.table_alias = alias
	md.join_type = JOIN_RJOIN
	md.join_on = join_on
	if len(fields) > 0 {
		fields = parselTableSelect(alias, fields)
		if len(fields) > 0 {
			md.Select(fields...)
		}
	}
	lk.md_arr = append(lk.md_arr, md)
	return lk
}

func (lk *DbJoint)SetWhere(sqlw *DbWhere) *DbJoint {
	lk.db_where.Init(sqlw.Tpl_sql, sqlw.Values...)
	return lk
}
func (lk *DbJoint)Where(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.Init(sql, vals...)
	return lk
}
func (lk *DbJoint)WhereIf(cond bool, sql string, vals ...interface{}) *DbJoint {
	if cond {
		lk.db_where.Init(sql, vals...)
	}
	return lk
}
func (lk *DbJoint)WhereIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.Reset()
	lk.db_where.AndIn(sql, vals...)
	return lk
}
func (lk *DbJoint)WhereNotIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.Reset()
	lk.db_where.AndNotIn(sql, vals...)
	return lk
}
func (lk *DbJoint)And(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.And(sql, vals...)
	return lk
}
func (lk *DbJoint)AndIf(cond bool, sql string, vals ...interface{}) *DbJoint {
	lk.db_where.AndIf(cond, sql, vals...)
	return lk
}
func (lk *DbJoint)OrIf(cond bool, sql string, vals ...interface{}) *DbJoint {
	lk.db_where.OrIf(cond, sql, vals...)
	return lk
}
func (lk *DbJoint)Or(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.Or(sql, vals...)
	return lk
}
func (lk *DbJoint)AndExp(sqlw_sub *DbWhere) *DbJoint {
	lk.db_where.AndExp(sqlw_sub)
	return lk
}
func (lk *DbJoint)OrExp(sqlw_sub *DbWhere) *DbJoint {
	lk.db_where.OrExp(sqlw_sub)
	return lk
}
func (lk *DbJoint)AndIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.AndIn(sql, vals...)
	return lk
}
func (lk *DbJoint)AndNotIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.AndNotIn(sql, vals...)
	return lk
}
func (lk *DbJoint)OrIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.AndIn(sql, vals...)
	return lk
}
func (lk *DbJoint)OrNotIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.OrNotIn(sql, vals...)
	return lk
}
func (lk *DbJoint)OrderBy(val string) *DbJoint {
	lk.order_by = val
	return lk
}
func (lk *DbJoint)Top(rows int64) *DbJoint {
	lk.db_limit = rows
	return lk
}
func (lk *DbJoint)Limit(rows int64) *DbJoint {
	lk.db_limit = rows
	return lk
}
func (lk *DbJoint)Offset(offset int64) *DbJoint {
	lk.db_offset = offset
	return lk
}
func (lk *DbJoint)Page(rows int64, page_no int64) *DbJoint {
	lk.db_offset = page_no*rows
	lk.db_limit = rows
	return lk
}

func (lk *DbJoint)gen_select() string {
	var buffer bytes.Buffer

	buffer.WriteString("select ")
	buffer.WriteString(lk.md_ptr.gen_select_fields())
	for _, table := range  lk.md_arr {
		str := table.gen_select_fields()
		if len(str) < 1 {
			continue
		}
		buffer.WriteString(",")
		buffer.WriteString(table.gen_select_fields())
	}

	buffer.WriteString(" from ")
	buffer.WriteString(lk.md_ptr.table_name)
	if len(lk.md_ptr.table_alias) > 0 {
		buffer.WriteString(" ")
		buffer.WriteString(lk.md_ptr.table_alias)
	}

	for _, table := range lk.md_arr {
		buffer.WriteString(" ")
		buffer.WriteString(get_join_type_str(table.join_type))
		buffer.WriteString(" ")
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

	if len(lk.db_where.Tpl_sql)>0 {
		buffer.WriteString(" where ")
		buffer.WriteString(lk.db_where.Tpl_sql)
	}

	if len(lk.order_by) > 0 {
		buffer.WriteString(" order by ")
		buffer.WriteString(lk.order_by)
	}

	if lk.db_limit > 0 {
		//buffer.WriteString(" limit ")
		//buffer.WriteString(fmt.Sprintf("%d,%d", lk.db_offset, lk.db_limit))
		dialect := lk.db_ptr.engine.db_dialect
		str_val := dialect.LimitSql(lk.db_offset, lk.db_limit)
		buffer.WriteString(str_val)
	}

	return buffer.String()
}

func (lk *DbJoint)get_scan_valus() []interface{} {
	vals:= lk.md_ptr.get_scan_valus()
	for _, table := range lk.md_arr {
		vals = append(vals, table.get_scan_valus()...)
	}
	return vals
}

func call_after_query(md *DbModel) {
	hook, has_hook := md.ent_ptr.(AfterQueryInterface)
	if has_hook {
		hook.AfterQuery(md.ctx)
	}
}

func (lk *DbJoint)Scan() (bool, error) {
	if lk.Err != nil {
		return false, lk.Err
	}
	sql_str := lk.gen_select() + " limit 1"
	rows, err := lk.db_ptr.ExecQuery(&lk.SqlContex, sql_str, lk.db_where.Values...)
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		rows.Close()
		return false, nil
	}

	scan_vals:= lk.get_scan_valus()
	err = rows.Scan(scan_vals...)
	if err != nil {
		rows.Close()
		return false, err
	}

	call_after_query(lk.md_ptr)
	for _, table := range lk.md_arr {
		call_after_query(table)
	}

	rows.Close()
	return true, nil
}

//通过vo对象来选择需要查询的字段
func (lk *DbJoint)select_field_by_vo(vo_ptr VoLoader)  {
	//调用：LoadFromModel获取vo对应的字段
	//vo_ptr.LoadFromModel(lk.md_ptr, lk.md_ptr.ent_ptr)
	getSelectFieldsByVo(lk.md_ptr, vo_ptr)
	for _, table := range lk.md_arr {
		//vo_ptr.LoadFromModel(table, table.ent_ptr)
		getSelectFieldsByVo(table, vo_ptr)
	}
}

//通过model的类型来获取一个model地址
func (lk *DbJoint)get_model_by_ent_type(ent_type reflect.Type) *DbModel {
	if reflect.TypeOf(lk.md_ptr.ent_ptr).Elem() == ent_type {
		return lk.md_ptr
	}
	num := len(lk.md_arr)
	for i:=0; i < num; i++ {
		t_table_ent := reflect.TypeOf(lk.md_arr[i].ent_ptr).Elem()
		if t_table_ent == ent_type {
			return lk.md_arr[i]
		}
	}
	return nil
}

//通过名称以及类型获取model对象以及mode的字段的序号
func (lk *DbJoint)get_model_field_by_ent_type(fname string, ent_type reflect.Type) (*DbModel, int) {
	index := lk.md_ptr.get_field_index_byname(fname)
	if index >=0 && lk.md_ptr.flds_info[index].FieldType == ent_type {
		return lk.md_ptr, index
	}
	num := len(lk.md_arr)
	for i:=0; i < num; i++ {
		index := lk.md_arr[i].get_field_index_byname(fname)
		if index >=0 && lk.md_arr[i].flds_info[index].FieldType == ent_type {
			return lk.md_arr[i], index
		}
	}
	return nil, -1
}

//绑定scan地址到目标对象
func (lk *DbJoint)BindAddr2Struct(v_ent reflect.Value) {
	t_num := v_ent.NumField()
	for t:=0; t < t_num; t++ {
		v_field := v_ent.Field(t)
		t_field := v_field.Type()
		if t_field.Kind() == reflect.Struct {
			md := lk.get_model_by_ent_type(t_field)
			if md != nil {
				rebindEntAddrs(md.flds_info, v_field, md.flds_addr)
				continue
			}
		}
	}
}

func (lk *DbJoint)Get(args ...interface{}) (bool, error) {
	if lk.Err != nil {
		return false, lk.Err
	}

	if len(args) > 1 {
		return false, errors.New("arg number can not great 1")
	}

	//参数为空, 调用Scan()
	if len(args) < 1 {
		return lk.Scan()
	}

	ent_ptr := args[0]
	if ent_ptr == nil{
		return false, errors.New("ent_ptr must be *Struct")
	}
	//ent_ptr必须是一个指针
	v_ent := reflect.ValueOf(ent_ptr)
	if v_ent.Kind() != reflect.Ptr {
		return false, errors.New("ent_ptr must be *Struct")
	}
	//ent_ptr必须是一个结构体指针
	v_ent = reflect.Indirect(v_ent)
	if v_ent.Kind() != reflect.Struct {
		return false, errors.New("ent_ptr must be *Struct")
	}

	//若目标对象是一个vo，则通过vo来选择字段
	//若目标对象不是一个vo，则需要重新进行地址绑定，将scan的地址绑定到目标对象
	vo_ptr, isvo := ent_ptr.(VoLoader)
	if isvo {
		lk.select_field_by_vo(vo_ptr)
	} else {
		lk.BindAddr2Struct(v_ent)
	}

	sql_str := lk.gen_select() + " limit 1"
	rows, err := lk.db_ptr.ExecQuery(&lk.SqlContex, sql_str, lk.db_where.Values...)
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		rows.Close()
		return false, nil
	}

	scan_vals:= lk.get_scan_valus()
	err = rows.Scan(scan_vals...)
	if err != nil {
		rows.Close()
		return false, err
	}

	call_after_query(lk.md_ptr)
	for _, table := range lk.md_arr {
		call_after_query(table)
	}

	//若目标对象是一个vo，则调用LoadFromModel，给vo赋值
	if isvo {
		vo_ptr.LoadFromModel(nil, lk.md_ptr.ent_ptr)
		for _, table := range lk.md_arr {
			vo_ptr.LoadFromModel(nil, table.ent_ptr)
		}
	}

	rows.Close()
	return true, nil
}

func (lk *DbJoint)Find(arr_ptr interface{}) error {
	if lk.Err != nil {
		return lk.Err
	}

	v_arr := reflect.ValueOf(arr_ptr)
	if v_arr.Kind() != reflect.Ptr {
		return  errors.New("arr_ptr must be *Slice")
	}
	t_arr := GetDirectType(v_arr.Type())
	if t_arr.Kind() != reflect.Slice {
		return  errors.New("arr_ptr must be *Slice")
	}
	//获取数组成员的类型
	t_item := GetDirectType(t_arr.Elem())
	if t_item.Kind() != reflect.Struct {
		return errors.New("array item muse be Struct")
	}

	//若目标对象是一个vo，则通过vo来选择字段
	//若目标对象不是一个vo，则需要重新进行地址绑定，将scan的地址绑定到目标对象
	v_item := reflect.New(t_item)
	v_item_base := v_item.Elem()
	vo_ptr, isvo := v_item.Interface().(VoLoader)
	if isvo {
		lk.select_field_by_vo(vo_ptr)
	} else {
		lk.BindAddr2Struct(v_item_base)
	}

	sql_str := lk.gen_select()
	vals:= lk.get_scan_valus()
	rows, err := lk.db_ptr.ExecQuery(&lk.SqlContex, sql_str, lk.db_where.Values...)
	if err != nil {
		return err
	}

	v_arr_base := reflect.Indirect(v_arr)
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			rows.Close()
			return err
		}

		call_after_query(lk.md_ptr)
		for _, table := range lk.md_arr {
			call_after_query(table)
		}

		//若目标对象是一个vo，则调用LoadFromModel，给vo赋值
		if isvo {
			vo_ptr.LoadFromModel(nil, lk.md_ptr.ent_ptr)
			for _, table := range lk.md_arr {
				vo_ptr.LoadFromModel(nil, table.ent_ptr)
			}
		}
		v_arr_base.Set(reflect.Append(v_arr_base, v_item_base))
	}

	rows.Close()
	return nil
}

/*
func firstToUpper(s string) string {
	if len(s) < 1 {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if i==0 && 'a' <= c && c <= 'z' {
			c -= 'a' - 'A'
		}
		b.WriteByte(c)
	}

	return b.String()
}
func gen_sturct_field(table *DbModel) reflect.StructField {
	v_ent := reflect.ValueOf(table.ent_ptr)
	v_ent = reflect.Indirect(v_ent)
	f_name := firstToUpper(table.table_name)
	return reflect.StructField{Name:f_name, Type:v_ent.Type()}
}
//返回动态结构体
func (lk *DbJoint)GetInterfaceValue() (interface{}, error) {
	f_num := len(lk.md_arr) + 1
	sfs := make([]reflect.StructField, f_num)
	sfs[0] = gen_sturct_field(lk.md_ptr)
	for i, table := range lk.md_arr {
		sfs[i+1] = gen_sturct_field(table)
	}

	v_ent_tmp := reflect.StructOf(sfs)
	v_ent := reflect.New(v_ent_tmp).Elem()

	sql_str := lk.gen_select() + " limit 1"
	rows, err := lk.db_ptr.ExecQuery(lk.getContext(), sql_str, lk.db_where.Values...)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		debug_log.Debug("SQL_DATA:rows_count=0")
		rows.Close()
		return nil, nil
	}

	scan_vals:= lk.get_scan_valus()
	err = rows.Scan(scan_vals...)
	if err != nil {
		rows.Close()
		return nil, err
	}

	v_ent.Field(0).Set(reflect.ValueOf(lk.md_ptr.ent_ptr).Elem())
	for i, table := range lk.md_arr {
		v_ent.Field(i+1).Set(reflect.ValueOf(table.ent_ptr).Elem())
	}

	rows.Close()
	return v_ent.Interface(), nil
}
*/
