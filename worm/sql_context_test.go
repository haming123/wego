package worm

import (
	"wego/worm/conf"
	"testing"
)

func TestUseMaster (t *testing.T) {
	InitEngine4Test()
	db_slave, err := conf.OpenSalveDb()
	if err != nil {
		t.Error(err)
		return
	}
	AddSlave(db_slave, "slave1", 1)

	user := User{}
	_, err = Model(&user).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)

	user = User{}
	_, err = Model(&user).UseMaster(true).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}

func TestUsePrepare (t *testing.T) {
	InitEngine4Test()
	UsePrepare(false)

	user := User{}
	_, err := Model(&user).UsePrepare(true).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)

	user = User{}
	_, err = Model(&user).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}

func TestShowLog (t *testing.T) {
	InitEngine4Test()

	user := User{}
	_, err := Model(&user).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)

	user = User{}
	_, err = Model(&user).ShowLog(false).Where("id=?", 1).Get()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}
