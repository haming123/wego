package wego

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPathParamSet(t *testing.T) {
	var params PathParam
	params.SetValue("name", "lisi")
	params.SetValue("Age", "22")
	t.Log(params.GetString("name"))
	t.Log(params.GetInt("Age"))
}

func TestGetParam(t *testing.T) {
	web, _ := NewWeb()
	web.GET("/user", func(c *WebContext) {
		name := c.Param.GetString("name")
		if name.Error != nil {
			t.Error(name.Error)
		}
		t.Log(name.Value)
		age := c.Param.GetInt("age")
		if age.Error != nil {
			t.Error(age.Error)
		}
		t.Log(age.Value)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user?name=lisi&age=12", nil)
	web.ServeHTTP(w, req)
}

func TestParamMust(t *testing.T) {
	web, _ := NewWeb()
	web.GET("/user", func(c *WebContext) {
		name := c.Param.MustString("name")
		t.Log(name)
		age := c.Param.MustInt("age")
		t.Log(age)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user?name=lisi&age=12", nil)
	web.ServeHTTP(w, req)
}

func TestGetPathParam(t *testing.T) {
	web, _ := NewWeb()
	web.GET("/hello/:id", func(c *WebContext) {
		ret := c.RouteParam.GetString("id")
		if ret.Error != nil {
			t.Error(ret.Error)
		}
		t.Log(ret.Value)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello/123", nil)
	web.ServeHTTP(w, req)
}

func TestGetPathParamMust(t *testing.T) {
	web, _ := NewWeb()
	web.GET("/hello/:id", func(c *WebContext) {
		val := c.RouteParam.MustInt("id", 0)
		t.Log(val)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello/123", nil)
	web.ServeHTTP(w, req)
}

func TestGetQueryParam(t *testing.T) {
	web, _ := NewWeb()
	web.GET("/user", func(c *WebContext) {
		name := c.QueryParam.GetString("name").Value
		age := c.QueryParam.GetInt("age").Value
		t.Log(name)
		t.Log(age)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user?name=lisi&age=12", nil)
	web.ServeHTTP(w, req)
}

func TestGetFromParam(t *testing.T) {
	web, _ := NewWeb()
	web.POST("/user", func(c *WebContext) {
		name := c.FormParam.GetString("name").Value
		age := c.FormParam.GetInt("age").Value
		t.Log(name)
		t.Log(age)
	})

	var buff bytes.Buffer
	buff.WriteString("name=lisi&age=12")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user", &buff)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	web.ServeHTTP(w, req)
}

type FormUser struct {
	Name 	string 		`form:"name"`
	Age 	int
}

func TestGetStructRouteParam(t *testing.T) {
	web, _ := NewWeb()
	web.POST("/user/:name/:Age", func(c *WebContext) {
		var user FormUser
		err := c.RouteParam.GetStruct(&user)
		if err != nil {
			t.Log(err)
		}
		t.Log(user)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user/lisi/12", nil)
	web.ServeHTTP(w, req)
}

func TestGetStructQueryParam(t *testing.T) {
	web, _ := NewWeb()
	web.POST("/user", func(c *WebContext) {
		var user FormUser
		err := c.QueryParam.GetStruct(&user)
		if err != nil {
			t.Log(err)
		}
		t.Log(user)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user?name=lisi&Age=12", nil)
	web.ServeHTTP(w, req)
}

func TestGetStructFormParam(t *testing.T) {
	web, _ := NewWeb()
	web.POST("/user", func(c *WebContext) {
		var user FormUser
		err := c.FormParam.GetStruct(&user)
		if err != nil {
			t.Log(err)
		}
		t.Log(user)
	})

	var buff bytes.Buffer
	buff.WriteString("name=lisi&Age=12")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user", &buff)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	web.ServeHTTP(w, req)
}

func TestGetStructParam(t *testing.T) {
	web, _ := NewWeb()
	web.POST("/user", func(c *WebContext) {
		var user FormUser
		err := c.Param.GetStruct(&user)
		if err != nil {
			t.Log(err)
		}
		t.Log(user)
	})

	var buff bytes.Buffer
	buff.WriteString("name=lisi&Age=12")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user", &buff)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	web.ServeHTTP(w, req)
}

func TestReadJson(t *testing.T) {
	web, _ := NewWeb()
	web.POST("/user", func(c *WebContext) {
		var user2 User
		err := c.ReadJSON(&user2)
		if err != nil {
			t.Log(err)
		}
		t.Log(user2)
	})

	user := User{}
	user.ID = 1
	user.Name = "lisi"
	user.Age = 12
	data, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user",  bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	web.ServeHTTP(w, req)
}

func TestReadXML(t *testing.T) {
	web, _ := NewWeb()
	web.POST("/user", func(c *WebContext) {
		var user2 User
		err := c.ReadXML(&user2)
		if err != nil {
			t.Log(err)
		}
		t.Log(user2)
	})

	user := User{}
	user.ID = 1
	user.Name = "lisi"
	user.Age = 12
	data, _ := xml.Marshal(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user",  bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	web.ServeHTTP(w, req)
}
