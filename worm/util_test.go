package worm

import (
	"testing"
)

func TestGetIntArray (t *testing.T) {
	arr := []User{}
	if 1==1 {
		user := User{DB_id:2, Age: 31, DB_name: "demo9"}
		arr = append(arr, user)
	}
	if 1==1 {
		user := User{DB_id:5, Age: 31, DB_name: "demo9"}
		arr = append(arr, user)
	}
	ret, err := CreateIntArray(arr, "DB_id")
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}