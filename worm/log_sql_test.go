package worm

import (
	"context"
	"testing"
)

func TestSqlLogInsert (t *testing.T) {
	sql := "insert into user (name,age) values (?,?)"
	ret := parsel_debug_sql_insert(sql)
	t.Log(ret)
	sql = "insert into user (name,age) values (?)"
	ret = parsel_debug_sql_insert(sql)
	t.Log(ret)
}

func TestSqlLogHolder (t *testing.T) {
	sql := "insert into user (name,age) values (?,?)"
	ret := parsel_debug_sql_holder(sql, []interface{}{"lisi", 20}, 0)
	t.Log(ret)
	ret = parsel_debug_sql_holder(sql, []interface{}{"12345678", 20}, 5)
	t.Log(ret)
}

var g_log *SimpleLogger = NewSimpleLogger()
func my_print_sql_log(ctx *LogContex) {
	g_log.Info(ctx.GetPrintSQL())
	g_log.Info(ctx.GetPrintResult())
}

func TestSqlLogCB(t *testing.T) {
	InitEngine4Test()
	SetSqlLogCB(my_print_sql_log)
	SetDebugLogLevel(LOG_DEBUG)

	var ent User
	_, err := Model(&ent).Select("id", "name", "age").Where("id=? or age>?", 9, 0).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ent)
}

func my_print_sql_log2(ctx *LogContex) {
	g_log.Info(ctx.GetPrintSQL())
	g_log.Info(ctx.GetPrintResult())
	g_log.Debug(ctx.Ctx.Value("test_key"))
}

func TestSqlLogCBContext(t *testing.T) {
	InitEngine4Test()
	SetSqlLogCB(my_print_sql_log2)
	SetDebugLogLevel(LOG_DEBUG)

	var ent User
	ctx := context.WithValue(context.Background(), "test_key", "test_val")
	_, err := Model(&ent).Context(ctx).Select("id", "name", "age").Where("id=? or age>?", 9, 0).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ent)
}
