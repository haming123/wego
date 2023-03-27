### 主要特点
gows是一款方便易用的Go语言websocket库，gows使用简单，能够支持大量客户端连接。具体特征如下：
* 采用更加经济的内存分配机制，使得每台服务器可接入更多的客户端。
* 内置消息接收循环，简化了消息的接收以及处理。
* 支持permessage-deflate压缩，支持设置压缩阈值，只有大于阈值的消息才会压缩发送。
* 支持Close handshake，可以优雅地关闭websocket连接。
* 支持并发消息发送。

### 安装
go get github.com/haming123/wego/gows

### 快速上手
* 注册页面路由，启动web服务
```go
package main
import (
	"github.com/haming123/wego/gows"
	"html/template"
	"net/http"
)
func main() {
	http.HandleFunc("/ws", HandlerWebSocket)
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")
		t, _ := template.ParseFiles("./index.html")
		t.Execute(w, user)
	})
	http.ListenAndServe(":8080", nil)
}
```
其中：
1）/index是一个web页面，在该页面中连接websocket，发送并接收websocket消息。
2）在/ws页面中进行websocket握手并升级为websocket协议。

* 定义websocket连接对应的结构体
在进行websocket握手前需要首先实现一个代表websocket连接的结构体，该结构体用于websocket事件的处理。该结构体需要实现以下接口方法：
```go
OnClose(ws *WebSocket)
OnMessage(ws *WebSocket, opcode int, buff *ByteBuffer) error
```
当接收的消息时会调用OnMessage方法，消息的数据可以从buff中获取。websocket连接关闭时会调用OnClose方法。
```go
type Client struct {
    ws   *gows.WebSocket
    user string
}

func (c *Client) OnClose(ws *gows.WebSocket) {
    log.Printf("OnClose: %s\n", c.user)
    c.ws = nil
}

func (c *Client) OnMessage(ws *gows.WebSocket, opcode int, vbuff *gows.ByteBuffer) error {
    log.Println("收到消息：", vbuff.GetString())
    c.ws.WriteString(vbuff.GetString())
    return nil
}
```

* 实现websocket握手处理器函数，将http连接升级为websocket连接
```go
func HandlerWebSocket(w http.ResponseWriter, r *http.Request) {
    user := r.FormValue("user")
    ws, err := gows.Accept(w, r, nil, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    client := &Client{ws: ws, user: user}
    ws.ServeRead(client)
}
```
其中的Client结构体代表客户端链接，必须实现OnMessage方法以及OnClose方法。创建Client对象后调用ws.Serve来启动消息接收循环。

### 使用wego的web框架实现webSocket握手
```go
func HandlerWebSocket(c *wego.WebContext) {
    user := c.Param.MustString("user")
    ws, err := c.AcceptWebsocket(nil, nil)
    if err != nil {
        c.AbortWithError(500, err)
        return
    }
    client := &Client{ws: ws, user: user}
    ws.ServeRead(client)
}
```

### 发送消息
* 发送文本消息
若要发送文本格式的消息，请调用WriteText或WriteString函数：
```go
func (ws *WebSocket) WriteText(data []byte) error {
    ......
}
func (ws *WebSocket) WriteString(data string) error {
    ......
}
```

* 若要发送二进制消息，使用WriteBinary方法
```go
func (ws *WebSocket) WriteBinary(data []byte) error {
    ......
}
```

* 发送JSON数据
```go
type Book struct {
    Name  string
    Price float64
}

err := ws.WriteJSON(Book{"Golang", 30.2})
if err != nil {
    log.Error(err)
    return
}
```

* 通过WriteCloser接口发送消息
```go
writer := ws.NextWriter(gows.Frame_Text)
defer writer.Close()

_, err := writer.WriteString("hello")
if err != nil {
    log.Error(err)
    return
}

_, err := writer.WriteString(" world")
if err != nil {
    log.Error(err)
    return
}
```

### 启用消息压缩处理
```go
package main
import (
	"github.com/haming123/wego/gows"
	"html/template"
	"net/http"
)
func main() {
	//开启发送压缩
	gows.UseFlate()
	//设置压缩的阈值，只有大于阈值的消息才会被压缩
	gows.SetMinCompressSize(512)
	http.HandleFunc("/ws", HandlerWebSocket)
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")
		t, _ := template.ParseFiles("./index.html")
		t.Execute(w, user)
	})
	http.ListenAndServe(":8080", nil)
}
```

### 跨域处理函数
gows缺省不支持跨域访问，若需要开启跨域访问，则需要自定义跨域处理函数（通过：SetOriginCheckFunc来设置）。例如:
```go
package main
import (
	"github.com/haming123/wego/gows"
	"html/template"
	"net/http"
)
func main() {
	//允许跨域访问
	gows.SetOriginCheckFunc(func(r *http.Request) bool {
		return true
	})
	http.HandleFunc("/ws", HandlerWebSocket)
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		user := r.FormValue("user")
		t, _ := template.ParseFiles("./index.html")
		t.Execute(w, user)
	})
	http.ListenAndServe(":8080", nil)
}
```

### 关闭客户端链接
关闭websocket时需要进行关闭协商：接任一端想关闭websocket，就发一个close frame给对端，对端收到该frame后，若之前没有发过close frame，则必须回复一个close frame。
采用协商方式关闭websocket连接，请调用CloseHandshake方法：
```go
func (ws *WebSocket) CloseHandshake(code CloseCode, text string) error {
    ......
}
```
CloseHandshake方法中的code参数为webSocket关闭状态码。webSocket关闭状态码的列表如下所示：
```go
const (
//正常关闭; 无论为何目的而创建, 该链接都已成功完成任务.
CloseNormalClosure CloseCode = 1000
//终端离开：可能因为服务端错误, 也可能因为浏览器正从打开连接的页面跳转离开.
CloseGoingAway CloseCode = 1001
//协议错误：由于协议错误而中断连接.
CloseProtocolError CloseCode = 1002
//数据格式错误：由于接收到不允许的数据类型而断开连接
CloseUnsupportedData CloseCode = 1003
//保留
CloseReserved CloseCode = 1004
//没有收到预期的状态码.
CloseNoCloseRcvd CloseCode = 1005
//异常关闭：用于期望收到状态码时连接非正常关闭 (也就是说, 没有发送关闭帧).
CloseAbnormalClosure CloseCode = 1006
//由于收到了格式不符的数据而断开连接 (如文本消息中包含了非 UTF-8 数据).
CloseInvalidPayload CloseCode = 1007
//由于收到不符合约定的数据而断开连接.
ClosePolicyViolation CloseCode = 1008
//由于收到过大的数据帧而断开连接.
CloseMessageTooBig CloseCode = 1009
//缺少扩展：客户端终止连接，因为期望一个或多个拓展, 但服务器没有.
CloseMandatoryExtension CloseCode = 1010
//内部错误：服务器终止连接，因为遇到异常
CloseInternalError CloseCode = 1011
//服务重启：服务器由于重启而断开连接.
CloseServiceRestart CloseCode = 1012
//稍后再试：服务器由于临时原因断开连接。
CloseTryAgainLater CloseCode = 1013
//错误的网关.
CloseBadGateway CloseCode = 1014
//握手错误：表示连接由于无法完成 TLS 握手而关闭 (例如无法验证服务器证书).
CloseTLSHandshake CloseCode = 1015
)
```
