package wego

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"
)

/*
go test -v -run=none -bench="BenchmarkGetRoute" -benchmem
pkg:
//参数路由
BenchmarkGetRoute-6     17706348                67.59 ns/op            0 B/op          0 allocs/op
BenchmarkGetRoute-6     14551132                82.11 ns/op            0 B/op          0 allocs/op
//固定路由
BenchmarkGetRoute-6     18674397                64.41 ns/op            0 B/op          0 allocs/op
BenchmarkGetRoute-6     16346502                71.50 ns/op            0 B/op          0 allocs/op

pkg:
//参数路由
Benchmark404Many-6      17321390                72.51 ns/op            0 B/op          0 allocs/op
Benchmark404Many-6      14443422                82.05 ns/op            0 B/op          0 allocs/op
//固定路由
Benchmark404Many-6      14622055                80.74 ns/op            0 B/op          0 allocs/op
Benchmark404Many-6      13066329                89.46 ns/op            0 B/op          0 allocs/op
*/
func BenchmarkGetRoute(B *testing.B) {
	web, _ := NewWeb()
	web.Config.ShowUrlLog = false
	web.SetHandler404(func(c *WebContext) {})

	/*
	for i:=0; i < 9; i++ {
		no := rand.Intn(100)
		part1 := fmt.Sprintf("/part1_%d", no)
		for j:=0; j < 10; j++ {
			no2 := rand.Intn(10000)
			part2 := fmt.Sprintf("/part2_%d", no2)
			fmt.Printf("web.GET(\"%s/%s\", func(c *WebContext) {})\n", part1, part2)
		}
	}*/

	web.GET("/", func(c *WebContext) {})
	web.GET("/path_82", func(c *WebContext) {})
	web.GET("/path_887", func(c *WebContext) {})
	web.GET("/path_847", func(c *WebContext) {})
	web.GET("/path_59", func(c *WebContext) {})
	web.GET("/path_81", func(c *WebContext) {})
	web.GET("/path_318", func(c *WebContext) {})
	web.GET("/path_425", func(c *WebContext) {})
	web.GET("/path_540", func(c *WebContext) {})
	web.GET("/path_456", func(c *WebContext) {})
	web.GET("/path_300", func(c *WebContext) {})
	web.GET("/part1_81/part2_7887", func(c *WebContext) {})
	web.GET("/part1_81/part2_1847", func(c *WebContext) {})
	web.GET("/part1_81/part2_4059", func(c *WebContext) {})
	web.GET("/part1_81/part2_2081", func(c *WebContext) {})
	web.GET("/part1_81/part2_1318", func(c *WebContext) {})
	web.GET("/part1_81/part2_4425", func(c *WebContext) {})
	web.GET("/part1_81/part2_2540", func(c *WebContext) {})
	web.GET("/part1_81/part2_456", func(c *WebContext) {})
	web.GET("/part1_81/part2_3300", func(c *WebContext) {})
	web.GET("/part1_81/part2_694", func(c *WebContext) {})
	web.GET("/part1_11/part2_8162", func(c *WebContext) {})
	web.GET("/part1_11/part2_5089", func(c *WebContext) {})
	web.GET("/part1_11/part2_4728", func(c *WebContext) {})
	web.GET("/part1_11/part2_3274", func(c *WebContext) {})
	web.GET("/part1_11/part2_1211", func(c *WebContext) {})
	web.GET("/part1_11/part2_1445", func(c *WebContext) {})
	web.GET("/part1_11/part2_3237", func(c *WebContext) {})
	web.GET("/part1_11/part2_9106", func(c *WebContext) {})
	web.GET("/part1_11/part2_495", func(c *WebContext) {})
	web.GET("/part1_11/part2_5466", func(c *WebContext) {})
	web.GET("/part1_28/part2_6258", func(c *WebContext) {})
	web.GET("/part1_28/part2_8047", func(c *WebContext) {})
	web.GET("/part1_28/part2_9947", func(c *WebContext) {})
	web.GET("/part1_28/part2_8287", func(c *WebContext) {})
	web.GET("/part1_28/part2_2888", func(c *WebContext) {})
	web.GET("/part1_28/part2_2790", func(c *WebContext) {})
	web.GET("/part1_28/part2_3015", func(c *WebContext) {})
	web.GET("/part1_28/part2_5541", func(c *WebContext) {})
	web.GET("/part1_28/part2_408", func(c *WebContext) {})
	web.GET("/part1_28/part2_7387", func(c *WebContext) {})
	web.GET("/part1_31/part2_5429", func(c *WebContext) {})
	web.GET("/part1_31/part2_5356", func(c *WebContext) {})
	web.GET("/part1_31/part2_1737", func(c *WebContext) {})
	web.GET("/part1_31/part2_631", func(c *WebContext) {})
	web.GET("/part1_31/part2_1485", func(c *WebContext) {})
	web.GET("/part1_31/part2_5026", func(c *WebContext) {})
	web.GET("/part1_31/part2_6413", func(c *WebContext) {})
	web.GET("/part1_31/part2_3090", func(c *WebContext) {})
	web.GET("/part1_31/part2_5194", func(c *WebContext) {})
	web.GET("/part1_31/part2_563", func(c *WebContext) {})
	web.GET("/part1_33/part2_4147", func(c *WebContext) {})
	web.GET("/part1_33/part2_4078", func(c *WebContext) {})
	web.GET("/part1_33/part2_4324", func(c *WebContext) {})
	web.GET("/part1_33/part2_6159", func(c *WebContext) {})
	web.GET("/part1_33/part2_1353", func(c *WebContext) {})
	web.GET("/part1_33/part2_1957", func(c *WebContext) {})
	web.GET("/part1_33/part2_3721", func(c *WebContext) {})
	web.GET("/part1_33/part2_7189", func(c *WebContext) {})
	web.GET("/part1_33/part2_2199", func(c *WebContext) {})
	web.GET("/part1_33/part2_3000", func(c *WebContext) {})
	web.GET("/part1_5/part2_2888", func(c *WebContext) {})
	web.GET("/part1_5/part2_4538", func(c *WebContext) {})
	web.GET("/part1_5/part2_9703", func(c *WebContext) {})
	web.GET("/part1_5/part2_9355", func(c *WebContext) {})
	web.GET("/part1_5/part2_2451", func(c *WebContext) {})
	web.GET("/part1_5/part2_8510", func(c *WebContext) {})
	web.GET("/part1_5/part2_2605", func(c *WebContext) {})
	web.GET("/part1_5/part2_156", func(c *WebContext) {})
	web.GET("/part1_5/part2_8266", func(c *WebContext) {})
	web.GET("/part1_5/part2_9828", func(c *WebContext) {})
	web.GET("/part1_61/part2_7202", func(c *WebContext) {})
	web.GET("/part1_61/part2_4783", func(c *WebContext) {})
	web.GET("/part1_61/part2_5746", func(c *WebContext) {})
	web.GET("/part1_61/part2_1563", func(c *WebContext) {})
	web.GET("/part1_61/part2_4376", func(c *WebContext) {})
	web.GET("/part1_61/part2_9002", func(c *WebContext) {})
	web.GET("/part1_61/part2_9718", func(c *WebContext) {})
	web.GET("/part1_61/part2_5447", func(c *WebContext) {})
	web.GET("/part1_61/part2_5094", func(c *WebContext) {})
	web.GET("/part1_61/part2_1577", func(c *WebContext) {})
	web.GET("/part1_63/part2_7996", func(c *WebContext) {})
	web.GET("/part1_63/part2_6420", func(c *WebContext) {})
	web.GET("/part1_63/part2_8623", func(c *WebContext) {})
	web.GET("/part1_63/part2_953", func(c *WebContext) {})
	web.GET("/part1_63/part2_1137", func(c *WebContext) {})
	web.GET("/part1_63/part2_3133", func(c *WebContext) {})
	web.GET("/part1_63/part2_9241", func(c *WebContext) {})
	web.GET("/part1_63/part2_59", func(c *WebContext) {})
	web.GET("/part1_63/part2_3033", func(c *WebContext) {})
	web.GET("/part1_63/part2_8643", func(c *WebContext) {})
	web.GET("/part1_91/part2_2002", func(c *WebContext) {})
	web.GET("/part1_91/part2_8878", func(c *WebContext) {})
	web.GET("/part1_91/part2_9336", func(c *WebContext) {})
	web.GET("/part1_91/part2_2546", func(c *WebContext) {})
	web.GET("/part1_91/part2_9107", func(c *WebContext) {})
	web.GET("/part1_91/part2_7940", func(c *WebContext) {})
	web.GET("/part1_91/part2_6503", func(c *WebContext) {})
	web.GET("/part1_91/part2_552", func(c *WebContext) {})
	web.GET("/part1_91/part2_9843", func(c *WebContext) {})
	web.GET("/hello/:id", func(c *WebContext) {})

	web.initModule()
	runRequest(B, web, "/part1_91/part2_9843")
}

