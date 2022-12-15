package worm

import (
	"reflect"
	"strings"
	"sync"
)

/*
type DB_User struct {
	DB_id      	int64		`db:";autoincr"`
	DB_name    	string
	Age     	int			`db:"age"`
	Creatat    time.Time	`db:"creatat;insert_only"`
}
*/

const (
	STR_AUTOINCR     string = "autoincr"
	STR_NOT_AUTOINCR string = "n_autoincr"
	STR_NOT_INSERT   string = "n_insert"
	STR_NOT_UPDATE   string = "n_update"
	STR_NOT_SELECT   string = "n_select"
)

type TableName interface {
	TableName() string
}

type FieldInfo struct {
	FieldIndex int
	FieldName  string
	FieldType  reflect.Type
	DbName     string
	AutoIncr   bool
	NotInsert  bool
	NotUpdate  bool
	NotSelect  bool
}

type ModelInfo struct {
	Fields    []FieldInfo
	TableName string
	FieldID   string
	NameMapDb map[string]int
	NameMapGo map[string]int
}

func (mi *ModelInfo) get_field_index_dbname1(dbname string) int {
	index := -1
	num := len(mi.Fields)
	for i := 0; i < num; i++ {
		if mi.Fields[i].DbName == dbname {
			index = i
			break
		}
	}
	return index
}

func (mi *ModelInfo) get_field_index_dbname2(dbname string) int {
	index, ok := mi.NameMapDb[dbname]
	if ok == false {
		return -1
	}
	return index
}

func (mi *ModelInfo) get_field_index_dbname(dbname string) int {
	if len(mi.Fields) < 20 {
		return mi.get_field_index_dbname1(dbname)
	} else {
		return mi.get_field_index_dbname2(dbname)
	}
}

func (mi *ModelInfo) get_field_index_goname1(goname string) int {
	index := -1
	num := len(mi.Fields)
	for i := 0; i < num; i++ {
		if mi.Fields[i].FieldName == goname {
			index = i
			break
		}
	}
	return index
}

func (mi *ModelInfo) get_field_index_goname2(goname string) int {
	index, ok := mi.NameMapGo[goname]
	if ok == false {
		return -1
	}
	return index
}

func (mi *ModelInfo) get_field_index_goname(goname string) int {
	if len(mi.Fields) < 20 {
		return mi.get_field_index_goname1(goname)
	} else {
		return mi.get_field_index_goname2(goname)
	}
}

//struct信息的缓存
var g_model_cache map[reflect.Type]*ModelInfo = make(map[reflect.Type]*ModelInfo)
var g_model_mutex sync.Mutex

//获取model的信息数据
//若存在缓存数据，直接返回缓存数据
//否则，生成model的信息数据，并添加到缓存中
func getModelInfo(t_ent reflect.Type) *ModelInfo {
	g_model_mutex.Lock()
	defer g_model_mutex.Unlock()

	info, ok := g_model_cache[t_ent]
	if ok {
		return info
	}

	info = genModelInfo(t_ent)
	g_model_cache[t_ent] = info
	return info
}

//生成model的信息数据
func genModelInfo(t_ent reflect.Type) *ModelInfo {
	minfo := ModelInfo{}
	minfo.Fields = make([]FieldInfo, t_ent.NumField())
	minfo.NameMapDb = make(map[string]int)
	minfo.NameMapGo = make(map[string]int)
	for i := 0; i < t_ent.NumField(); i++ {
		ff := t_ent.Field(i)

		finfo := FieldInfo{}
		finfo.FieldIndex = -1

		//获取字段的数据库名称
		field_name := ff.Name
		db_name := getFieldName(ff)
		if len(db_name) < 1 {
			minfo.Fields[i] = finfo
			continue
		}

		finfo.FieldIndex = i
		finfo.FieldName = field_name
		finfo.FieldType = ff.Type
		finfo.DbName = db_name

		//获取字段的tag属性
		if strings.ToLower(db_name) == "id" {
			minfo.FieldID = db_name
			finfo.AutoIncr = true
			finfo.NotUpdate = true
		}
		parselFeildTag(&finfo, ff)

		minfo.Fields[i] = finfo
		minfo.NameMapDb[db_name] = i
		minfo.NameMapGo[ff.Name] = i
	}

	return &minfo
}

//获取model结构体对应的数据库表名称
func getTableName(v_ent reflect.Value, t_ent reflect.Type) string {
	var tpTableName = reflect.TypeOf((*TableName)(nil)).Elem()
	if t_ent.Implements(tpTableName) {
		return v_ent.Interface().(TableName).TableName()
	}

	if v_ent.Kind() == reflect.Ptr {
		v_ent = v_ent.Elem()
		if v_ent.Type().Implements(tpTableName) {
			return v_ent.Interface().(TableName).TableName()
		}
	} else if v_ent.CanAddr() {
		v1 := v_ent.Addr()
		if v1.Type().Implements(tpTableName) {
			return v1.Interface().(TableName).TableName()
		}
	}

	t_name := t_ent.Name()
	t_name = strings.ToLower(t_name)
	ind := strings.Index(t_name, "db_")
	if ind >= 0 {
		ind += 3
		t_name = t_name[ind:]
	}
	return t_name
}

//获取model结构体字段对应的数据库字段名称
//优先从tag中获取数据库字段名称
//若不存在tag，则从结构体字段名称中解析数据库字段名称
func getFieldName(ff reflect.StructField) string {
	f_name := ""
	ind := strings.Index(ff.Name, "DB_")
	if ind >= 0 {
		ind += 3
		f_name = ff.Name[ind:]
	}

	tag := ff.Tag.Get("db")
	parts := strings.Split(tag, ";")
	part0 := strings.Trim(parts[0], " ")
	if part0 == "-" {
		f_name = ""
	} else if part0 != "" {
		f_name = part0
	}

	return f_name
}

//获取model结构体字段的tag属性
func parselFeildTag(finfo *FieldInfo, ff reflect.StructField) {
	tag := ff.Tag.Get("db")
	if tag == "" {
		return
	}

	parts := strings.Split(tag, ";")
	for i, item := range parts {
		//first part is field name
		if i == 0 {
			continue
		}

		item = strings.Trim(item, " ")
		if item == STR_AUTOINCR {
			finfo.AutoIncr = true
		} else if item == STR_NOT_AUTOINCR {
			finfo.AutoIncr = false
		} else if item == STR_NOT_INSERT {
			finfo.NotInsert = true
		} else if item == STR_NOT_UPDATE {
			finfo.NotUpdate = true
		} else if item == STR_NOT_SELECT {
			finfo.NotSelect = true
		}
	}
}
