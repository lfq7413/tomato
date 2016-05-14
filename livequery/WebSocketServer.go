package livequery

import (
	"net/http"

	"golang.org/x/net/websocket"
)

type webSocketHandler interface {
	onConnect(ws *webSocket)
	onMessage(ws *webSocket, msg interface{})
	onDisconnect(ws *webSocket)
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
	socket := &webSocket{
		ws:       ws,
		clientID: 0,
	}
	handler.onConnect(socket)
	var v string
	for {
		err := socket.receive(&v)
		if err != nil {
			handler.onDisconnect(socket)
			return
		}
		handler.onMessage(socket, v)
	}
}

type webSocket struct {
	ws       *websocket.Conn
	clientID int
}

func (w *webSocket) receive(v interface{}) error {
	return websocket.Message.Receive(w.ws, v)
}

func (w *webSocket) send(v interface{}) error {
	return websocket.Message.Send(w.ws, v)
}
