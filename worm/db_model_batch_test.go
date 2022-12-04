package worm

import "testing"

func TestBatchInsertMo(t *testing.T) {
	InitEngine4Test()

	users := []User{User{DB_name: "batch1", Age: 33}, User{DB_name: "batch2", Age: 33}}
	res, err := Model(&User{}).BatchInsert(&users)
	if err != nil {
		t.Error(err)
		return
	}
	num, _ := res.RowsAffected()
	t.Logf("Affected=%d", num)

	num, err = Model(&User{}).Where("age=?", 33).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("delete num=%d", num)
}

func TestBatchInsertVo(t *testing.T) {
	InitEngine4Test()

	users := []UserVo{UserVo{DB_name: "batch1", Age: 33}, UserVo{DB_name: "batch2", Age: 33}}
	res, err := Model(&User{}).BatchInsert(&users)
	if err != nil {
		t.Error(err)
		return
	}
	num, _ := res.RowsAffected()
	t.Logf("Affected=%d", num)

	num, err = Model(&User{}).Where("age=?", 33).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("delete num=%d", num)
}

func TestBatchInsertVo2(t *testing.T) {
	InitEngine4Test()

	users := []UserVo{UserVo{DB_name: "batch1", Age: 33}, UserVo{DB_name: "batch2", Age: 33}}
	res, err := Model(&User{}).BatchInsert(&users)
	if err != nil {
		t.Error(err)
		return
	}
	num, _ := res.RowsAffected()
	t.Logf("Affected=%d", num)

	num, err = Model(&User{}).Where("age=?", 33).Delete()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("delete num=%d", num)
}
