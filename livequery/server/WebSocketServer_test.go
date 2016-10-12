package server

import (
	"fmt"
	"testing"
)

func Test_WebSocketServer(t *testing.T) {
	h := handle{}
	RunWebSocketServer("/livequery", ":8089", h)
}

type handle struct{}

func (h handle) OnConnect(ws *WebSocket) {
	fmt.Println("OnConnect")
}

func (h handle) OnMessage(ws *WebSocket, msg interface{}) {
	fmt.Println("OnMessage", msg)
	ws.send(msg)
}

func (h handle) OnDisconnect(ws *WebSocket) {
	fmt.Println("OnDisconnect")
}
