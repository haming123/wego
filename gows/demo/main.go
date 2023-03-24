package main

import (
	"github.com/haming123/wego/gows"
	"html/template"
	"log"
	"net/http"
)

func main() {
	//开启发送压缩
	gows.UseFlate()
	//设置压缩的阈值，只有大于阈值的消息才会被压缩
	gows.SetMinCompressSize(512)

	http.HandleFunc("/conn", HandlerWebSocket)
	//在浏览器一个标签页中输入：http://127.0.0.1:8080/index?user=a
	//在浏览器另外一个标签页中输入：http://127.0.0.1:8080/index?user=b
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")
		t, _ := template.ParseFiles("./index.html")
		t.Execute(w, user)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
