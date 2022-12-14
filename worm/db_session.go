package worm

import (
	"database/sql"
	"time"
)

type DbSession struct {
	engine *DbEngine
	//事务会话
	tx_raw *sql.Tx
	//是否显示日志
	show_sql_log bool
	//是否开启stmt缓存
	use_prepare bool
}

func NewDbSession(engine *DbEngine) *DbSession {
	session := &DbSession{}
	session.engine = engine
	session.tx_raw = nil
	session.use_prepare = engine.use_prepare
	session.show_sql_log = engine.show_sql_log
	return session
}

//是否显示日志
func (dbs *DbSession) ShowLog(flag bool) {
	dbs.show_sql_log = flag
}

//是否启用预处理
func (dbs *DbSession) UsePrepare(flag bool) {
	dbs.use_prepare = flag
}

func (dbs *DbSession) DB() *sql.DB {
	return dbs.engine.db_raw
}

func (dbs *DbSession) Tx() *sql.Tx {
	return dbs.tx_raw
}

func (dbs *DbSession) need_print_sql_log(ctx *SqlContex) bool {
	if ctx.show_log.Valid {
		return ctx.show_log.Bool
	}
	return dbs.show_sql_log
}

func (dbs *DbSession) need_prepare(ctx *SqlContex) bool {
	if ctx.use_prepare.Valid {
		return ctx.use_prepare.Bool
	}
	return dbs.use_prepare
}

//开启数据库事务
func (dbs *DbSession) TxBegin() error {
	if dbs.tx_raw != nil {
		panic("transaction has begun")
	}

	tx, err := dbs.engine.db_raw.Begin()
	if err != nil {
		return err
	}
	dbs.tx_raw = tx

	dbs.engine.log_print_cb(LOG_INFO, "[SQL] TxBegin")
	return nil
}

//执行事务回滚
func (dbs *DbSession) TxRollback() error {
	if dbs.tx_raw != nil {
		err := dbs.tx_raw.Rollback()
		if err != nil {
			return err
		}
		dbs.tx_raw = nil
	}

	dbs.engine.log_print_cb(LOG_INFO, "[SQL] TxRollback")
	return nil
}

//执行事务提交
func (dbs *DbSession) TxCommit() error {
	if dbs.tx_raw != nil {
		err := dbs.tx_raw.Commit()
		if err != sql.ErrTxDone && err != nil {
			return err
		}
		dbs.tx_raw = nil
	}

	dbs.engine.log_print_cb(LOG_INFO, "[SQL] TxCommit")
	return nil
}

//提交一个sql预处理
//如果是事务会话，则使用tx_raw.Prepare
//否则使用db_raw.Prepare提交预处理
func (dbs *DbSession) Prepare(sql_tpl string) (*sql.Stmt, error) {
	sql_str := dbs.engine.db_dialect.ParsePlaceholder(sql_tpl)
	if dbs.tx_raw != nil {
		return dbs.tx_raw.Prepare(sql_str)
	} else {
		return dbs.engine.db_raw.Prepare(sql_str)
	}
}

