package worm

import (
	"testing"
)

func TestWeightSlave(t *testing.T)  {
	db, err := OpenDb()
	if err != nil {
		t.Error(err)
		return
	}

	eng, err := NewEngine(&dialectMysql{}, db)
	if err != nil {
		t.Error(err)
		return
	}

	eng.AddSlave(db, "db2", 1)
	eng.AddSlave(db, "db3", 2)
	ret := getSlaveByWeight(eng.slaves)
	t.Log(ret.db_name)
}
