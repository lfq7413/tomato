package livequery

import (
	"encoding/json"

	"golang.org/x/net/websocket"
)

var server *liveQueryServer

// Run ...
func Run(args map[string]string) {
	server = &liveQueryServer{}
	server.initServer(args)
	server.run()
}

type liveQueryServer struct {
	pattern    string
	addr       string
	subscriber subscriber
}

// initServer 初始化 liveQuery 服务
func (l *liveQueryServer) initServer(args map[string]string) {
	l.pattern = args["pattern"]
	l.addr = args["addr"]

	l.subscriber = createSubscriber("", "")
	l.subscriber.subscribe("afterSave")
	l.subscriber.subscribe("afterDelete")
	var h HandlerType
	h = func(args ...string) {
		channel := args[0]
		messageStr := args[1]
		var message M
		err := json.Unmarshal([]byte(messageStr), &message)
		if err != nil {
			return
		}
		l.inflateParseObject(message)
		if channel == "afterSave" {
			l.onAfterSave(message)
		} else if channel == "afterDelete" {
			l.onAfterDelete(message)
		} else {

		}
	}
	l.subscriber.on("message", h)
}

// run 启动 WebSocket 服务
func (l *liveQueryServer) run() {
	runWebSocketServer(l.pattern, l.addr, l)
}

// onConnect 当有客户端连接成功时调用
func (l *liveQueryServer) onConnect(ws *websocket.Conn) {

}

// onMessage 当接收到客户端发来的消息时调用
func (l *liveQueryServer) onMessage(ws *websocket.Conn, msg interface{}) {

}

// onDisconnect 当客户端断开时调用
func (l *liveQueryServer) onDisconnect(ws *websocket.Conn) {

}

// inflateParseObject 展开对象
func (l *liveQueryServer) inflateParseObject(message M) {

}

// onAfterDelete 有对象删除时调用
func (l *liveQueryServer) onAfterDelete(message M) {

}

// onAfterSave 有对象保存时调用
func (l *liveQueryServer) onAfterSave(message M) {

}
