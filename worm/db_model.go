package worm

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type DbModel struct {
	SqlContex
	db_ptr    *DbSession
	ent_ptr   interface{}
	ent_type  reflect.Type
	ent_value reflect.Value

	model_info  *ModelInfo
	table_name  string
	table_alias string
	field_id    string
	flds_info   []FieldInfo
	flds_addr   []FieldValue
	name_map_db map[string]int
	name_map_go map[string]int

	db_where  DbWhere
	group_by  string
	order_by  string
	db_limit  int64
	db_offset int64

	join_type int
	join_on   string
	Err       error

	md_pool  *ModelPool
	auto_put bool

	//自动人工选择标志
	flag_edit bool
	//字段自动选择标志
	flag_auto bool
	//通过Vo选择的字段的缓存数据
	VoFields *PublicFields
}

func NewModel(dbs *DbSession, ent_ptr interface{}, flag bool) *DbModel {
	if ent_ptr == nil {
		panic("ent_ptr must be *Struct")
	}
	v_ent := reflect.ValueOf(ent_ptr)
	if v_ent.Kind() != reflect.Ptr {
		panic("ent_ptr must be *Struct")
	}
	v_ent = reflect.Indirect(v_ent)
	if v_ent.Type().Kind() != reflect.Struct {
		panic("ent_ptr must be *Struct")
	}

	md := &DbModel{}
	md.db_ptr = dbs
	md.ent_ptr = ent_ptr
	md.ent_type = v_ent.Type()
	md.ent_value = v_ent

	minfo := getModelInfo(md.ent_type)
	minfo.TableName = getTableName(md.ent_value, md.ent_type)
	md.flds_addr = getEntFieldAddrs(minfo.Fields, v_ent, flag)
	md.model_info = minfo
	md.flds_info = minfo.Fields
	md.table_name = minfo.TableName
	md.field_id = minfo.FieldID
	md.name_map_db = minfo.NameMapDb
	md.name_map_go = minfo.NameMapGo
	return md
}

//重置model状态，保留以下字段的内容：
//ent_ptr、flds_info、flds_addr、table_name、field_id、name_map_db、name_map_go
func (md *DbModel) Reset() {
	md.db_ptr = nil
	md.table_alias = ""
	md.group_by = ""
	md.order_by = ""
	md.db_limit = 0
	md.db_offset = 0
	md.join_type = 0
	md.join_on = ""
	md.flag_edit = false
	md.flag_auto = false
	md.auto_put = false
	md.md_pool = nil
	md.Err = nil
	md.SqlContex.Reset()
	md.db_where.Reset()
	md.SelectALL()
}

func (md *DbModel) SetDbSession(dbs *DbSession) *DbModel {
	md.db_ptr = dbs
	return md
}

func (md *DbModel) WithModelPool(pool *ModelPool, auto_put ...bool) *DbModel {
	md.md_pool = pool
	if len(auto_put) > 0 {
		md.auto_put = auto_put[0]
	}
	return md
}

func (md *DbModel) PutToPool() {
	if md.md_pool != nil {
		md.md_pool.Put(md)
	}
}

func (md *DbModel) split_pool() *ModelPool {
	pool := md.md_pool
	md.md_pool = nil
	return pool
}

func (md *DbModel) put_pool(pool *ModelPool) {
	if pool != nil {
		md.md_pool = pool
		pool.Put(md)
	}
}

func (md *DbModel) GetModelEnt() interface{} {
	return md.ent_ptr
}

func (md *DbModel) GetContext() context.Context {
	return md.ctx
}

func (md *DbModel) Context(ctx context.Context) *DbModel {
	md.ctx = ctx
	return md
}

func (md *DbModel) UsePrepare(val bool) *DbModel {
	md.use_prepare.Valid = true
	md.use_prepare.Bool = val
	return md
}

func (md *DbModel) ShowLog(val bool) *DbModel {
	md.show_log.Valid = true
	md.show_log.Bool = val
	return md
}

func (md *DbModel) UseMaster(val bool) *DbModel {
	md.use_master.Valid = true
	md.use_master.Bool = val
	return md
}

func (md *DbModel) TableName(val string) *DbModel {
	md.table_name = val
	return md
}

func (md *DbModel) TableAlias(alias string) *DbJoint {
	md.table_alias = alias

	lk := &DbJoint{}
	lk.db_ptr = md.db_ptr
	lk.tables = append(lk.tables, md)
	//TG.SetWhere(&md.db_where)
	return lk
}

func (md *DbModel) get_field_index_byindex(no int) int {
	if no < 0 {
		return -1
	} else if no >= len(md.flds_info) {
		return -1
	} else {
		return no
	}
}

