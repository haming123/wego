package wego

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"testing"
)

//go test -v -run=none -bench="BenchmarkWriteJsonWego" -benchmem
func BenchmarkWriteJsonWego(B *testing.B) {
	web, _ := NewWeb()
	web.Config.ShowUrlLog = false
	web.initModule()

	web.GET("/json", func(c *WebContext) {
		var user User
		user.Name = "lisi"
		user.Age = 12
		c.WriteJSON(200, user)
	})

	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		req, err := http.NewRequest("GET", "/json", nil)
		if err != nil {
			B.Log(err)
			return
		}
		web.ServeHTTP(w, req)
	}
}

//go test -v -run=none -bench="BenchmarkWriteJsonGolang" -benchmem
func BenchmarkWriteJsonGolang(B *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/json", func (w http.ResponseWriter, r *http.Request) {
		var user User
		user.Name = "lisi"
		user.Age = 12
		data, err := json.Marshal(user)
		if err != nil {
			B.Log(err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(data)
	})

	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		req, err := http.NewRequest("GET", "/json", nil)
		if err != nil {
			B.Log(err)
			return
		}
		mux.ServeHTTP(w, req)
	}
}

//go test -v -run=none -bench="BenchmarkWriteHTMLWego" -benchmem
func BenchmarkWriteHTMLWego(B *testing.B) {
	web, _ := NewWeb()
	web.Config.ShowUrlLog = false
	web.initModule()

	web.GET("/html", func(c *WebContext) {
		c.WriteHTML(200, "./test.html", "demo")
	})

	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		req, err := http.NewRequest("GET", "/html", nil)
		if err != nil {
			B.Log(err)
			return
		}
		web.ServeHTTP(w, req)
	}
}

//go test -v -run=none -bench="BenchmarkWriteHTMLGolang" -benchmem
func BenchmarkWriteHTMLGolang(B *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/html", func (w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./test.html")
		if err!=nil {
			B.Log(err)
			return
		}
		t.Execute(w, "demo")
	})

	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		req, err := http.NewRequest("GET", "/html", nil)
		if err != nil {
			B.Log(err)
			return
		}
		mux.ServeHTTP(w, req)
	}
}

//go test -v -run=none -bench="BenchmarkGetParamWego" -benchmem
func BenchmarkGetParamWego(B *testing.B) {
	web, _ := NewWeb()
	web.Config.ShowUrlLog = false
	web.initModule()

	web.GET("/text", func(c *WebContext) {
		name := c.QueryParam.GetString("name")
		c.QueryParam.GetString("name_9871")
		c.QueryParam.GetString("name_9872")
		c.QueryParam.GetString("name_9873")
		c.QueryParam.GetString("name_9874")
		c.QueryParam.GetString("name_9875")
		c.QueryParam.GetString("name_9876")
		c.QueryParam.GetString("name_9877")
		c.QueryParam.GetString("name_9878")
		c.QueryParam.GetString("name_9879")
		c.WriteText(200, name.Value)
	})

	var buf bytes.Buffer
	buf.WriteString("/text?name=lisi")
	buf.WriteString("&name_9871=value_9871")
	buf.WriteString("&name_9872=value_9872")
	buf.WriteString("&name_9873=value_9873")
	buf.WriteString("&name_9874=value_9874")
	buf.WriteString("&name_9875=value_9875")
	buf.WriteString("&name_9876=value_9876")
	buf.WriteString("&name_9877=value_9877")
	buf.WriteString("&name_9878=value_9878")
	buf.WriteString("&name_9879=value_9879")
	query := buf.String()

	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		req, err := http.NewRequest("GET", query, nil)
		if err != nil {
			B.Log(err)
			return
		}
		web.ServeHTTP(w, req)
	}
}

