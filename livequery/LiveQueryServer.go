package livequery

import (
	"encoding/json"
	"strconv"

	"github.com/lfq7413/tomato/livequery/pubsub"
	"github.com/lfq7413/tomato/livequery/server"
	"github.com/lfq7413/tomato/livequery/t"
	"github.com/lfq7413/tomato/livequery/utils"
)

var s *liveQueryServer

// Run ...
func Run(args map[string]string) {
	s = &liveQueryServer{}
	s.initServer(args)
	s.run()
}

type liveQueryServer struct {
	pattern           string
	addr              string
	clientID          int
	clients           map[int]*server.Client
	subscriptions     map[string]map[string]*server.Subscription // className -> (queryHash -> subscription) TODO 增加并发锁
	keyPairs          map[string]string
	subscriber        pubsub.Subscriber
	sessionTokenCache *server.SessionTokenCache
}

// initServer 初始化 liveQuery 服务
func (l *liveQueryServer) initServer(args map[string]string) {
	l.pattern = args["pattern"]
	l.addr = args["addr"]

	l.clientID = 1
	l.clients = map[int]*server.Client{}
	l.subscriptions = map[string]map[string]*server.Subscription{}

	// 设置日志级别
	if level, ok := args["logLevel"]; ok {
		utils.TLog.Level = level
	} else {
		utils.TLog.Level = "NONE"
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
	utils.TLog.Verbose("Support key pairs", l.keyPairs)

	server.TomatoInfo["serverURL"] = args["serverURL"]
	server.TomatoInfo["appId"] = args["appId"]
	server.TomatoInfo["clientKey"] = args["clientKey"]
	server.TomatoInfo["masterKey"] = args["masterKey"]

	l.subscriber = pubsub.CreateSubscriber("", "")
	l.subscriber.Subscribe("afterSave")
	l.subscriber.Subscribe("afterDelete")
	// 向 subscriber 订阅 "message" ，当有对象操作消息时，以下处理函数将会被调起
	var h pubsub.HandlerType
	h = func(args ...string) {
		channel := args[0]
		messageStr := args[1]
		utils.TLog.Verbose("Subscribe messsage", messageStr)
		var message t.M
		err := json.Unmarshal([]byte(messageStr), &message)
		if err != nil {
			utils.TLog.Error("json.Unmarshal error")
			return
		}
		l.inflateParseObject(message)
		if channel == "afterSave" {
			l.onAfterSave(message)
		} else if channel == "afterDelete" {
			l.onAfterDelete(message)
		} else {
			utils.TLog.Error("Get message", message, "from unknown channel", channel)
		}
	}
	l.subscriber.On("message", h)
	l.sessionTokenCache = server.NewSessionTokenCache()
}

// run 启动 WebSocket 服务
func (l *liveQueryServer) run() {
	server.RunWebSocketServer(l.pattern, l.addr, l)
}

// OnConnect 当有客户端连接成功时调用
func (l *liveQueryServer) OnConnect(ws *server.WebSocket) {

}

// OnMessage 当接收到客户端发来的消息时调用
func (l *liveQueryServer) OnMessage(ws *server.WebSocket, msg interface{}) {
	var request t.M
	if message, ok := msg.(string); ok {
		err := json.Unmarshal([]byte(message), &request)
		if err != nil {
			return
		}
	}
	utils.TLog.Verbose("Request:", request)

	err := server.Validate(request, "general")
	if err != nil {
		server.PushError(ws, 1, err.Error(), true)
		utils.TLog.Error("Connect message error", err.Error())
		return
	}
	err = server.Validate(request, request["op"].(string))
	if err != nil {
		server.PushError(ws, 1, err.Error(), true)
		utils.TLog.Error("Connect message error", err.Error())
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
		server.PushError(ws, 3, "Get unknown operation", true)
		utils.TLog.Error("Get unknown operation", op)
	}
}

// OnDisconnect 当客户端断开时调用
func (l *liveQueryServer) OnDisconnect(ws *server.WebSocket) {
	utils.TLog.Log("Client disconnect:", ws.ClientID)

	clientID := ws.ClientID
	if _, ok := l.clients[clientID]; ok == false {
		utils.TLog.Error("Can not find client", clientID, "on disconnect")
		return
	}

	client := l.clients[clientID]
	delete(l.clients, clientID)

	for requestID, subscriptionInfo := range client.SubscriptionInfos {
		subscription := subscriptionInfo.Subscription
		subscription.DeleteClientSubscription(clientID, requestID)

		classSubscriptions := l.subscriptions[subscription.ClassName]
		if subscription.HasSubscribingClient() == false {
			delete(classSubscriptions, subscription.Hash)
		}

		if len(classSubscriptions) == 0 {
			delete(l.subscriptions, subscription.ClassName)
		}
	}

	utils.TLog.Verbose("Current clients", len(l.clients))
	utils.TLog.Verbose("Current subscriptions", len(l.subscriptions))
}

// inflateParseObject 展开对象
func (l *liveQueryServer) inflateParseObject(message t.M) {

}

// onAfterDelete 有对象删除时调用
func (l *liveQueryServer) onAfterDelete(message t.M) {
	utils.TLog.Verbose("afterDelete is triggered")

	deletedParseObject := message["currentParseObject"].(map[string]interface{})
	className := deletedParseObject["className"].(string)
	utils.TLog.Verbose("ClassName:", className, "| ObjectId:", deletedParseObject["objectId"])
	utils.TLog.Verbose("Current client number :", len(l.clients))

	classSubscriptions := l.subscriptions[className]
	if classSubscriptions == nil {
		utils.TLog.Error("Can not find subscriptions under this class", className)
		return
	}
	for _, subscription := range classSubscriptions {
		isSubscriptionMatched := l.matchesSubscription(deletedParseObject, subscription)
		if isSubscriptionMatched == false {
			continue
		}
		for clientID, requestIDs := range subscription.ClientRequestIDs {
			client := l.clients[clientID]
			if client == nil {
				continue
			}
			for _, requestID := range requestIDs {
				acl := deletedParseObject["ACL"].(map[string]interface{})
				isMatched, err := l.matchesACL(acl, client, requestID)
				if err != nil {
					utils.TLog.Error("Matching ACL error :", err)
				}
				if isMatched == false {
					return
				}
				client.PushDelete(requestID, deletedParseObject)
			}
		}
	}
}

// onAfterSave 有对象保存时调用
func (l *liveQueryServer) onAfterSave(message t.M) {
	utils.TLog.Verbose("afterSave is triggered")

	var originalParseObject t.M
	if message["originalParseObject"] != nil {
		originalParseObject = message["originalParseObject"].(map[string]interface{})
	}
	currentParseObject := message["currentParseObject"].(map[string]interface{})
	className := currentParseObject["className"].(string)
	utils.TLog.Verbose("ClassName:", className, "| ObjectId:", currentParseObject["objectId"])
	utils.TLog.Verbose("Current client number :", len(l.clients))

	classSubscriptions := l.subscriptions[className]
	if classSubscriptions == nil {
		utils.TLog.Error("Can not find subscriptions under this class", className)
		return
	}

	for _, subscription := range classSubscriptions {
		isOriginalSubscriptionMatched := l.matchesSubscription(originalParseObject, subscription)
		isCurrentSubscriptionMatched := l.matchesSubscription(currentParseObject, subscription)
		for clientID, requestIDs := range subscription.ClientRequestIDs {
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
					var originalACL t.M
					if originalParseObject != nil {
						originalACL = originalParseObject["ACL"].(map[string]interface{})
					}
					isOriginalMatched, err = l.matchesACL(originalACL, client, requestID)
					if err != nil {
						utils.TLog.Error("Matching ACL error :", err)
						continue
					}
				}

				var isCurrentMatched bool
				if isCurrentSubscriptionMatched == false {
					isCurrentMatched = false
				} else {
					currentACL := currentParseObject["ACL"].(map[string]interface{})
					isCurrentMatched, err = l.matchesACL(currentACL, client, requestID)
					if err != nil {
						utils.TLog.Error("Matching ACL error :", err)
						continue
					}
				}

				utils.TLog.Verbose("Original", originalParseObject,
					"| Current", currentParseObject,
					"| Match:", isOriginalSubscriptionMatched, isCurrentSubscriptionMatched, isOriginalMatched, isCurrentMatched,
					"| Query:", subscription.Hash)

				if isOriginalMatched && isCurrentMatched {
					client.PushUpdate(requestID, currentParseObject)
				} else if isOriginalMatched && !isCurrentMatched {
					client.PushLeave(requestID, currentParseObject)
				} else if !isOriginalMatched && isCurrentMatched {
					if originalParseObject != nil {
						client.PushEnter(requestID, currentParseObject)
					} else {
						client.PushCreate(requestID, currentParseObject)
					}
				} else {
					continue
				}
			}
		}
	}
}