func (md *DbModel) get_field_index_dbname1(dbname string) int {
	index := -1
	num := len(md.flds_info)
	for i := 0; i < num; i++ {
		if md.flds_info[i].DbName == dbname {
			index = i
			break
		}
	}
	return index
}

func (md *DbModel) get_field_index_dbname2(dbname string) int {
	index, ok := md.name_map_db[dbname]
	if ok == false {
		return -1
	}
	return index
}

func (md *DbModel) get_field_index_dbname(dbname string) int {
	if len(md.flds_info) < 20 {
		return md.get_field_index_dbname1(dbname)
	} else {
		return md.get_field_index_dbname2(dbname)
	}
}

func (md *DbModel) get_field_index(dbname string) int {
	return md.get_field_index_dbname(dbname)
}

func (md *DbModel) get_field_index_goname1(goname string) int {
	index := -1
	num := len(md.flds_info)
	for i := 0; i < num; i++ {
		if md.flds_info[i].FieldName == goname {
			index = i
			break
		}
	}
	return index
}

func (md *DbModel) get_field_index_goname2(goname string) int {
	index, ok := md.name_map_go[goname]
	if ok == false {
		return -1
	}
	return index
}

func (md *DbModel) get_field_index_goname(goname string) int {
	if len(md.flds_info) < 20 {
		return md.get_field_index_goname1(goname)
	} else {
		return md.get_field_index_goname2(goname)
	}
}

func (md *DbModel) get_field_index_byaddr(fldg_ptr interface{}) int {
	index := -1
	num := len(md.flds_addr)
	for i := 0; i < num; i++ {
		if md.flds_addr[i].VAddr == fldg_ptr {
			index = i
			break
		}
	}
	return index
}

func (md *DbModel) set_flag_by_index(no int, flag bool) bool {
	index := md.get_field_index_byindex(no)
	if index < 0 {
		return false
	}
	md.flds_addr[index].Flag = flag
	return true
}

func (md *DbModel) set_flag_by_addr(fldg_ptr interface{}, flag bool) bool {
	index := md.get_field_index_byaddr(fldg_ptr)
	if index < 0 {
		return false
	}
	md.flds_addr[index].Flag = flag
	return true
}

func (md *DbModel) SelectALL() *DbModel {
	md.flag_edit = true
	num := len(md.flds_addr)
	for i := 0; i < num; i++ {
		md.flds_addr[i].Flag = true
	}
	return md
}

func (md *DbModel) OmitALL() *DbModel {
	md.flag_edit = true
	num := len(md.flds_addr)
	for i := 0; i < num; i++ {
		md.flds_addr[i].Flag = false
	}
	return md
}

//追加选中一批字段
func (md *DbModel) AddField(fields ...string) *DbModel {
	md.flag_edit = true
	for _, field := range fields {
		//field = strings.Trim(field, " ")
		ind := md.get_field_index(field)
		if ind >= 0 {
			md.flds_addr[ind].Flag = true
		} else {
			md.Err = errors.New("field not find")
		}
	}
	return md
}

//追加选中一批字段
func (md *DbModel) AddFieldX(fields ...interface{}) *DbModel {
	md.flag_edit = true
	for _, fld_ptr := range fields {
		if fld_ptr == nil {
			md.Err = errors.New("field addr is nil")
			return md
		}
		ret := md.set_flag_by_addr(fld_ptr, true)
		if ret == false {
			md.Err = errors.New("field not find")
			return md
		}
	}
	return md
}

//选中一批字段
func (md *DbModel) Select(fields ...string) *DbModel {
	md.flag_edit = true
	//每次都要清空当前选择集
	md.OmitALL()
	return md.AddField(fields...)
}

//选中一批字段
func (md *DbModel) SelectX(fields ...interface{}) *DbModel {
	md.flag_edit = true
	//每次都要清空当前选择集
	md.OmitALL()
	return md.AddFieldX(fields...)
}

//排除若干字段，其余全部选中
func (md *DbModel) Omit(fields ...string) *DbModel {
	md.flag_edit = true
	md.SelectALL()
	for _, field := range fields {
		//field = strings.Trim(field, " ")
		ind := md.get_field_index(field)
		if ind >= 0 {
			md.flds_addr[ind].Flag = false
		} else {
			md.Err = errors.New("field not find")
		}
	}
	return md
}

//排除若干字段，其余全部选中
func (md *DbModel) OmitX(fields ...interface{}) *DbModel {
	md.flag_edit = true
	md.SelectALL()
	for _, fld_ptr := range fields {
		if fld_ptr == nil {
			md.Err = errors.New("field is nil")
			return md
		}
		if reflect.TypeOf(fld_ptr).Kind() != reflect.Ptr {
			md.Err = errors.New("field must be reflect.Ptr")
			return md
		}
		ret := md.set_flag_by_addr(fld_ptr, false)
		if ret == false {
			md.Err = errors.New("field not find")
			return md
		}
	}
	return md
}

