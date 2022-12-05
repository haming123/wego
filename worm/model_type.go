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
	STR_AUTOINCR   string = "autoincr"
	STR_NOT_INSERT string = "n_insert"
	STR_NOT_UPDATE string = "n_update"
	STR_NOT_SELECT string = "n_select"
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
	NameMap   map[string]int
	TableName string
	FieldID   string
}

//struct信息的缓存
var g_model_cache map[reflect.Type]*ModelInfo = make(map[reflect.Type]*ModelInfo)
var g_model_mutex sync.Mutex

func getModelInfoUseCache(v_ent reflect.Value) *ModelInfo {
	g_model_mutex.Lock()
	defer g_model_mutex.Unlock()

	v_ent = reflect.Indirect(v_ent)
	t_ent := v_ent.Type()

	info, ok := g_model_cache[t_ent]
	if ok {
		return info
	}

	info = getModelInfo(v_ent)
	g_model_cache[t_ent] = info
	return info
}

func getModelInfo(v_ent reflect.Value) *ModelInfo {
	minfo := ModelInfo{}
	minfo.TableName = getTableName(v_ent)

	v_ent = reflect.Indirect(v_ent)
	t_ent := v_ent.Type()

	f_num := t_ent.NumField()
	minfo.Fields = make([]FieldInfo, f_num)
	minfo.NameMap = make(map[string]int)
	for i := 0; i < f_num; i++ {
		finfo := FieldInfo{}
		finfo.FieldIndex = -1

		ff := t_ent.Field(i)
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
		parselFeildTag(&finfo, ff)
		if strings.ToLower(db_name) == "id" {
			minfo.FieldID = db_name
			finfo.AutoIncr = true
		}

		minfo.Fields[i] = finfo
		minfo.NameMap[db_name] = i
	}

	return &minfo
}

/*
//获取model信息（支持匿名字段）
func getModelInfoAnonymous(v_ent reflect.Value) *ModelInfo {
	minfo := ModelInfo{}
	minfo.TableName = getTableName(v_ent)

	v_ent = reflect.Indirect(v_ent)
	t_ent := v_ent.Type()
	getModelInfoNest(&minfo, t_ent, nil)

	return &minfo
}
//获取model信息递归调用
func getModelInfoNest(minfo *ModelInfo, t_ent reflect.Type, pos []int) {
	f_num := t_ent.NumField()
	for i := 0; i < f_num; i++ {
		ff := t_ent.Field(i)
		if ff.Anonymous == true {
			getModelInfoNest(minfo, ff.Type, append(pos, i))
			continue
		}

		field_name := ff.Name
		db_name := getFieldName(ff)
		if len(db_name) < 1 {
			continue
		}

		finfo := FieldInfo{}
		finfo.FieldIndex = i
		if pos != nil {
			finfo.FieldPos = append(pos, i)
		}
		finfo.FieldName = field_name
		finfo.FieldType = ff.Type
		finfo.DbName = db_name
		parselFeildTag(&finfo, ff)
		if strings.ToLower(db_name) == "id" {
			minfo.FieldID = db_name
			finfo.AutoIncr = true
		}

		minfo.Fields = append(minfo.Fields, finfo)
	}
}
*/

func getTableName(v_ent reflect.Value) string {
	var t_ent = v_ent.Type()
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
		} else if item == STR_NOT_INSERT {
			finfo.NotInsert = true
		} else if item == STR_NOT_UPDATE {
			finfo.NotUpdate = true
		} else if item == STR_NOT_SELECT {
			finfo.NotSelect = true
		}
	}
}
