package main

import (
	"sync"
)

type ChatGroup struct {
	muxLock sync.Mutex
	clients map[string]*Client
}

var user_map = ChatGroup{
	clients: make(map[string]*Client),
}

func AddClient(uuid string, client *Client) {
	user_map.muxLock.Lock()
	user_map.clients[uuid] = client
	user_map.muxLock.Unlock()
}

func BroadcastMeesage(data []byte) {
	user_map.muxLock.Lock()
	for _, client := range user_map.clients {
		if client.ws != nil {
			client.ws.WriteText(data)
		}
	}
	user_map.muxLock.Unlock()
}
