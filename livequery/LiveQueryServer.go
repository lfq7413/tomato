package livequery

import "encoding/json"

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
func (l *liveQueryServer) onConnect(ws *webSocket) {

}

// onMessage 当接收到客户端发来的消息时调用
func (l *liveQueryServer) onMessage(ws *webSocket, msg interface{}) {
	var request M
	if message, ok := msg.(string); ok {
		err := json.Unmarshal([]byte(message), &request)
		if err != nil {
			return
		}
	}

	op := request["op"].(string)
	if op == "" {
		return
	}

	switch op {
	case "connect":
		l.handleConnect(ws, request)
	case "subscribe":
		l.handleSubscribe(ws, request)
	case "unsubscribe":
		l.handleUnsubscribe(ws, request)
	default:

	}
}

// onDisconnect 当客户端断开时调用
func (l *liveQueryServer) onDisconnect(ws *webSocket) {

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

// handleConnect 处理客户端 Connect 操作
func (l *liveQueryServer) handleConnect(ws *webSocket, request M) {

}

// handleSubscribe 处理客户端 Subscribe 操作
func (l *liveQueryServer) handleSubscribe(ws *webSocket, request M) {

}

// handleUnsubscribe 处理客户端 Unsubscribe 操作
func (l *liveQueryServer) handleUnsubscribe(ws *webSocket, request M) {

}