/*
$ go test -v -run=none -bench="BenchmarkWriteText" -benchmem
cd /d/GoApp/src/github.com/gin-gonic/gin
cd /d/GoApp/src/wego
pkg: wego
BenchmarkWriteText-6     		 8584384               143.4 ns/op            16 B/op          1 allocs/op
pkg:
BenchmarkOneRouteString-6        7704174               159.4 ns/op            48 B/op          1 allocs/op
*/
func BenchmarkWriteText(B *testing.B) {
	web, _ := NewWeb()
	web.Config.ShowUrlLog = false
	web.GET("/text", func(c *WebContext) {
		c.WriteText(http.StatusOK, "this is a plain text")
	})

	web.initModule()
	runRequest(B, web,  "/text")
}

/*
$ go test -v -run=none -bench="BenchmarkWriteText2" -benchmem
pkg: wego
BenchmarkWriteText2-6            1605830               745.3 ns/op           480 B/op          7 allocs/op
pkg: wego2
BenchmarkWriteText2-6            5060521               237.7 ns/op            16 B/op          1 allocs/op
pkg:
BenchmarkWriteText2-6            1437736               830.4 ns/op           544 B/op          8 allocs/op
*/
func BenchmarkWriteText2(B *testing.B) {
	web, _ := NewWeb()
	web.Config.ShowUrlLog = false
	web.GET("/text", func(c *WebContext) {
		name := c.QueryParam.GetString("name")
		c.WriteTextF(200, "%s", name.Value)
	})

	web.initModule()
	runRequest(B, web,  "/text?name=lisi&age=12")
}