//执行一个sql命令
func (dbs *DbSession) ExecSQL(ctx *SqlContex, sql_tpl string, args ...interface{}) (res sql.Result, err error) {
	log_info := &LogContex{}
	log_info.Session = dbs
	log_info.Start = time.Now()
	log_info.SqlType = "exec"
	log_info.SQL = sql_tpl
	log_info.Args = args
	log_info.Ctx = ctx.ctx

	//将sql语句中的占位符替换为对应数据库的格式
	sql_str := dbs.engine.db_dialect.ParsePlaceholder(sql_tpl)

	//若开启了事务，直接执行sql命令
	//若开启了预处理, 并且有参数，则先调用prepare，然后通过stmt执行命令
	//没有开启预处理，或者参数为0，直接执行sql
	use_stmt_cache := 0
	use_tx := false
	if dbs.tx_raw != nil {
		res, err = dbs.tx_raw.Exec(sql_str, args...)
		use_tx = true
	} else if dbs.need_prepare(ctx) && len(args) > 0 {
		var stmt *sql.Stmt
		var use_stmt int
		stmt, use_stmt, err = dbs.engine.stmt_cache.Prepare(dbs, dbs.engine.db_raw, sql_str)
		if err == nil {
			res, err = stmt.Exec(args...)
		}
		use_stmt_cache = use_stmt
	} else {
		res, err = dbs.engine.db_raw.Exec(sql_str, args...)
	}

	if dbs.need_print_sql_log(ctx) {
		log_info.ExeTime = time.Now().Sub(log_info.Start)
		log_info.Result = res
		log_info.Err = err
		log_info.IsTx = use_tx
		log_info.UseStmt = use_stmt_cache
		dbs.engine.sql_print_cb(log_info)
	}
	return res, err
}

//执行数据库查询
func (dbs *DbSession) ExecQuery(ctx *SqlContex, sql_tpl string, args ...interface{}) (res *sql.Rows, err error) {
	log_info := &LogContex{}
	log_info.Start = time.Now()
	log_info.Session = dbs
	log_info.SqlType = "query"
	log_info.SQL = sql_tpl
	log_info.Args = args
	log_info.Ctx = ctx.ctx

	//将sql语句中的占位符替换为对应数据库的格式
	sql_str := dbs.engine.db_dialect.ParsePlaceholder(sql_tpl)

	//若开启事务，则不使用从库，也不使用预处理
	use_stmt_cache := 0
	use_tx := false
	db_name := ""
	if dbs.tx_raw != nil {
		res, err = dbs.tx_raw.Query(sql_str, args...)
		use_tx = true
	} else {
		//若有存库，选择一个从库进行查询
		var db_raw = dbs.engine.db_raw
		if dbs.engine.slaves != nil && len(dbs.engine.slaves) > 0 {
			slave := getSlaveByWeight(dbs.engine.slaves)
			db_raw = slave.db_raw
			db_name = slave.db_name
		}

		//若明确指定使用master
		if ctx.use_master.Valid && ctx.use_master.Bool {
			db_raw = dbs.engine.db_raw
			db_name = ""
		}

		//若开启了预处理, 并且有参数，则先调用prepare，然后通过stmt执行查询
		//若没有开启预处理，或者参数为0，则直接执行查询
		if dbs.need_prepare(ctx) && len(args) > 0 {
			var stmt *sql.Stmt
			var use_stmt int
			stmt, use_stmt, err = dbs.engine.stmt_cache.Prepare(dbs, db_raw, sql_str)
			if err == nil {
				res, err = stmt.Query(args...)
			}
			use_stmt_cache = use_stmt
		} else {
			res, err = db_raw.Query(sql_str, args...)
		}
	}

	if dbs.need_print_sql_log(ctx) {
		log_info.ExeTime = time.Now().Sub(log_info.Start)
		log_info.Err = err
		log_info.IsTx = use_tx
		log_info.UseStmt = use_stmt_cache
		log_info.DbName = db_name
		dbs.engine.sql_print_cb(log_info)
	}
	return res, err
}

func (dbs *DbSession) NewModel(ent_ptr interface{}, flag bool) *DbModel {
	return NewModel(dbs, ent_ptr, flag)
}

func (dbs *DbSession) Model(ent_ptr interface{}) *DbModel {
	return NewModel(dbs, ent_ptr, true)
}

func (dbs *DbSession) Joint(ent_ptr interface{}, alias string, fields ...string) *DbJoint {
	return NewJoint(dbs, ent_ptr, alias, fields...)
}

func (dbs *DbSession) SQL(sql_str string, args ...interface{}) *DbSQL {
	return NewDbSQL(dbs, sql_str, args...)
}

func (dbs *DbSession) Table(table_name string) *DbTable {
	return NewDbTable(dbs, table_name)
}