//go test -v -run=none -bench="BenchmarkGetParamGolang" -benchmem
func BenchmarkGetParamGolang(B *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/text", func (w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		r.FormValue("name_9871")
		r.FormValue("name_9872")
		r.FormValue("name_9873")
		r.FormValue("name_9874")
		r.FormValue("name_9875")
		r.FormValue("name_9876")
		r.FormValue("name_9877")
		r.FormValue("name_9878")
		r.FormValue("name_9879")
		w.Write([]byte(name))
	})

	var buf bytes.Buffer
	buf.WriteString("/text?name=lisi")
	buf.WriteString("&name_9871=value_9871")
	buf.WriteString("&name_9872=value_9872")
	buf.WriteString("&name_9873=value_9873")
	buf.WriteString("&name_9874=value_9874")
	buf.WriteString("&name_9875=value_9875")
	buf.WriteString("&name_9876=value_9876")
	buf.WriteString("&name_9877=value_9877")
	buf.WriteString("&name_9878=value_9878")
	buf.WriteString("&name_9879=value_9879")
	query := buf.String()

	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		req, err := http.NewRequest("GET", query, nil)
		if err != nil {
			B.Log(err)
			return
		}
		mux.ServeHTTP(w, req)
	}
}

//go test -v -run=none -bench="BenchmarkPostParamWego" -benchmem
func BenchmarkPostParamWego(B *testing.B) {
	web, _ := NewWeb()
	web.Config.ShowUrlLog = false
	web.initModule()

	web.POST("/text", func(c *WebContext) {
		name := c.QueryParam.GetString("name")
		c.FormParam.GetString("name_9871")
		c.FormParam.GetString("name_9872")
		c.FormParam.GetString("name_9873")
		c.FormParam.GetString("name_9874")
		c.FormParam.GetString("name_9875")
		c.FormParam.GetString("name_9876")
		c.FormParam.GetString("name_9877")
		c.FormParam.GetString("name_9878")
		c.FormParam.GetString("name_9879")
		c.WriteText(200, name.Value)
	})

	var buf bytes.Buffer
	buf.WriteString("&name_9871=value_9871")
	buf.WriteString("&name_9872=value_9872")
	buf.WriteString("&name_9873=value_9873")
	buf.WriteString("&name_9874=value_9874")
	buf.WriteString("&name_9875=value_9875")
	buf.WriteString("&name_9876=value_9876")
	buf.WriteString("&name_9877=value_9877")
	buf.WriteString("&name_9878=value_9878")
	buf.WriteString("&name_9879=value_9879")

	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		req, err := http.NewRequest("POST", "/text?name=lisi", &buf)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset:utf-8;")
		if err != nil {
			B.Log(err)
			return
		}
		web.ServeHTTP(w, req)
	}
}

//go test -v -run=none -bench="BenchmarkPostParamGolang" -benchmem
func BenchmarkPostParamGolang(B *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/text", func (w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		r.Form.Get("name_9871")
		r.Form.Get("name_9872")
		r.Form.Get("name_9873")
		r.Form.Get("name_9874")
		r.Form.Get("name_9875")
		r.Form.Get("name_9876")
		r.Form.Get("name_9877")
		r.Form.Get("name_9878")
		r.Form.Get("name_9879")
		w.Write([]byte(name))
	})

	var buf bytes.Buffer
	buf.WriteString("name_9871=value_9871")
	buf.WriteString("&name_9872=value_9872")
	buf.WriteString("&name_9873=value_9873")
	buf.WriteString("&name_9874=value_9874")
	buf.WriteString("&name_9875=value_9875")
	buf.WriteString("&name_9876=value_9876")
	buf.WriteString("&name_9877=value_9877")
	buf.WriteString("&name_9878=value_9878")
	buf.WriteString("&name_9879=value_9879")

	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		req, err := http.NewRequest("POST", "/text?name=lisi", &buf)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset:utf-8;")
		if err != nil {
			B.Log(err)
			return
		}
		mux.ServeHTTP(w, req)
	}
}