//go test -v -run=none -bench="BenchmarkQueryParam" -benchmem
func BenchmarkQueryParam(B *testing.B) {
	web, _ := NewWeb()
	web.Config.ShowUrlLog = false
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
		c.QueryParam.GetString("name_9879")/*
		c.QueryParam.GetString("name_9881")
		c.QueryParam.GetString("name_9882")
		c.QueryParam.GetString("name_9883")
		c.QueryParam.GetString("name_9884")
		c.QueryParam.GetString("name_9885")
		c.QueryParam.GetString("name_9886")
		c.QueryParam.GetString("name_9887")
		c.QueryParam.GetString("name_9888")
		c.QueryParam.GetString("name_9889")
		c.QueryParam.GetString("name_9171")
		c.QueryParam.GetString("name_9172")
		c.QueryParam.GetString("name_9173")
		c.QueryParam.GetString("name_9174")
		c.QueryParam.GetString("name_9175")
		c.QueryParam.GetString("name_9176")
		c.QueryParam.GetString("name_9177")
		c.QueryParam.GetString("name_9178")
		c.QueryParam.GetString("name_9179")
		c.QueryParam.GetString("name_9281")
		c.QueryParam.GetString("name_9282")
		c.QueryParam.GetString("name_9283")
		c.QueryParam.GetString("name_9284")
		c.QueryParam.GetString("name_9285")
		c.QueryParam.GetString("name_9286")
		c.QueryParam.GetString("name_9287")
		c.QueryParam.GetString("name_9288")
		c.QueryParam.GetString("name_9289")
		c.QueryParam.GetString("name_9371")
		c.QueryParam.GetString("name_9372")
		c.QueryParam.GetString("name_9373")
		c.QueryParam.GetString("name_9374")
		c.QueryParam.GetString("name_9375")
		c.QueryParam.GetString("name_9376")
		c.QueryParam.GetString("name_9377")
		c.QueryParam.GetString("name_9378")
		c.QueryParam.GetString("name_9379")
		c.QueryParam.GetString("name_9481")
		c.QueryParam.GetString("name_9482")
		c.QueryParam.GetString("name_9483")
		c.QueryParam.GetString("name_9484")
		c.QueryParam.GetString("name_9485")
		c.QueryParam.GetString("name_9486")
		c.QueryParam.GetString("name_9487")
		c.QueryParam.GetString("name_9488")
		c.QueryParam.GetString("name_9489")
		c.QueryParam.GetString("name_9571")
		c.QueryParam.GetString("name_9572")
		c.QueryParam.GetString("name_9573")
		c.QueryParam.GetString("name_9574")
		c.QueryParam.GetString("name_9575")
		c.QueryParam.GetString("name_9576")
		c.QueryParam.GetString("name_9577")
		c.QueryParam.GetString("name_9578")
		c.QueryParam.GetString("name_9579")
		c.QueryParam.GetString("name_9681")
		c.QueryParam.GetString("name_9682")
		c.QueryParam.GetString("name_9683")
		c.QueryParam.GetString("name_9684")
		c.QueryParam.GetString("name_9685")
		c.QueryParam.GetString("name_9686")
		c.QueryParam.GetString("name_9687")
		c.QueryParam.GetString("name_9688")
		c.QueryParam.GetString("name_9689")*/
		c.WriteTextF(200, "%s", name.Value)
	})

	web.initModule()
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
	buf.WriteString("&name_9879=value_9879")/*
	buf.WriteString("&name_9881=value_9871")
	buf.WriteString("&name_9882=value_9872")
	buf.WriteString("&name_9883=value_9873")
	buf.WriteString("&name_9884=value_9874")
	buf.WriteString("&name_9885=value_9875")
	buf.WriteString("&name_9886=value_9876")
	buf.WriteString("&name_9887=value_9877")
	buf.WriteString("&name_9888=value_9878")
	buf.WriteString("&name_9889=value_9879")
	buf.WriteString("&name_9171=value_9871")
	buf.WriteString("&name_9172=value_9872")
	buf.WriteString("&name_9173=value_9873")
	buf.WriteString("&name_9174=value_9874")
	buf.WriteString("&name_9175=value_9875")
	buf.WriteString("&name_9176=value_9876")
	buf.WriteString("&name_9177=value_9877")
	buf.WriteString("&name_9178=value_9878")
	buf.WriteString("&name_9179=value_9879")
	buf.WriteString("&name_9281=value_9871")
	buf.WriteString("&name_9282=value_9872")
	buf.WriteString("&name_9283=value_9873")
	buf.WriteString("&name_9284=value_9874")
	buf.WriteString("&name_9285=value_9875")
	buf.WriteString("&name_9286=value_9876")
	buf.WriteString("&name_9287=value_9877")
	buf.WriteString("&name_9288=value_9878")
	buf.WriteString("&name_9289=value_9879")
	buf.WriteString("&name_9371=value_9871")
	buf.WriteString("&name_9372=value_9872")
	buf.WriteString("&name_9373=value_9873")
	buf.WriteString("&name_9374=value_9874")
	buf.WriteString("&name_9375=value_9875")
	buf.WriteString("&name_9376=value_9876")
	buf.WriteString("&name_9377=value_9877")
	buf.WriteString("&name_9378=value_9878")
	buf.WriteString("&name_9379=value_9879")
	buf.WriteString("&name_9481=value_9871")
	buf.WriteString("&name_9482=value_9872")
	buf.WriteString("&name_9483=value_9873")
	buf.WriteString("&name_9484=value_9874")
	buf.WriteString("&name_9485=value_9875")
	buf.WriteString("&name_9486=value_9876")
	buf.WriteString("&name_9487=value_9877")
	buf.WriteString("&name_9488=value_9878")
	buf.WriteString("&name_9489=value_9879")
	buf.WriteString("&name_9571=value_9871")
	buf.WriteString("&name_9572=value_9872")
	buf.WriteString("&name_9573=value_9873")
	buf.WriteString("&name_9574=value_9874")
	buf.WriteString("&name_9575=value_9875")
	buf.WriteString("&name_9576=value_9876")
	buf.WriteString("&name_9577=value_9877")
	buf.WriteString("&name_9578=value_9878")
	buf.WriteString("&name_9579=value_9879")
	buf.WriteString("&name_9681=value_9871")
	buf.WriteString("&name_9682=value_9872")
	buf.WriteString("&name_9683=value_9873")
	buf.WriteString("&name_9684=value_9874")
	buf.WriteString("&name_9685=value_9875")
	buf.WriteString("&name_9686=value_9876")
	buf.WriteString("&name_9687=value_9877")
	buf.WriteString("&name_9688=value_9878")
	buf.WriteString("&name_9689=value_9879")*/

	runGetRequest(B, web, buf.String())
}

