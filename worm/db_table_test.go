package worm

import (
	"fmt"
	"testing"
)

func TestSqlBuildExpr (t *testing.T) {
	InitEngine4Test()

	tb := Table("user")
	tb.Value("name", "test1")
	tb.Value("age", Expr("id+?", 2))
	tb.Value("created", nil)
	id, err := tb.Insert()
	if err != nil{
		t.Error(err)
		return
	}
	t.Logf("insert id=%d", id)

	tb = Table("user")
	tb.Value("age", Expr("age+?", 1)).Value("created", nil)
	ret, err := tb.Where("id=?", 1).Update()
	if err != nil{
		t.Error(err)
		return
	}
	t.Logf("update num=%d", ret)
}

func TestSqlBuildIUD (t *testing.T) {
	InitEngine4Test()

	id, err := Table("user").Value("name", "test1").Value("age", 11).Insert()
	if err != nil{
		t.Error(err)
		return
	}
	t.Logf("insert id=%d", id)

	ret, err := Table("user").Value("age", 20).Value("name", "zhangsan").Where("id=?", id).Update()
	if err != nil{
		t.Error(err)
		return
	}
	t.Logf("update num=%d", ret)

	ret, err = Table("user").Where("id=?", id).Delete()
	if err != nil{
		t.Error(err)
		return
	}
	t.Logf("delete num=%d", ret)
}

func TestSqlBuildSelectSql (t *testing.T) {
	InitEngine4Test()

	sql := Table("user").Select("*").Where("id>?", 0).OrderBy("name desc").Limit(5).Offset(2)
	t.Log(sql.gen_select())

	sql = Table("user").Select("*").Where("id>?", 0).OrderBy("name desc").Having("age>?", 20)
	t.Log(sql.gen_select())
}

func TestSQLBuilderRows (t *testing.T) {
	InitEngine4Test()

	rows, err := Table("user").Select("*").Where("id>?", 0).OrderBy("name desc").Limit(5).Offset(2).Rows()
	if err != nil{
		t.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next(){
		var user User
		err = ScanModel(rows, &user)
		if err != nil{
			t.Error(err)
		}
		t.Log(user)
	}
}

func TestSQLBuilderExist (t *testing.T) {
	InitEngine4Test()

	has, err := Table("user").Select("*").Where("id=?", 199).Exist()
	if err != nil{
		t.Error(err)
		return
	}
	t.Logf("has=%v\n", has)
}

func TestSQLBuilderJoinRows (t *testing.T) {
	InitEngine4Test()

	rows, err := Table("user").Alias("u").Select("*").Where("u.id>?", 0).OrderBy("u.name desc").Join("book", "b", "b.author=u.id").Rows()
	if err != nil{
		t.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next(){
		data, err := ScanStringRow(rows)
		if err != nil{
			t.Error(err)
		}
		t.Log(data)
	}
}

func TestSQLBuilderGetValue (t *testing.T) {
	InitEngine4Test()

	name := ""
	age := 0
	_, err := Table("user").Select("name,age").Where("id=?", 1).Get(&name, &age)
	if err != nil{
		t.Error(err)
		return
	}
	t.Log(name)
	t.Log(age)
}

func TestSQLBuilderGetString (t *testing.T) {
	InitEngine4Test()

	data, err := Table("user").Select("name").Where("id=?", 1).GetString()
	if err != nil{
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestSQLBuilderGetModel (t *testing.T) {
	InitEngine4Test()

	var ent User
	_, err := Table("user").Select("*").Where("id=?", 4).GetModel(&ent)
	if err != nil{
		t.Error(err)
		return
	}
	t.Log(ent)
}

func TestSQLBuilderGetRow (t *testing.T) {
	InitEngine4Test()

	ent, err := Table("user").Select("*").Where("id=?", 1).GetRow()
	if err != nil{
		t.Error(err)
		return
	}
	t.Log(ent)
}

func TestSQLBuilderFindString (t *testing.T) {
	InitEngine4Test()

	data, err := Table("user").Select("name").Where("id>?", 0).And("name is not null").FindString()
	if err != nil{
		t.Error(err)
		return
	}
	t.Log(data)
}

func TestSQLBuilderFindModel (t *testing.T) {
	InitEngine4Test()

	var users []User
	err := Table("user").Select("*").Where("id>?", 0).FindModel(&users)
	if err != nil{
		t.Error(err)
		return
	}
	for _, user := range users {
		t.Log(user)
	}
}

func TestSQLBuilderFindRow (t *testing.T) {
	InitEngine4Test()

	ret, err := Table("user").Select("*").Where("id>?", 0).FindRow()
	if err != nil{
		t.Error(err)
		return
	}
	rr := ret.GetRowCount()
	t.Logf("row num = %d", rr)
	for i :=0; i < rr; i++ {
		fmt.Println(ret.GetRowData(i))
	}
}

func TestSQLBuilderCount (t *testing.T) {
	InitEngine4Test()

	cc, err := Table("user").Select("name").Where("id>?", 0).GroupBy("name,id").DistinctCount("name")
	if err != nil{
		t.Error(err)
		return
	}
	t.Logf("count=%d\n", cc)
}
