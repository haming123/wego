package worm

import (
	"testing"
	"wego/worm/conf"
)

func TestNewSession (t *testing.T) {
	InitEngine4Test()
	dbs := NewSession()
	dbs.ShowLog(false)

	var users []User
	err := dbs.Model(&User{}).Select("created", "name").AndIn("id", 5,6).Find(&users)
	if err != nil{
		t.Error(err)
		return
	}
	for _, item := range users {
		t.Log(item)
	}
}

func TestSessionTxRollback (t *testing.T) {
	InitEngine4Test()
	dbs := NewSession()
	dbs.TxBegin()

	var user = User{DB_name:"model", Age: 13 }
	id, err := dbs.Model(&user).Insert()
	if err != nil{
		t.Error(err)
		return
	}
	t.Logf("insert id=%d", id)

	user = User{Age: 31, DB_name: "model2"}
	num, err := dbs.Model(&user).Where("iid=?", id).Update()
	if err != nil{
		t.Log(err)
		dbs.TxRollback()
		return
	}
	t.Logf("update num=%d", num)

	num, err = Model(&User{}).Where("id=?", id).Delete()
	if err != nil{
		t.Error(err)
		dbs.TxRollback()
		return
	}
	t.Logf("delete num=%d", num)
	dbs.TxCommit()
}

func TestSessionTxCommit (t *testing.T) {
	InitEngine4Test()
	dbs := NewSession()
	dbs.TxBegin()

	var user = User{DB_name:"model", Age: 13 }
	id, err := dbs.Model(&user).Insert()
	if err != nil{
		t.Error(err)
		dbs.TxRollback()
		return
	}

	_, err = dbs.Model(&user).Where("id=?", id).Get()
	if err != nil {
		t.Error(err)
		dbs.TxRollback()
		return
	}
	t.Log(user)

	user = User{Age: 31, DB_name: "model2"}
	_, err = dbs.Model(&user).Where("id=?", id).Update()
	if err != nil{
		t.Error(err)
		dbs.TxRollback()
		return
	}

	_, err = dbs.Model(&user).Where("id=?", id).Get()
	if err != nil {
		t.Error(err)
		dbs.TxRollback()
		return
	}
	t.Log(user)

	_, err = dbs.Model(&User{}).Where("id=?", id).Delete()
	if err != nil{
		t.Error(err)
		dbs.TxRollback()
		return
	}
	dbs.TxCommit()
}

func TestSessionUsePrepare (t *testing.T) {
	InitEngine4Test()
	dbs := NewSession()
	dbs.UsePrepare(true)

	user := User{Age: 31, DB_name: "model111"}
	_, err := dbs.Model(&user).Where("id=?", 1).Update()
	if err != nil{
		t.Error(err)
		return
	}

	_, err = dbs.Model(&user).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)

	user = User{Age: 31, DB_name: "model222"}
	_, err = dbs.Model(&user).Where("id=?", 1).Update()
	if err != nil{
		t.Error(err)
		return
	}

	_, err = dbs.Model(&user).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}

func TestSessionUseSlave (t *testing.T) {
	InitEngine4Test()
	db_slave, err := conf.OpenSalveDb()
	if err != nil {
		t.Error(err)
		return
	}
	AddSlave(db_slave, "slave1", 1)

	user := User{Age: 31, DB_name: "model111"}
	_, err = Model(&user).Where("id=?", 1).Update()
	if err != nil{
		t.Error(err)
		return
	}

	user = User{}
	_, err = Model(&user).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)

	user = User{Age: 31, DB_name: "model222"}
	_, err = Model(&user).Where("id=?", 1).Update()
	if err != nil{
		t.Error(err)
		return
	}

	user = User{}
	_, err = Model(&user).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}

func TestSessionUseSlaveAndPrepare (t *testing.T) {
	InitEngine4Test()
	UsePrepare(true)
	db_slave, err := conf.OpenSalveDb()
	if err != nil {
		t.Error(err)
		return
	}
	AddSlave(db_slave, "slave1", 1)

	user := User{Age: 31, DB_name: "model111"}
	_, err = Model(&user).Where("id=?", 1).Update()
	if err != nil{
		t.Error(err)
		return
	}

	user = User{}
	_, err = Model(&user).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)

	user = User{Age: 31, DB_name: "model222"}
	_, err = Model(&user).Where("id=?", 1).Update()
	if err != nil{
		t.Error(err)
		return
	}

	user = User{}
	_, err = Model(&user).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}