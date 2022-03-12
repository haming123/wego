package worm

import (
	"database/sql"
	"errors"
)

type DbEngine struct {
	//缺省DbSession
	def_session *DbSession
	//数据库驱动名称
	db_driver string
	//数据库方言
	db_dialect Dialect
	//数据库
	db_raw *sql.DB
	//数据库从库
	slaves []*DbSlave
	//是否开启stmt缓存
	use_prepare bool
	//数据库stmt缓存
	stmt_cache *StmtCache
	//是否显示SQL日志
	show_sql_log 	bool
	//sql日志打印回调函数
	sql_print_cb 	SqlPrintCB
	//debug日志打印回调函数
	log_print_cb	LogPrintCB
	//sql日志中最大的字段长度
	max_log_field_len int
	//sql日志中最大的select长度
	select_log_len int
	//是否修改insert日志的形式
	show_pretty_log bool
}

//dialect：数据库驱动
//db数据库链接
func NewEngine(dialect Dialect, db *sql.DB) (*DbEngine, error) {
	engine := new (DbEngine)
	engine.db_raw = db
	engine.db_driver = dialect.GetName()
	engine.db_dialect = dialect

	engine.use_prepare = false
	engine.stmt_cache = NewStmtCache(0)

	engine.show_sql_log = true
	engine.sql_print_cb = print_sql_log

	engine.log_print_cb = print_debug_log
	engine.max_log_field_len = 30
	engine.show_pretty_log = true

	engine.def_session = engine.NewSession()

	return engine, nil
}

func NewMysql(db *sql.DB) (*DbEngine, error) {
	return NewEngine(&dialectMysql{}, db)
}

func NewPostgres(db *sql.DB) (*DbEngine, error) {
	return NewEngine(&postgresDialect{}, db)
}

func NewSqlite3(db *sql.DB) (*DbEngine, error) {
	return NewEngine(&dialectSqlite{}, db)
}

func NewSqlServer(db *sql.DB) (*DbEngine, error) {
	return NewEngine(&dialectMssql{}, db)
}

//获取数据库连接池
func (engine *DbEngine)DB() *sql.DB {
	return engine.db_raw
}

//添加一个从库
func (engine *DbEngine)AddSlave(db *sql.DB, db_name string, weight int) error {
	if weight < 1 {
		return errors.New("weight must > 0")
	}

	var slave DbSlave
	slave.db_raw = db
	slave.db_name = db_name
	slave.db_weight = weight
	engine.slaves = append(engine.slaves, &slave)

	return nil
}

//获取数据库方言对象
func (engine *DbEngine)GetDialect() Dialect {
	return engine.db_dialect
}

//是否启用预处理
func (dbs *DbEngine)UsePrepare(flag bool) {
	dbs.use_prepare = flag
}

//设置可缓存的最大stmt数量
func (engine *DbEngine)SetMaxStmtCacheNum(max_len int) {
	if max_len < 100 { max_len = 100 }
	engine.stmt_cache.max_size = max_len
}

//是否显示Sql日志
func (engine *DbEngine)ShowSqlLog(flag bool)  {
	engine.show_sql_log = flag
}

//设置sql log回调函数
func (engine *DbEngine)SetSqlLogCB(cb SqlPrintCB) {
	engine.sql_print_cb = cb
}

func (engine *DbEngine)NewSession() *DbSession {
	return NewDbSession(engine)
}

func (engine *DbEngine)Model(ent_ptr interface{}) *DbModel {
	dbs := engine.NewSession()
	return dbs.Model(ent_ptr)
}

func (engine *DbEngine)Joint(ent_ptr interface{}, alias string, fields ...string) *DbJoint {
	dbs := engine.NewSession()
	return dbs.Joint(ent_ptr, alias, fields...)
}

func (engine *DbEngine)SQL(sql_str string, args ...interface{}) *DbSQL {
	dbs := engine.NewSession()
	return dbs.SQL(sql_str, args...)
}

func (engine *DbEngine)Table(table_name string) *DbTable {
	dbs := engine.NewSession()
	return dbs.Table(table_name)
}
