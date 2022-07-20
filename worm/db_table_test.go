package worm

import (
	"fmt"
	"testing"
)

func TestSqlBuildSelectSql(t *testing.T) {
	InitEngine4Test()

	tb := Table("user").Select("*").Where("id>?", 0).OrderBy("name desc").Limit(5).Offset(2)
	t.Log(tb.db_ptr.engine.db_dialect.GenTableFindSql(tb))

	tb = Table("user").Select("*").Where("id>?", 0).OrderBy("name desc").Having("age>?", 20)
	t.Log(tb.db_ptr.engine.db_dialect.GenTableFindSql(tb))
}

func TestSqlBuildIUD(t *testing.T) {
	InitEngine4Test()

	id, err := Table("user").Value("name", "test1").Value("age", 11).Insert()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("insert id=%d", id)

	ret, err := Table("user").Value("age", 20).Value("name", "zhangsan").Where("id=?", id).Update()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("update num=%d", ret)

	ret, err = Table("user").Where("id=?", id).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("delete num=%d", ret)
}

func TestSqlBuildExpr(t *testing.T) {
	InitEngine4Test()

	tb := Table("user")
	tb.Value("name", "test1")
	tb.Value("age", Expr("id+?", 2))
	tb.Value("created", nil)
	id, err := tb.Insert()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("insert id=%d", id)

	tb = Table("user")
	tb.Value("age", Expr("age+?", 1)).Value("created", nil)
	ret, err := tb.Where("id=?", 1).Update()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("update num=%d", ret)
}

func TestSQLBuilderGetValues(t *testing.T) {
	InitEngine4Test()

	name := ""
	age := 0
	_, err := Table("user").Select("name,age").Where("id=?", 1).GetValues(&name, &age)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(name)
	t.Log(age)
}

func TestSQLBuilderGetString(t *testing.T) {
	InitEngine4Test()

	data, err := Table("user").Select("name").Where("id=?", 1).GetString()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestSQLBuilderGetTime(t *testing.T) {
	InitEngine4Test()

	data, err := Table("user").Select("created").Where("id=?", 12).GetTime()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestSQLBuilderGetModel(t *testing.T) {
	InitEngine4Test()

	var ent User
	_, err := Table("user").Select("*").Where("id=?", 4).GetModel(&ent)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ent)
}

func TestSQLBuilderRows(t *testing.T) {
	InitEngine4Test()

	rows, err := Table("user").Select("*").Where("id>?", 0).OrderBy("name desc").Limit(5).Offset(2).Rows()
	if err != nil {
		t.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err = ScanModel(rows, &user)
		if err != nil {
			t.Error(err)
		}
		t.Log(user)
	}
}

func TestSQLBuilderFindValues(t *testing.T) {
	InitEngine4Test()

	var ids []int64
	var names []string
	num, err := Table("user").Select("id", "name").Where("id>?", 0).FindValues(&ids, &names)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < num; i++ {
		str := fmt.Sprintf("id=%d, name=%s", ids[i], names[i])
		t.Log(str)
	}
}

func TestSQLBuilderFindString(t *testing.T) {
	InitEngine4Test()

	data, err := Table("user").Select("name").Where("id>?", 0).And("name is not null").FindString()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestSQLBuilderFindModel(t *testing.T) {
	InitEngine4Test()

	var users []User
	err := Table("user").Select("*").Where("id>?", 0).FindModel(&users)
	if err != nil {
		t.Error(err)
		return
	}
	for _, user := range users {
		t.Log(user)
	}
}

func TestSQLBuilderJoinRows(t *testing.T) {
	InitEngine4Test()

	rows, err := Table("user").Alias("u").Select("*").Where("u.id>?", 0).OrderBy("u.name desc").Join("book", "b", "b.author=u.id").Rows()
	if err != nil {
		t.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		data, err := ScanStringRow(rows)
		if err != nil {
			t.Error(err)
		}
		t.Log(data)
	}
}

func TestSQLBuilderExist(t *testing.T) {
	InitEngine4Test()

	has, err := Table("user").Select("*").Where("id=?", 199).Exist()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("has=%v\n", has)
}

func TestSQLBuilderCount(t *testing.T) {
	InitEngine4Test()

	cc, err := Table("user").Select("name").Where("id>?", 0).GroupBy("name,id").DistinctCount("name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("count=%d\n", cc)
}
