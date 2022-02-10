package wego

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	log "wego/dlog"
)

func CreateTestContext(w http.ResponseWriter, req *http.Request) (*WebEngine, *WebContext) {
	engine := newEngine()
	c := engine.ctxPool.Get().(*WebContext)
	c.reset()
	c.engine = engine
	c.Input.Request = req
	c.Output.ResponseWriter = w
	c.Path = req.URL.Path
	return engine, c
}

func TestContextGetStruct4Query(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/?name=lisi&Age=22&Str=111", bytes.NewBufferString("name=zhangsan"))
	_, c := CreateTestContext(w, req)

	var user struct {
		Name 	string 		`form:"name"`
		Age 	int
		Def 	time.Time	`form:"birth;default=2018-02-01"`
		Ptr		*string 	`form:"name"`
		tmp		int
	}
	err := c.QueryParam.GetStruct(&user)
	if err != nil {
		t.Error(err)
	}
	t.Log(log.JsonMarshal(user))
}

func TestContextGetStruct4Error(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/?name=lisi&Age=22&Str=111", bytes.NewBufferString("name=zhangsan"))
	_, c := CreateTestContext(w, req)

	var user struct {
		Name 	string 		`form:"name"`
		Age 	int
		Str		sql.NullString
	}
	err := c.QueryParam.GetStruct(&user)
	if err != nil {
		t.Log(err)
	}
}

func TestContextGetStruct4Param(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/?name=lisi&Age=22", bytes.NewBufferString("name=zhangsan"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, c := CreateTestContext(w, req)

	var user struct {
		Name 	[]string 	`form:"name"`
		Age 	int
	}
	err := c.Param.GetStruct(&user)
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
}