// handleConnect 处理客户端 Connect 操作
func (l *liveQueryServer) handleConnect(ws *server.WebSocket, request t.M) {
	if l.validateKeys(request, l.keyPairs) == false {
		server.PushError(ws, 4, "Key in request is not valid", true)
		utils.TLog.Error("Key in request is not valid")
		return
	}

	client := server.NewClient(l.clientID, ws)
	ws.ClientID = l.clientID
	l.clientID++
	l.clients[ws.ClientID] = client
	utils.TLog.Log("Create new client:", ws.ClientID)
	client.PushConnect(0, nil)
}

// handleSubscribe 处理客户端 Subscribe 操作
func (l *liveQueryServer) handleSubscribe(ws *server.WebSocket, request t.M) {
	if ws.ClientID == 0 {
		server.PushError(ws, 2, "Can not find this client, make sure you connect to server before subscribing", true)
		utils.TLog.Error("Can not find this client, make sure you connect to server before subscribing")
		return
	}

	client := l.clients[ws.ClientID]

	query := request["query"].(map[string]interface{})
	subscriptionHash := utils.QueryHash(query)
	className := query["className"].(string)
	if _, ok := l.subscriptions[className]; ok == false {
		l.subscriptions[className] = map[string]*server.Subscription{}
	}
	classSubscriptions := l.subscriptions[className]
	var subscription *server.Subscription
	if s, ok := classSubscriptions[subscriptionHash]; ok {
		subscription = s
	} else {
		where := query["where"].(map[string]interface{})
		subscription = server.NewSubscription(className, where, subscriptionHash)
		classSubscriptions[subscriptionHash] = subscription
	}

	subscriptionInfo := &server.SubscriptionInfo{
		Subscription: subscription,
	}

	if fields, ok := query["fields"]; ok {
		subscriptionInfo.Fields = fields.([]string)
	}
	if sessionToken, ok := request["sessionToken"]; ok {
		subscriptionInfo.SessionToken = sessionToken.(string)
	}
	requestID := int(request["requestId"].(float64))
	client.AddSubscriptionInfo(requestID, subscriptionInfo)

	subscription.AddClientSubscription(ws.ClientID, requestID)

	client.PushSubscribe(requestID, nil)

	utils.TLog.Verbose("Create client", ws.ClientID, "new subscription:", requestID)
	utils.TLog.Verbose("Current client number:", len(l.clients))
}

