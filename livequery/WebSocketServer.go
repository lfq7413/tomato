package livequery

import (
	"net/http"

	"golang.org/x/net/websocket"
)

type webSocketHandler interface {
	onConnect(ws *websocket.Conn)
	onMessage(ws *websocket.Conn, msg interface{})
	onDisconnect(ws *websocket.Conn)
}

var handler webSocketHandler

func runWebSocketServer(pattern, addr string, h webSocketHandler) {
	handler = h
	http.Handle(pattern, websocket.Handler(httpHandler))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func httpHandler(ws *websocket.Conn) {
	handler.onConnect(ws)
	var v string
	for {
		err := receive(ws, &v)
		if err != nil {
			handler.onDisconnect(ws)
			return
		}
		handler.onMessage(ws, v)
	}
}

func receive(ws *websocket.Conn, v interface{}) error {
	return websocket.Message.Receive(ws, v)
}

func send(ws *websocket.Conn, v interface{}) error {
	return websocket.Message.Send(ws, v)
}
