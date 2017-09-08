package server

import (
	"net/http"
	"strings"

	"github.com/astaxie/beego"
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
	handlerFunc := http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			s := websocket.Server{Handler: websocket.Handler(httpHandler)}
			s.ServeHTTP(w, req)
		})
	// 如果未设置监听地址，则与 beego 共用
	if addr == "" {
		// http://127.0.0.1:8080/v1 ==>> pattern = /v1
		serverURL := TomatoInfo["serverURL"]
		i := strings.Index(serverURL, `//`)
		if i < 0 {
			panic("RunWebSocketServer: invalid serverURL: " + serverURL)
		}
		serverURL = serverURL[(i + 2):]
		i = strings.Index(serverURL, `/`)
		if i < 0 {
			panic("RunWebSocketServer: invalid serverURL: " + serverURL)
		}
		pattern = serverURL[i:]
		beego.Handler(pattern, handlerFunc)
		return
	}
	// 如果设置了地址，则开启新服务去处理 WebSocket
	http.Handle(pattern, handlerFunc)
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