// handleUnsubscribe 处理客户端 Unsubscribe 操作
func (l *liveQueryServer) handleUnsubscribe(ws *server.WebSocket, request t.M) {
	if ws.ClientID == 0 {
		server.PushError(ws, 2, "Can not find this client, make sure you connect to server before unsubscribing", true)
		utils.TLog.Error("Can not find this client, make sure you connect to server before unsubscribing")
		return
	}

	requestID := int(request["requestId"].(float64))

	client := l.clients[ws.ClientID]
	if client == nil {
		server.PushError(ws, 2, "Cannot find client with clientId "+strconv.Itoa(ws.ClientID)+". Make sure you connect to live query server before unsubscribing.", true)
		utils.TLog.Error("Can not find this client", ws.ClientID)
		return
	}

	subscriptionInfo := client.GetSubscriptionInfo(requestID)
	if subscriptionInfo == nil {
		server.PushError(ws, 2, "Cannot find subscription with clientId "+strconv.Itoa(ws.ClientID)+" subscriptionId "+strconv.Itoa(requestID)+". Make sure you subscribe to live query server before unsubscribing.", true)
		utils.TLog.Error("Can not find subscription with clientId", ws.ClientID, "subscriptionId", requestID)
		return
	}

	client.DeleteSubscriptionInfo(requestID)

	subscription := subscriptionInfo.Subscription
	className := subscription.ClassName
	subscription.DeleteClientSubscription(ws.ClientID, requestID)

	classSubscriptions := l.subscriptions[className]
	if subscription.HasSubscribingClient() == false {
		delete(classSubscriptions, subscription.Hash)
	}

	if len(classSubscriptions) == 0 {
		delete(l.subscriptions, className)
	}

	client.PushUnsubscribe(requestID, nil)
}

func (l *liveQueryServer) matchesSubscription(object t.M, subscription *server.Subscription) bool {
	if object == nil {
		return false
	}

	return utils.MatchesQuery(object, subscription.Query)
}

func (l *liveQueryServer) matchesACL(acl t.M, client *server.Client, requestID int) (bool, error) {
	if acl == nil {
		return true, nil
	}

	if getPublicReadAccess(acl) {
		return true, nil
	}

	subscriptionInfo := client.GetSubscriptionInfo(requestID)
	if subscriptionInfo == nil {
		return false, nil
	}

	subscriptionSessionToken := subscriptionInfo.SessionToken
	userID := l.sessionTokenCache.GetUserID(subscriptionSessionToken)
	if userID == "" {
		return false, nil
	}
	isSubscriptionSessionTokenMatched := getReadAccess(acl, userID)
	if isSubscriptionSessionTokenMatched {
		return true, nil
	}

	return false, nil
}

func (l *liveQueryServer) validateKeys(request t.M, validKeyPairs map[string]string) bool {
	if validKeyPairs == nil || len(validKeyPairs) == 0 {
		return true
	}
	isValid := false
	for key, secret := range validKeyPairs {
		if request[key] == nil {
			continue
		}
		if request[key].(string) != secret {
			continue
		}
		isValid = true
		break
	}

	return isValid
}

func getPublicReadAccess(acl t.M) bool {
	return getReadAccess(acl, "*")
}

// getReadAccess 需要解析的格式如下
// {
// 	"id":{
// 		"read":true,
// 		"write":true
// 	}
// 	"*":{
// 		"read":true
// 	}
// }
func getReadAccess(acl t.M, id string) bool {
	if p, ok := acl[id]; ok {
		if per, ok := p.(map[string]interface{}); ok {
			if _, ok := per["read"]; ok {
				return true
			}
			return false
		}
		return false
	}
	return false
}
