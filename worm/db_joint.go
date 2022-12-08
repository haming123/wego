package worm

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
)

const (
	JOIN_INNER int = 0
	JOIN_LJOIN int = 1
	JOIN_RJOIN int = 2
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

type JointEoFieldCache struct {
	models []PublicFields
}

func newJointEoFieldCache(md_arr []*DbModel) *JointEoFieldCache {
	var cache JointEoFieldCache
	cache.models = make([]PublicFields, len(md_arr))
	for m := 0; m < len(md_arr); m++ {
		cache.models[m].ModelField = -1
		cache.models[m].Fields = make([]FieldIndex, len(md_arr[m].flds_info))
		cache.models[m].Fields = cache.models[m].Fields[:0]
	}
	return &cache
}

//vo、mo字段交集缓存
var g_joint_field_cache map[string]*JointEoFieldCache = make(map[string]*JointEoFieldCache)
var g_joint_field_mutex sync.Mutex

type DbJoint struct {
	SqlContex
	db_ptr *DbSession
	//md_ptr    *DbModel
	md_arr    []*DbModel
	db_where  DbWhere
	order_by  string
	db_limit  int64
	db_offset int64
	Err       error
	ctx       context.Context
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
	//lk.md_ptr = md
	lk.md_arr = append(lk.md_arr, md)
	return lk
}

func (lk *DbJoint) GetContext() context.Context {
	return lk.ctx
}

func (lk *DbJoint) Context(ctx context.Context) *DbJoint {
	lk.ctx = ctx
	return lk
}

func (lk *DbJoint) UsePrepare(val bool) *DbJoint {
	lk.use_prepare.Valid = true
	lk.use_prepare.Bool = val
	return lk
}

func (lk *DbJoint) ShowLog(val bool) *DbJoint {
	lk.show_log.Valid = true
	lk.show_log.Bool = val
	return lk
}

func (lk *DbJoint) UseMaster(val bool) *DbJoint {
	lk.use_master.Valid = true
	lk.use_master.Bool = val
	return lk
}

func (lk *DbJoint) get_table_index(ent_type reflect.Type) int {
	index := -1
	num := len(lk.md_arr)
	for i := 0; i < num; i++ {
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
	for i := 0; i < len(fields); i++ {
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

func (lk *DbJoint) Join(ent_ptr interface{}, alias string, join_on string, fields ...string) *DbJoint {
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

func (lk *DbJoint) LeftJoin(ent_ptr interface{}, alias string, join_on string, fields ...string) *DbJoint {
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

func (lk *DbJoint) RightJoin(ent_ptr interface{}, alias string, join_on string, fields ...string) *DbJoint {
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

func (lk *DbJoint) SetWhere(sqlw *DbWhere) *DbJoint {
	lk.db_where.Init(sqlw.Tpl_sql, sqlw.Values...)
	return lk
}
func (lk *DbJoint) Where(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.Init(sql, vals...)
	return lk
}
func (lk *DbJoint) WhereIf(cond bool, sql string, vals ...interface{}) *DbJoint {
	if cond {
		lk.db_where.Init(sql, vals...)
	}
	return lk
}
func (lk *DbJoint) WhereIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.Reset()
	lk.db_where.AndIn(sql, vals...)
	return lk
}
func (lk *DbJoint) WhereNotIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.Reset()
	lk.db_where.AndNotIn(sql, vals...)
	return lk
}
func (lk *DbJoint) And(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.And(sql, vals...)
	return lk
}
func (lk *DbJoint) AndIf(cond bool, sql string, vals ...interface{}) *DbJoint {
	lk.db_where.AndIf(cond, sql, vals...)
	return lk
}
func (lk *DbJoint) OrIf(cond bool, sql string, vals ...interface{}) *DbJoint {
	lk.db_where.OrIf(cond, sql, vals...)
	return lk
}
func (lk *DbJoint) Or(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.Or(sql, vals...)
	return lk
}
func (lk *DbJoint) AndExp(sqlw_sub *DbWhere) *DbJoint {
	lk.db_where.AndExp(sqlw_sub)
	return lk
}
func (lk *DbJoint) OrExp(sqlw_sub *DbWhere) *DbJoint {
	lk.db_where.OrExp(sqlw_sub)
	return lk
}
func (lk *DbJoint) AndIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.AndIn(sql, vals...)
	return lk
}
func (lk *DbJoint) AndNotIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.AndNotIn(sql, vals...)
	return lk
}
func (lk *DbJoint) OrIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.AndIn(sql, vals...)
	return lk
}
func (lk *DbJoint) OrNotIn(sql string, vals ...interface{}) *DbJoint {
	lk.db_where.OrNotIn(sql, vals...)
	return lk
}
func (lk *DbJoint) OrderBy(val string) *DbJoint {
	lk.order_by = val
	return lk
}
func (lk *DbJoint) Top(rows int64) *DbJoint {
	lk.db_limit = rows
	return lk
}
func (lk *DbJoint) Limit(rows int64) *DbJoint {
	lk.db_limit = rows
	return lk
}
func (lk *DbJoint) Offset(offset int64) *DbJoint {
	lk.db_offset = offset
	return lk
}
func (lk *DbJoint) Page(rows int64, page_no int64) *DbJoint {
	lk.db_offset = page_no * rows
	lk.db_limit = rows
	return lk
}

func (lk *DbJoint) get_scan_valus() []interface{} {
	//vals := lk.md_ptr.get_scan_valus()
	var vals []interface{}
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

func (lk *DbJoint) Scan() (bool, error) {
	if lk.Err != nil {
		return false, lk.Err
	}
	sql_str := lk.db_ptr.engine.db_dialect.GenJointGetSql(lk)
	rows, err := lk.db_ptr.ExecQuery(&lk.SqlContex, sql_str, lk.db_where.Values...)
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		rows.Close()
		return false, nil
	}

	scan_vals := lk.get_scan_valus()
	err = rows.Scan(scan_vals...)
	if err != nil {
		rows.Close()
		return false, err
	}

	//call_after_query(lk.md_ptr)
	for _, table := range lk.md_arr {
		call_after_query(table)
	}

	rows.Close()
	return true, nil
}

/*
//通过model的类型来获取一个model地址
func (lk *DbJoint) get_model_by_ent_type(ent_type reflect.Type) *DbModel {
	if reflect.TypeOf(lk.md_ptr.ent_ptr).Elem() == ent_type {
		return lk.md_ptr
	}
	num := len(lk.md_arr)
	for i := 0; i < num; i++ {
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
func (lk *DbJoint) BindAddr2Struct(v_ent reflect.Value) {
	t_num := v_ent.NumField()
	for t := 0; t < t_num; t++ {
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

func (lk *DbJoint) Get(args ...interface{}) (bool, error) {
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
	if ent_ptr == nil {
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

	sql_str := lk.db_ptr.engine.db_dialect.GenJointGetSql(lk)
	rows, err := lk.db_ptr.ExecQuery(&lk.SqlContex, sql_str, lk.db_where.Values...)
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		rows.Close()
		return false, nil
	}

	scan_vals := lk.get_scan_valus()
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
*/

//查找与Model名称、类型一致的字段，选中该字段，记录该字段的索引位置
func (lk *DbJoint) genPubField4VoMoNest(cache *JointEoFieldCache, t_vo reflect.Type, pos FieldPos, deep int) {
	//超过最大层次，则退出
	if deep >= len(pos) {
		return
	}

	f_num := t_vo.NumField()
	for ff := 0; ff < f_num; ff++ {
		ft_vo := t_vo.Field(ff)

		//首先查找与Model类型一致的字段，并将字段的索引赋值给Model
		//只有第1层(deep=0)字段才判断是否为Model类型
		if deep < 1 && ft_vo.Type.Kind() == reflect.Struct {
			is_model_field := false
			for m := 0; m < len(lk.md_arr); m++ {
				md := lk.md_arr[m]
				if ft_vo.Type == md.ent_type {
					cache.models[m].ModelField = ff
					is_model_field = true
					break
				}
			}
			if is_model_field {
				continue
			}
		}

		//若是匿名字段,则递归调用
		if ft_vo.Anonymous {
			pos[deep] = ff
			lk.genPubField4VoMoNest(cache, ft_vo.Type, pos, deep+1)
			continue
		}

		for m := 0; m < len(lk.md_arr); m++ {
			md := lk.md_arr[m]
			if cache.models[m].ModelField >= 0 {
				continue
			}

			//通过eo字段名称查询model中字段的位置
			mo_index := md.get_field_index_byname2(ft_vo.Name)
			if mo_index < 0 {
				continue
			}

			//只有类型与名称一致才算匹配上
			if md.flds_info[mo_index].FieldType != ft_vo.Type {
				continue
			}

			var item FieldIndex
			item.FieldName = ft_vo.Name
			item.MoIndex = mo_index
			item.VoField = pos
			item.VoField[deep] = ff
			item.VoIndex = item.VoField[0 : deep+1]
			cache.models[m].Fields = append(cache.models[m].Fields, item)
			break
		}
	}
}

//首先从缓存中获取字段交集
//若缓存中不存在，则生成字段交集
func (lk *DbJoint) getPubField4VoMo(t_vo reflect.Type) {
	g_joint_field_mutex.Lock()
	defer g_joint_field_mutex.Unlock()

	//生成缓存的key
	cache_key := t_vo.String()
	for m := 0; m < len(lk.md_arr); m++ {
		md := lk.md_arr[m]
		cache_key += md.ent_type.String()
	}

	//从缓存中获取字段交集
	cache, ok := g_joint_field_cache[cache_key]
	if !ok {
		//没有字段交集的缓存，调用genPubField4VoMoNest生成字段交集
		var pos FieldPos
		cache = newJointEoFieldCache(lk.md_arr)
		lk.genPubField4VoMoNest(cache, t_vo, pos, 0)
		g_joint_field_cache[cache_key] = cache
	}

	//将公共字段添加到选择集中
	for m := 0; m < len(lk.md_arr) && m < len(cache.models); m++ {
		pflds := &cache.models[m]
		md := lk.md_arr[m]
		md.VoFields = pflds

		//若进行了字段的人工选择，则不需要进行字段的自动选择
		if md.flag_edit {
			continue
		}

		//若存在model类型一致的字段，不用额外选择字段（缺省选择全部）
		if pflds.ModelField < 0 {
			for _, item := range pflds.Fields {
				md.auto_add_field_index(item.MoIndex)
			}
		}
	}
}

//把Model中地址的值赋值给vo对象
func (lk *DbJoint) CopyModelData2Eo(v_vo reflect.Value) {
	for m := 0; m < len(lk.md_arr); m++ {
		md := lk.md_arr[m]
		pflds := md.VoFields

		//若vo中存在Model字段，只需要赋值Model对应的字段即可
		if pflds.ModelField >= 0 {
			fv_vo := v_vo.Field(pflds.ModelField)
			if fv_vo.CanSet() == true {
				fv_vo.Set(md.ent_value)
			}
			continue
		}

		for _, item := range pflds.Fields {
			fv_vo := v_vo.FieldByIndex(item.VoIndex)
			fv_mo := md.ent_value.Field(item.MoIndex)
			if fv_vo.CanSet() == false {
				continue
			}
			fv_vo.Set(fv_mo)
		}
	}
}

//把Model中地址的值赋值给vo对象
func (lk *DbJoint) CopyModelData2Vo(vo_ptr VoLoader) {
	for _, md := range lk.md_arr {
		vo_ptr.LoadFromModel(nil, md.ent_ptr)
	}
}

//通过eo(struct)对象来选择需要查询的字段
func (lk *DbJoint) select_field_by_eo2(eo_ptr interface{}) {
	for _, md := range lk.md_arr {
		selectFieldsByEo(md, eo_ptr)
	}
}

//通过eo(struct)对象来选择需要查询的字段
func (lk *DbJoint) select_field_by_eo(t_vo reflect.Type) {
	lk.getPubField4VoMo(t_vo)
}

//通过vo对象来选择需要查询的字段
func (lk *DbJoint) select_field_by_vo(vo_ptr VoLoader) {
	for _, table := range lk.md_arr {
		selectFieldsByVo(table, vo_ptr)
	}
}

func (lk *DbJoint) Get(args ...interface{}) (bool, error) {
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
	if ent_ptr == nil {
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
	//若目标对象不是一个vo，则通过与eo的字段交集来选择字段
	vo_ptr, isvo := ent_ptr.(VoLoader)
	if isvo {
		lk.select_field_by_vo(vo_ptr)
	} else {
		lk.select_field_by_eo(v_ent.Type())
		//lk.select_field_by_eo2(ent_ptr)
	}

	sql_str := lk.db_ptr.engine.db_dialect.GenJointGetSql(lk)
	rows, err := lk.db_ptr.ExecQuery(&lk.SqlContex, sql_str, lk.db_where.Values...)
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		rows.Close()
		return false, nil
	}

	scan_vals := lk.get_scan_valus()
	err = rows.Scan(scan_vals...)
	if err != nil {
		rows.Close()
		return false, err
	}

	//call_after_query(lk.md_ptr)
	for _, table := range lk.md_arr {
		call_after_query(table)
	}

	//若目标对象是一个vo，则调用LoadFromModel，给vo赋值
	if isvo {
		//for _, table := range lk.md_arr {
		//	vo_ptr.LoadFromModel(nil, table.ent_ptr)
		//}
		lk.CopyModelData2Vo(vo_ptr)
	} else {
		//for _, table := range lk.md_arr {
		//	CopyDataFromModel(nil, ent_ptr, table.ent_ptr)
		//}
		lk.CopyModelData2Eo(v_ent)
	}

	rows.Close()
	return true, nil
}

/*
func (lk *DbJoint) Find(arr_ptr interface{}) error {
	if lk.Err != nil {
		return lk.Err
	}

	v_arr := reflect.ValueOf(arr_ptr)
	if v_arr.Kind() != reflect.Ptr {
		return errors.New("arr_ptr must be *Slice")
	}
	t_arr := GetDirectType(v_arr.Type())
	if t_arr.Kind() != reflect.Slice {
		return errors.New("arr_ptr must be *Slice")
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

	sql_str := lk.db_ptr.engine.db_dialect.GenJointFindSql(lk)
	vals := lk.get_scan_valus()
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
*/

func (lk *DbJoint) Find(arr_ptr interface{}) error {
	if lk.Err != nil {
		return lk.Err
	}

	v_arr := reflect.ValueOf(arr_ptr)
	if v_arr.Kind() != reflect.Ptr {
		return errors.New("arr_ptr must be *Slice")
	}
	t_arr := GetDirectType(v_arr.Type())
	if t_arr.Kind() != reflect.Slice {
		return errors.New("arr_ptr must be *Slice")
	}
	//获取数组成员的类型
	t_item := GetDirectType(t_arr.Elem())
	if t_item.Kind() != reflect.Struct {
		return errors.New("array item muse be Struct")
	}

	//若目标对象是一个vo，则通过vo来选择字段
	//若目标对象不是一个vo，则通过与eo的字段交集来选择字段
	v_item := reflect.New(t_item)
	v_item_base := v_item.Elem()
	ent_ptr := v_item.Interface()
	vo_ptr, isvo := v_item.Interface().(VoLoader)
	if isvo {
		lk.select_field_by_vo(vo_ptr)
	} else {
		lk.select_field_by_eo(t_item)
	}

	sql_str := lk.db_ptr.engine.db_dialect.GenJointFindSql(lk)
	vals := lk.get_scan_valus()
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

		//call_after_query(lk.md_ptr)
		for _, table := range lk.md_arr {
			call_after_query(table)
		}

		//若目标对象是一个vo，则调用LoadFromModel，给v_item_base赋值
		//若目标对象是一个eo，则调用CopyDataFromModel，给v_item_base赋值
		if isvo {
			//vo_ptr.LoadFromModel(nil, lk.md_ptr.ent_ptr)
			for _, table := range lk.md_arr {
				vo_ptr.LoadFromModel(nil, table.ent_ptr)
			}
		} else {
			//CopyDataFromModel(nil, ent_ptr, lk.md_ptr.ent_ptr)
			for _, table := range lk.md_arr {
				CopyDataFromModel(nil, ent_ptr, table.ent_ptr)
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
