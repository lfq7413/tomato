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
	pattern       string
	addr          string
	clientID      int
	clients       map[int]*client
	subscriptions map[string]map[string]*subscription // className -> (queryHash -> subscription)
	keyPairs      map[string]string
	subscriber    subscriber
}

// initServer 初始化 liveQuery 服务
func (l *liveQueryServer) initServer(args map[string]string) {
	l.pattern = args["pattern"]
	l.addr = args["addr"]

	l.clientID = 1
	l.clients = map[int]*client{}
	l.subscriptions = map[string]map[string]*subscription{}

	// 设置日志级别
	if level, ok := args["logLevel"]; ok {
		TLog.level = level
	} else {
		TLog.level = "NONE"
	}

	// 设置 keyPairs ，用于校验客户端
	if k, ok := args["keyPairs"]; ok {
		var keyPairs map[string]string
		err := json.Unmarshal([]byte(k), &keyPairs)
		if err != nil {
			l.keyPairs = map[string]string{}
		}
		l.keyPairs = keyPairs
	} else {
		l.keyPairs = map[string]string{}
	}
	TLog.verbose("Support key pairs", l.keyPairs)

	l.subscriber = createSubscriber("", "")
	l.subscriber.subscribe("afterSave")
	l.subscriber.subscribe("afterDelete")
	// 向 subscriber 订阅 "message" ，当有对象操作消息时，以下处理函数将会被调起
	var h HandlerType
	h = func(args ...string) {
		channel := args[0]
		messageStr := args[1]
		TLog.verbose("Subscribe messsage", messageStr)
		var message M
		err := json.Unmarshal([]byte(messageStr), &message)
		if err != nil {
			TLog.error("json.Unmarshal error")
			return
		}
		l.inflateParseObject(message)
		if channel == "afterSave" {
			l.onAfterSave(message)
		} else if channel == "afterDelete" {
			l.onAfterDelete(message)
		} else {
			TLog.error("Get message", message, "from unknown channel", channel)
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
	TLog.verbose("Request:", request)

	err := validate(request, "general")
	if err != nil {
		pushError(ws, 1, err.Error(), true)
		TLog.error("Connect message error", err.Error())
		return
	}
	err = validate(request, request["op"].(string))
	if err != nil {
		pushError(ws, 1, err.Error(), true)
		TLog.error("Connect message error", err.Error())
		return
	}

	op := request["op"].(string)
	switch op {
	case "connect":
		l.handleConnect(ws, request)
	case "subscribe":
		l.handleSubscribe(ws, request)
	case "unsubscribe":
		l.handleUnsubscribe(ws, request)
	default:
		pushError(ws, 3, "Get unknown operation", true)
		TLog.error("Get unknown operation", op)
	}
}

// onDisconnect 当客户端断开时调用
func (l *liveQueryServer) onDisconnect(ws *webSocket) {
	TLog.log("Client disconnect:", ws.clientID)

	clientID := ws.clientID
	if _, ok := l.clients[clientID]; ok == false {
		TLog.error("Can not find client", clientID, "on disconnect")
		return
	}

	client := l.clients[clientID]
	delete(l.clients, clientID)

	for requestID, subscriptionInfo := range client.subscriptionInfos {
		subscription := subscriptionInfo.subscription
		subscription.deleteClientSubscription(clientID, requestID)

		classSubscriptions := l.subscriptions[subscription.className]
		if subscription.hasSubscribingClient() == false {
			delete(classSubscriptions, subscription.hash)
		}

		if len(classSubscriptions) == 0 {
			delete(l.subscriptions, subscription.className)
		}
	}

	TLog.verbose("Current clients", len(l.clients))
	TLog.verbose("Current subscriptions", len(l.subscriptions))
}

// inflateParseObject 展开对象
func (l *liveQueryServer) inflateParseObject(message M) {

}

// onAfterDelete 有对象删除时调用
func (l *liveQueryServer) onAfterDelete(message M) {
	TLog.verbose("afterDelete is triggered")

	deletedParseObject := message["currentParseObject"].(M)
	className := deletedParseObject["className"].(string)
	TLog.verbose("ClassName:", className, "| ObjectId:", deletedParseObject["objectId"])
	TLog.verbose("Current client number :", len(l.clients))

	classSubscriptions := l.subscriptions[className]
	if classSubscriptions == nil {
		TLog.error("Can not find subscriptions under this class", className)
		return
	}
	for _, subscription := range classSubscriptions {
		isSubscriptionMatched := l.matchesSubscription(deletedParseObject, subscription)
		if isSubscriptionMatched == false {
			continue
		}
		for clientID, requestIDs := range subscription.clientRequestIDs {
			client := l.clients[clientID]
			if client == nil {
				continue
			}
			for _, requestID := range requestIDs {
				acl := deletedParseObject["ACL"].(M)
				isMatched, err := l.matchesACL(acl, client, requestID)
				if err != nil {
					TLog.error("Matching ACL error :", err)
				}
				if isMatched == false {
					return
				}
				client.pushDelete(requestID, deletedParseObject)
			}
		}
	}
}

// onAfterSave 有对象保存时调用
func (l *liveQueryServer) onAfterSave(message M) {
	TLog.verbose("afterSave is triggered")

	var originalParseObject M
	if message["originalParseObject"] != nil {
		originalParseObject = message["originalParseObject"].(M)
	}
	currentParseObject := message["currentParseObject"].(M)
	className := currentParseObject["className"].(string)
	TLog.verbose("ClassName:", className, "| ObjectId:", currentParseObject["objectId"])
	TLog.verbose("Current client number :", len(l.clients))

	classSubscriptions := l.subscriptions[className]
	if classSubscriptions == nil {
		TLog.error("Can not find subscriptions under this class", className)
		return
	}

	for _, subscription := range classSubscriptions {
		isOriginalSubscriptionMatched := l.matchesSubscription(originalParseObject, subscription)
		isCurrentSubscriptionMatched := l.matchesSubscription(currentParseObject, subscription)
		for clientID, requestIDs := range subscription.clientRequestIDs {
			client := l.clients[clientID]
			if client == nil {
				continue
			}
			for _, requestID := range requestIDs {
				var err error
				var isOriginalMatched bool
				if isOriginalSubscriptionMatched == false {
					isOriginalMatched = false
				} else {
					var originalACL M
					if originalParseObject != nil {
						originalACL = originalParseObject["ACL"].(M)
					}
					isOriginalMatched, err = l.matchesACL(originalACL, client, requestID)
					if err != nil {
						TLog.error("Matching ACL error :", err)
						continue
					}
				}

				var isCurrentMatched bool
				if isCurrentSubscriptionMatched == false {
					isCurrentMatched = false
				} else {
					currentACL := currentParseObject["ACL"].(M)
					isCurrentMatched, err = l.matchesACL(currentACL, client, requestID)
					if err != nil {
						TLog.error("Matching ACL error :", err)
						continue
					}
				}

				TLog.verbose("Original", originalParseObject,
					"| Current", currentParseObject,
					"| Match:", isOriginalSubscriptionMatched, isCurrentSubscriptionMatched, isOriginalMatched, isCurrentMatched,
					"| Query:", subscription.hash)

				if isOriginalMatched && isCurrentMatched {
					client.pushUpdate(requestID, currentParseObject)
				} else if isOriginalMatched && !isCurrentMatched {
					client.pushLeave(requestID, currentParseObject)
				} else if !isOriginalMatched && isCurrentMatched {
					if originalParseObject != nil {
						client.pushEnter(requestID, currentParseObject)
					} else {
						client.pushCreate(requestID, currentParseObject)
					}
				} else {
					continue
				}
			}
		}
	}
}

// handleConnect 处理客户端 Connect 操作
func (l *liveQueryServer) handleConnect(ws *webSocket, request M) {
	if l.validateKeys(request, l.keyPairs) == false {
		pushError(ws, 4, "Key in request is not valid", true)
		TLog.error("Key in request is not valid")
		return
	}

	client := newClient(l.clientID, ws)
	ws.clientID = l.clientID
	l.clientID++
	l.clients[ws.clientID] = client
	TLog.log("Create new client:", ws.clientID)
	client.pushConnect(0, nil)
}

// handleSubscribe 处理客户端 Subscribe 操作
func (l *liveQueryServer) handleSubscribe(ws *webSocket, request M) {
	if ws.clientID == 0 {
		pushError(ws, 2, "Can not find this client, make sure you connect to server before subscribing", true)
		TLog.error("Can not find this client, make sure you connect to server before subscribing")
		return
	}

	client := l.clients[ws.clientID]

	query := request["query"].(map[string]interface{})
	subscriptionHash := queryHash(query)
	className := query["className"].(string)
	if _, ok := l.subscriptions[className]; ok == false {
		l.subscriptions[className] = map[string]*subscription{}
	}
	classSubscriptions := l.subscriptions[className]
	var subscription *subscription
	if s, ok := classSubscriptions[subscriptionHash]; ok {
		subscription = s
	} else {
		where := query["where"].(map[string]interface{})
		subscription = newSubscription(className, where, subscriptionHash)
		classSubscriptions[subscriptionHash] = subscription
	}

	subscriptionInfo := &subscriptionInfo{
		subscription: subscription,
	}

	if fields, ok := query["fields"]; ok {
		subscriptionInfo.fields = fields.([]string)
	}
	if sessionToken, ok := request["sessionToken"]; ok {
		subscriptionInfo.sessionToken = sessionToken.(string)
	}
	requestID := int(request["requestId"].(float64))
	client.addSubscriptionInfo(requestID, subscriptionInfo)

	subscription.addClientSubscription(ws.clientID, requestID)

	client.pushSubscribe(requestID, nil)

	TLog.verbose("Create client", ws.clientID, "new subscription:", requestID)
	TLog.verbose("Current client number:", len(l.clients))
}

// handleUnsubscribe 处理客户端 Unsubscribe 操作
func (l *liveQueryServer) handleUnsubscribe(ws *webSocket, request M) {

}

func (l *liveQueryServer) matchesSubscription(object M, subscription *subscription) bool {
	return false
}

func (l *liveQueryServer) matchesACL(acl M, client *client, requestID int) (bool, error) {
	return false, nil
}

func (l *liveQueryServer) validateKeys(request M, validKeyPairs map[string]string) bool {
	return false
}