//go test -v -run=none -bench="BenchmarkFormParam" -benchmem
func BenchmarkFormParam(B *testing.B) {
	web, _ := NewWeb()
	web.Config.ShowUrlLog = false
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
		c.FormParam.GetString("name_9881")
		c.FormParam.GetString("name_9882")
		c.FormParam.GetString("name_9883")
		c.FormParam.GetString("name_9884")
		c.FormParam.GetString("name_9885")
		c.FormParam.GetString("name_9886")
		c.FormParam.GetString("name_9887")
		c.FormParam.GetString("name_9888")
		c.FormParam.GetString("name_9889")
		c.WriteTextF(200, "%s", name.Value)
	})

	web.initModule()
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
	buf.WriteString("&name_9881=value_9871")
	buf.WriteString("&name_9882=value_9872")
	buf.WriteString("&name_9883=value_9873")
	buf.WriteString("&name_9884=value_9874")
	buf.WriteString("&name_9885=value_9875")
	buf.WriteString("&name_9886=value_9876")
	buf.WriteString("&name_9887=value_9877")
	buf.WriteString("&name_9888=value_9878")
	buf.WriteString("&name_9889=value_9879")

	runPostRequest(B, web, "/text?name=lisi", buf.Bytes(), "application/x-www-form-urlencoded; charset:utf-8;")
}

