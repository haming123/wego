package wego

import (
	"testing"
	"time"
)

func TestSessionInfo(t *testing.T) {
	session := SessionInfo{}
	session.Set("name", "hello")
	session.Set("id", 333)
	session.Set("tm", time.Now())
	session.Set("ok", true)
	t.Log(session.GetString("name"))
	t.Log(session.GetInt("id"))
	t.Log(session.GetTime("tm"))
	t.Log(session.GetBool("ok"))
}

type User struct {
	Name string
	ID int64
	Age int
}

func TestSessionInfoStruct(t *testing.T) {
	user := User{}
	user.ID = 1
	user.Name = "lisi"
	user.Age = 12
	t.Log(user)

	session := SessionInfo{}
	session.Set("user", user)

	var user2 User
	err := session.GetStuct("user", &user2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user2)
}

