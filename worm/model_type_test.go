package worm

import (
	log "github.com/haming123/wego/dlog"
	"reflect"
	"testing"
	"time"
)

type User0 struct {
	Id      int64     `db:"id;autoincr"`
	Name    string    `db:"name"`
	Passwd  string    `db:"passwd"`
	Age     int       `db:"age"`
	Created time.Time `db:"created;n_update"`
}

func TestModelUser0(t *testing.T) {
	var mo User0
	mi := genModelInfo(reflect.TypeOf(mo))
	log.ShowIndent(true)
	log.DebugJSON(mi)
}

type User1 struct {
	UserId  int64     `db:"uid;autoid"`
	Name    string    `db:"name"`
	Passwd  string    `db:"passwd"`
	Age     int       `db:"age"`
	Created time.Time `db:"created;n_update"`
}

func TesUser1(t *testing.T) {
	var mo User1
	mi := genModelInfo(reflect.TypeOf(mo))
	log.ShowIndent(true)
	log.DebugJSON(mi)
}

func TestModelUser1(t *testing.T) {
	InitEngine4Test()

	var user = User1{Name: "model1", Age: 13}
	id, err := Model(&User1{}).Table("user1").Insert(&user)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("insert id=%d", id)

	user = User1{Name: "model11", Age: 23}
	num, err := Model(&User1{}).Table("user1").Where("uid=?", id).Update(&user)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("update num=%d", num)

	num, err = Model(&User1{}).Table("user1").Where("uid=?", id).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("delete num=%d", num)
}

type User2 struct {
	UserId  int64     `db:"uid;intid"`
	Name    string    `db:"name"`
	Passwd  string    `db:"passwd"`
	Age     int       `db:"age"`
	Created time.Time `db:"created;n_update"`
}

func TestUser2(t *testing.T) {
	var mo User2
	mi := genModelInfo(reflect.TypeOf(mo))
	log.ShowIndent(true)
	log.DebugJSON(mi)
}

func TestModelUser2(t *testing.T) {
	InitEngine4Test()

	var user = User2{UserId: 22, Name: "model", Age: 13}
	id, err := Model(&User2{}).Table("user2").Insert(&user)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("insert id=%d", id)

	user = User2{UserId: 22, Name: "model2", Age: 23}
	num, err := Model(&User2{}).Table("user2").ID(22).Update(&user)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("update num=%d", num)

	num, err = Model(&User2{}).Table("user2").ID(22).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("delete num=%d", num)
}