func (md *DbModel) AndX(field_ptr interface{}, oper string, vals ...interface{}) *DbModel {
	index := md.get_field_index_byaddr(field_ptr)
	if index < 0 {
		md.Err = errors.New("field_ptr not find")
		return md
	}

	col_name := md.flds_addr[index].FName
	oper = strings.Trim(oper, " ")
	if strings.ToLower(oper) == "in" {
		md.db_where.AndIn(col_name, vals...)
	} else if strings.ToLower(oper) == "not in" {
		md.db_where.AndNotIn(col_name, vals...)
	} else {
		sql := fmt.Sprintf("%s %s ?", col_name, oper)
		md.db_where.And(sql, vals[0])
	}

	return md
}

func (md *DbModel) OrX(field_ptr interface{}, oper string, vals ...interface{}) *DbModel {
	index := md.get_field_index_byaddr(field_ptr)
	if index < 0 {
		md.Err = errors.New("field_ptr not find")
		return md
	}

	col_name := md.flds_addr[index].FName
	oper = strings.Trim(oper, " ")
	if strings.ToLower(oper) == "in" {
		md.db_where.OrIn(col_name, vals...)
	} else if strings.ToLower(oper) == "not in" {
		md.db_where.OrNotIn(col_name, vals...)
	} else {
		sql := fmt.Sprintf("%s %s ?", col_name, oper)
		md.db_where.Or(sql, vals[0])
	}

	return md
}

func (md *DbModel) SetWhere(sqlw *DbWhere) *DbModel {
	md.db_where.Init(sqlw.Tpl_sql, sqlw.Values...)
	return md
}
func (md *DbModel) Where(sql string, vals ...interface{}) *DbModel {
	md.db_where.Init(sql, vals...)
	return md
}
func (md *DbModel) WhereIf(cond bool, sql string, vals ...interface{}) *DbModel {
	if cond {
		md.db_where.Init(sql, vals...)
	}
	return md
}
func (md *DbModel) WhereIn(sql string, vals ...interface{}) *DbModel {
	md.db_where.Reset()
	md.db_where.AndIn(sql, vals...)
	return md
}
func (md *DbModel) WhereNotIn(sql string, vals ...interface{}) *DbModel {
	md.db_where.Reset()
	md.db_where.AndNotIn(sql, vals...)
	return md
}
func (md *DbModel) ID(val int64) *DbModel {
	md.db_where.Init("id=?", val)
	return md
}
func (md *DbModel) And(sql string, vals ...interface{}) *DbModel {
	md.db_where.And(sql, vals...)
	return md
}
func (md *DbModel) Or(sql string, vals ...interface{}) *DbModel {
	md.db_where.Or(sql, vals...)
	return md
}
func (md *DbModel) AndIf(cond bool, sql string, vals ...interface{}) *DbModel {
	md.db_where.AndIf(cond, sql, vals...)
	return md
}
func (md *DbModel) OrIf(cond bool, sql string, vals ...interface{}) *DbModel {
	md.db_where.OrIf(cond, sql, vals...)
	return md
}
func (md *DbModel) AndExp(sqlw_sub *DbWhere) *DbModel {
	md.db_where.AndExp(sqlw_sub)
	return md
}
func (md *DbModel) OrExp(sqlw_sub *DbWhere) *DbModel {
	md.db_where.OrExp(sqlw_sub)
	return md
}
func (md *DbModel) AndIn(sql string, vals ...interface{}) *DbModel {
	md.db_where.AndIn(sql, vals...)
	return md
}
func (md *DbModel) AndNotIn(sql string, vals ...interface{}) *DbModel {
	md.db_where.AndNotIn(sql, vals...)
	return md
}
func (md *DbModel) OrIn(sql string, vals ...interface{}) *DbModel {
	md.db_where.AndIn(sql, vals...)
	return md
}
func (md *DbModel) OrNotIn(sql string, vals ...interface{}) *DbModel {
	md.db_where.OrNotIn(sql, vals...)
	return md
}
func (md *DbModel) OrderBy(val string) *DbModel {
	md.order_by = val
	return md
}
func (md *DbModel) Top(rows int64) *DbModel {
	md.db_limit = rows
	return md
}
func (md *DbModel) Limit(rows int64) *DbModel {
	md.db_limit = rows
	return md
}
func (md *DbModel) Offset(offset int64) *DbModel {
	md.db_offset = offset
	return md
}
func (md *DbModel) Page(rows int64, page_no int64) *DbModel {
	md.db_offset = page_no * rows
	md.db_limit = rows
	return md
}
