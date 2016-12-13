package server

import (
	"net/http"

	"golang.org/x/net/websocket"
)

// WebSocketHandler ...
type WebSocketHandler interface {
	OnConnect(ws *WebSocket)
	OnMessage(ws *WebSocket, msg interface{})
	OnDisconnect(ws *WebSocket)
}

var handler WebSocketHandler

// RunWebSocketServer ...
func RunWebSocketServer(pattern, addr string, h WebSocketHandler) {
	handler = h
	http.Handle(pattern, http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			s := websocket.Server{Handler: websocket.Handler(httpHandler)}
			s.ServeHTTP(w, req)
		}))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func httpHandler(ws *websocket.Conn) {
	socket := &WebSocket{
		ws:       ws,
		ClientID: 0,
	}
	handler.OnConnect(socket)
	var v string
	for {
		err := socket.receive(&v)
		if err != nil {
			handler.OnDisconnect(socket)
			return
		}
		handler.OnMessage(socket, v)
	}
}

// WebSocket ...
type WebSocket struct {
	ws       *websocket.Conn
	ClientID int
}

func (w *WebSocket) receive(v interface{}) error {
	return websocket.Message.Receive(w.ws, v)
}

func (w *WebSocket) send(v interface{}) error {
	return websocket.Message.Send(w.ws, v)
}
