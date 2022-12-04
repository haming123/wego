package worm

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

/*
drop table user
CREATE TABLE user (
  id bigint(20) NOT NULL AUTO_INCREMENT,
  name varchar(255) DEFAULT NULL,
  age int(11) DEFAULT NULL,
  created datetime DEFAULT NULL,
  PRIMARY KEY (id)
)
*/

type User struct {
	DB_id   int64 `db:";autoincr"`
	DB_name string
	DB_Salt string `db:"-"`
	Passwd  string
	Age     int       `db:"age"`
	Created time.Time `db:"created;insert_only"`
}

func (ent *User) TableName() string {
	return "user"
}

var pool_user = ModelPool{}

func NewModel_user(dbs *DbSession, auto_put ...bool) *DbModel {
	md := pool_user.Get()
	if md == nil {
		return dbs.Model(&User{}).WithModelPool(&pool_user, auto_put...)
	} else {
		return md.SetDbSession(dbs)
	}
}

func (ent *User) BeforeInsert(ctx context.Context) {
	//fmt.Println("User.BeforeInsert")
}

func (ent *User) AfterInsert(ctx context.Context) {
	//fmt.Println("User.AfterInsert")
}

func (ent *User) BeforeUpdate(ctx context.Context) {
	//fmt.Println("User.BeforeUpdate")
}

func (ent *User) AfterUpdate(ctx context.Context) {
	//fmt.Println("User.AfterUpdate")
}

func (ent *User) BeforeDelete(ctx context.Context) {
	//fmt.Println("User.BeforeDelete")
}

func (ent *User) AfterDelete(ctx context.Context) {
	//fmt.Println("User.AfterDelete")
}

func (ent *User) BeforeQuery(ctx context.Context) {
	//fmt.Println("User.BeforeQuery")
	//if ctx != nil {
	//	fmt.Println(ctx.Value("test_key"))
	//}
}

func (ent *User) AfterQuery(ctx context.Context) {
	//fmt.Println("User.AfterQuery")
}

type DB_Book struct {
	DB_id     int64
	DB_author int64
	DB_name   string
	DB_remark string
}

func (ent *DB_Book) TableName() string {
	return "book"
}

func TestModelIUD(t *testing.T) {
	InitEngine4Test()

	var user = User{DB_name: "model", Age: 13}
	id, err := Model(&user).Insert()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("insert id=%d", id)

	user = User{Age: 31, DB_name: "model2"}
	num, err := Model(&user).Select("name", "age").ID(id).Update()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("update num=%d", num)

	num, err = Model(&User{}).Where("id=?", id).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("delete num=%d", num)
}

func TestModelMoIUD(t *testing.T) {
	InitEngine4Test()

	var user = User{DB_name: "model", Age: 13}
	id, err := Model(&User{}).Insert(&user)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("insert id=%d", id)

	user = User{Age: 31, DB_name: "model2"}
	num, err := Model(&User{}).ID(id).Update(&user)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("update num=%d", num)

	num, err = Model(&User{}).Where("id=?", id).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("delete num=%d", num)
}

func TestModelGet(t *testing.T) {
	InitEngine4Test()

	var ent User
	_, err := Model(&ent).Select("id", "name", "age").ID(1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ent)
}

func TestModelGetMo(t *testing.T) {
	InitEngine4Test()

	var ent User
	_, err := Model(&User{}).Select("id", "name", "age").ID(1).Get(&ent)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ent)
}

func TestModelExist(t *testing.T) {
	InitEngine4Test()

	has, err := Model(&User{}).Where("id=?", 199).Exist()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("has=%v\n", has)
}

func TestModelCount(t *testing.T) {
	InitEngine4Test()

	num, err := Model(&User{}).Where("id>?", 0).Count("name")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("count=%v\n", num)
}

func TestModelRows(t *testing.T) {
	InitEngine4Test()

	rows, err := Model(&User{}).Where("id>?", 0).Rows()
	if err != nil {
		t.Error(err)
		return
	}

	for rows.Next() {
		var user User
		err = rows.ScanModel(&user)
		if err != nil {
			t.Error(err)
		}
		t.Log(user)
	}
	rows.Close()
}

func TestModelFindMo(t *testing.T) {
	InitEngine4Test()

	var users []User
	err := Model(&User{}).Select("id", "name").AndIn("id", 5, 6).Find(&users)
	if err != nil {
		t.Error(err)
		return
	}
	for _, item := range users {
		t.Log(item)
	}
}
