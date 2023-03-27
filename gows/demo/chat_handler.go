package main

import (
	"github.com/haming123/wego/gows"
	"log"
	"net/http"
)

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
	BroadcastMeesage(vbuff.GetBytes())
	return nil
}

func HandlerWebSocket(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	ws, err := gows.Accept(w, r, nil, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &Client{ws: ws, user: user}
	AddClient(user, client)
	ws.ServeRead(client)
}