func BenchmarkMultipart(B *testing.B) {
	web, _ := NewWeb()
	web.Config.ShowUrlLog = false
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
		c.WriteTextF(200, "%s", name.Value)
	})

	web.initModule()
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	mw.WriteField("name_9871", "value_9871")
	mw.WriteField("name_9872", "value_9872")
	mw.WriteField("name_9873", "value_9873")
	mw.WriteField("name_9874", "value_9874")
	mw.WriteField("name_9875", "value_9875")
	mw.WriteField("name_9876", "value_9876")
	mw.WriteField("name_9877", "value_9877")
	mw.WriteField("name_9878", "value_9878")
	mw.WriteField("name_9879", "value_9879")
	mw.WriteField("name_9881", "value_9871")
	mw.WriteField("name_9882", "value_9872")
	mw.WriteField("name_9883", "value_9873")
	mw.WriteField("name_9884", "value_9874")
	mw.WriteField("name_9885", "value_9875")
	mw.WriteField("name_9886", "value_9876")
	mw.WriteField("name_9887", "value_9877")
	mw.WriteField("name_9888", "value_9878")
	mw.WriteField("name_9889", "value_9879")
	mw.Close()
	data := buf.Bytes()
	runPostRequest(B, web, "/text?name=lisi", data, mw.FormDataContentType())
}

func runRequest(B *testing.B, web *WebEngine, path string) {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		B.Log(err)
		return
	}

	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		web.ServeHTTP(w, req)
	}
}

func runGetRequest(B *testing.B, web *WebEngine, path string) {
	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			B.Log(err)
			return
		}
		web.ServeHTTP(w, req)
	}
}

func runPostRequest(B *testing.B, web *WebEngine, path string, body []byte, ct string) {
	var buff bytes.Buffer
	buff.Write(body)

	w := newMockWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		req, err := http.NewRequest("POST", path, &buff)
		if ct != "" {
			req.Header.Set("Content-Type", ct)
		}
		if err != nil {
			B.Log(err)
			return
		}
		web.ServeHTTP(w, req)
		buff.Reset()
		buff.Write(body)
	}
}
