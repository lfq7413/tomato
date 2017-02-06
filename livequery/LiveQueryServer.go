package livequery

import (
	"encoding/json"
	"strconv"
	"sync"

	"strings"

	"github.com/lfq7413/tomato/livequery/pubsub"
	"github.com/lfq7413/tomato/livequery/server"
	"github.com/lfq7413/tomato/livequery/t"
	"github.com/lfq7413/tomato/livequery/utils"
)

/*
Parse LiveQuery Protocol Specification
https://github.com/ParsePlatform/parse-server/wiki/Parse-LiveQuery-Protocol-Specification

Connect message:
request
{
	"op": "connect",
	"restAPIKey": "",  // Optional
	"javascriptKey": "", // Optional
	"clientKey": "", //Optional
	"windowsKey": "", //Optional
	"masterKey": "", // Optional
	"sessionToken": "" // Optional
}
response
{
	"op": "connected"
}

Subscribe message:
request
{
	"op": "subscribe",
	"requestId": 1,
	"query": {
		"className": "Player",
		"where": {"name": "test"},
		"fields": ["name"] // Optional
	},
	"sessionToken": "" // Optional
}
response
{
	"op": "subscribed",
	"requestId":1
}

Update message:
request
{
	"op": "update",
	"requestId": 1,
	"query": {
		"className": "Player",
		"where": {"name": "test"},
		"fields": ["name"] // Optional
	},
	"sessionToken": "" // Optional
}
response
{
	"op": "subscribed",
	"requestId":1
}

Event message:
Suppose you subscribe like this ->
{
	"op": "subscribe",
	"requestId": 1,
	"query": {
		"className": "Player",
		"where": {"name": "test"}
	}
}
response - Create event - create object name = test
{
	"op": "create",
	"requestId": 1,
	"object": {
		"className": "Player",
		"objectId": "",
		"createdAt": "",
		"updatedAt": "",
		...
	}
}
response - Enter event - update object name unknown -> test
{
	"op": "enter",
	"requestId": 1,
	"object": {
		"className": "Player",
		"objectId": "",
		"createdAt": "",
		"updatedAt": "",
		...
	}
}
response - Update event - update object name unchanged
{
	"op": "update",
	"requestId": 1,
	"object": {
		"className": "Player",
		"objectId": "",
		"createdAt": "",
		"updatedAt": "",
		...
	}
}
response - Leave event - update object name test -> unknown
{
	"op": "leave",
	"requestId": 1,
	"object": {
		"className": "Player",
		"objectId": "",
		"createdAt": "",
		"updatedAt": "",
		...
	}
}
response - Delete event - delete object name = test
{
	"op": "delete",
	"requestId": 1,
	"object": {
		"className": "Player",
		"objectId": "",
		"createdAt": "",
		"updatedAt": "",
		...
	}
}

Unsubscribe message:
request
{
	"op": "unsubscribe",
	"requestId":1
}
response
{
	"op": "unsubscribed",
	"requestId":1
}

Error message:
response
{
  "op": "error",
  "code": 1,
  "error": "",
  "reconnect": true
}
*/

var s *liveQueryServer

type liveQueryServer struct {
	mutex             sync.Mutex
	pattern           string                                     // WebSocket 所在子地址
	addr              string                                     // WebSocket 监听地址与端口
	clientID          int                                        // 客户端 id ，递增
	clients           map[int]*server.Client                     // 当前已连接的客户端，以 clientID 为索引 TODO 增加并发锁
	subscriptions     map[string]map[string]*server.Subscription // 当前所有的订阅对象 className -> (queryHash -> subscription) TODO 增加并发锁
	keyPairs          map[string]string                          // 用于客户端鉴权的键值对，如 secretKey:abcd
	subscriber        pubsub.Subscriber                          // 订阅者
	sessionTokenCache *server.SessionTokenCache                  // 缓存 sessionToken 对应的用户 id
}

// Run 初始化 server ，启动 WebSocket
// args 支持的参数包括：
// pattern WebSocket 运行路径，例如： /livequery
// addr WebSocket 监听地址，例如： 127.0.0.1:8089
// logLevel 日志级别，包含： VERBOSE DEBUG INFO ERROR NONE ，默认为 NONE
// keyPairs 用于校验客户端权限的键值对，JSON格式字符串 例如： {"clientKey":"test"}
// serverURL tomato 地址
// appId tomato 对应的 appId
// clientKey tomato 对应的 clientKey
// masterKey tomato 对应的 masterKey
// subType 订阅服务类型，支持 EventEmitter Redis
// subURL 订阅服务地址，如果是 EventEmitter 可不填写
func Run(args map[string]string) {
	s = &liveQueryServer{}
	s.initServer(args)
	s.run()
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

	// 初始化 tomato 服务参数，用于获取用户信息
	server.TomatoInfo["serverURL"] = args["serverURL"]
	server.TomatoInfo["appId"] = args["appId"]
	server.TomatoInfo["clientKey"] = args["clientKey"]
	server.TomatoInfo["masterKey"] = args["masterKey"]

	// 向 subscriber 订阅 afterSave 、 afterDelete 两个频道
	l.subscriber = pubsub.CreateSubscriber(args["subType"], args["subURL"], args["subConfig"])
	l.subscriber.Subscribe(server.TomatoInfo["appId"] + "afterSave")
	l.subscriber.Subscribe(server.TomatoInfo["appId"] + "afterDelete")

	// 设置从 subscriber 接收到消息时的处理函数
	var h pubsub.HandlerType
	h = func(args ...string) {
		if len(args) < 2 {
			return
		}
		channel := args[0]
		messageStr := args[1]
		utils.TLog.Verbose("Subscribe messsage", messageStr)
		var message t.M
		err := json.Unmarshal([]byte(messageStr), &message)
		if err != nil {
			utils.TLog.Error("unable to parse message", []byte(messageStr), err)
			return
		}
		l.inflateParseObject(message)
		if channel == server.TomatoInfo["appId"]+"afterSave" {
			l.onAfterSave(message)
		} else if channel == server.TomatoInfo["appId"]+"afterDelete" {
			l.onAfterDelete(message)
		} else {
			utils.TLog.Error("Get message", message, "from unknown channel", channel)
		}
	}
	l.subscriber.On("message", h)

	// 设置 cache
	l.sessionTokenCache = server.NewSessionTokenCache()
}

// run 启动 WebSocket 服务
// liveQueryServer 通过 OnConnect 、 OnMessage 、 OnDisconnect 处理与客户端的交互
func (l *liveQueryServer) run() {
	server.RunWebSocketServer(l.pattern, l.addr, l)
}

// OnConnect 当有客户端连接到 WebSocket 时调用
func (l *liveQueryServer) OnConnect(ws *server.WebSocket) {

}

// OnMessage 当 WebSocket 接收到客户端发来的消息时调用
func (l *liveQueryServer) OnMessage(ws *server.WebSocket, msg interface{}) {
	var request t.M
	if message, ok := msg.(string); ok {
		err := json.Unmarshal([]byte(message), &request)
		if err != nil {
			utils.TLog.Error("unable to parse request", []byte(message), err)
			return
		}
	} else {
		return
	}
	utils.TLog.Verbose("Request:", request)

	// 校验 op 操作否是否支持
	err := server.Validate(request, "general")
	if err != nil {
		server.PushError(ws, 1, err.Error(), true)
		utils.TLog.Error("Connect message error", err.Error())
		return
	}
	// 校验 指定的操作符 格式是否正确
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
	case "update":
		l.handleUpdateSubscription(ws, request)
	case "unsubscribe":
		l.handleUnsubscribe(ws, request, true)
	default:
		server.PushError(ws, 3, "Get unknown operation", true)
		utils.TLog.Error("Get unknown operation", op)
	}
}

// OnDisconnect 当客户端从 WebSocket 断开时调用
func (l *liveQueryServer) OnDisconnect(ws *server.WebSocket) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
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

// onAfterDelete 从 subscriber 中接收到对象删除消息时调用
func (l *liveQueryServer) onAfterDelete(message t.M) {
	utils.TLog.Verbose("afterDelete is triggered")

	deletedParseObject := message["currentParseObject"].(map[string]interface{})
	className := deletedParseObject["className"].(string)
	utils.TLog.Verbose("ClassName:", className, "| ObjectId:", deletedParseObject["objectId"])
	utils.TLog.Verbose("Current client number :", len(l.clients))

	// 取出当前类对应的订阅对象列表
	classSubscriptions := l.subscriptions[className]
	if classSubscriptions == nil {
		utils.TLog.Error("Can not find subscriptions under this class", className)
		return
	}
	for _, subscription := range classSubscriptions {
		// 检测要删除的对象是否符合订阅条件
		isSubscriptionMatched := l.matchesSubscription(deletedParseObject, subscription)
		if isSubscriptionMatched == false {
			continue
		}
		// 如果符合订阅条件，则向指定的 client 的 request 返回对象
		for clientID, requestIDs := range subscription.ClientRequestIDs {
			client := l.clients[clientID]
			if client == nil {
				continue
			}
			for _, requestID := range requestIDs {
				var acl map[string]interface{}
				if v, ok := deletedParseObject["ACL"].(map[string]interface{}); ok {
					acl = v
				}
				// 检测 client 是否有权限接收这条删除信息
				isMatched := l.matchesACL(acl, client, requestID)
				if isMatched == false {
					continue
				}
				// 向 client 发送删除的对象
				client.PushDelete(requestID, deletedParseObject)
			}
		}
	}
}

// onAfterSave 从 subscriber 中接收到对象保存消息时调用
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
	// 取出当前类对应的订阅信息列表
	classSubscriptions := l.subscriptions[className]
	if classSubscriptions == nil {
		utils.TLog.Error("Can not find subscriptions under this class", className)
		return
	}

	for _, subscription := range classSubscriptions {
		// 检测要保存的对象是否符合订阅条件
		isOriginalSubscriptionMatched := l.matchesSubscription(originalParseObject, subscription)
		isCurrentSubscriptionMatched := l.matchesSubscription(currentParseObject, subscription)
		// 均不符合则跳过
		if isOriginalSubscriptionMatched == false && isCurrentSubscriptionMatched == false {
			continue
		}
		for clientID, requestIDs := range subscription.ClientRequestIDs {
			client := l.clients[clientID]
			if client == nil {
				continue
			}
			for _, requestID := range requestIDs {
				// 检测 client 是否有权限接收这条保存信息
				var isOriginalMatched bool
				if isOriginalSubscriptionMatched == false {
					isOriginalMatched = false
				} else {
					var originalACL t.M
					if originalParseObject != nil {
						if v, ok := originalParseObject["ACL"].(map[string]interface{}); ok {
							originalACL = v
						}
					}
					isOriginalMatched = l.matchesACL(originalACL, client, requestID)
				}

				var isCurrentMatched bool
				if isCurrentSubscriptionMatched == false {
					isCurrentMatched = false
				} else {
					var currentACL t.M
					if v, ok := currentParseObject["ACL"].(map[string]interface{}); ok {
						currentACL = v
					}
					isCurrentMatched = l.matchesACL(currentACL, client, requestID)
				}

				utils.TLog.Verbose("Original", originalParseObject,
					"| Current", currentParseObject,
					"| Match:", isOriginalSubscriptionMatched, isCurrentSubscriptionMatched, isOriginalMatched, isCurrentMatched,
					"| Query:", subscription.Hash)

				if isOriginalMatched && isCurrentMatched {
					// 原对象与新对象均符合条件，则为 Update
					client.PushUpdate(requestID, currentParseObject)
				} else if isOriginalMatched && !isCurrentMatched {
					// 原对象符合条件，但是新对象不符合，则为 Leave
					client.PushLeave(requestID, currentParseObject)
				} else if !isOriginalMatched && isCurrentMatched {
					if originalParseObject != nil {
						// 原对象不符合条件，但是新对象符合，则为 Enter
						client.PushEnter(requestID, currentParseObject)
					} else {
						// 原对象不存在，同时新对象符合条件，则为 Create
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
	// 校验是否含有必要的键值对
	if l.validateKeys(request, l.keyPairs) == false {
		server.PushError(ws, 4, "Key in request is not valid", true)
		utils.TLog.Error("Key in request is not valid")
		return
	}

	// 创建新的 client 并更新 l.clientID
	client := server.NewClient(l.clientID, ws)
	ws.ClientID = l.clientID
	l.clientID++
	l.mutex.Lock()
	l.clients[ws.ClientID] = client
	l.mutex.Unlock()
	utils.TLog.Log("Create new client:", ws.ClientID)
	client.PushConnect(0, nil)
}

// handleSubscribe 处理客户端 Subscribe 操作
func (l *liveQueryServer) handleSubscribe(ws *server.WebSocket, request t.M) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if ws.ClientID == 0 {
		server.PushError(ws, 2, "Can not find this client, make sure you connect to server before subscribing", true)
		utils.TLog.Error("Can not find this client, make sure you connect to server before subscribing")
		return
	}

	client := l.clients[ws.ClientID]
	if client == nil {
		server.PushError(ws, 2, "Can not find this client, make sure you connect to server before subscribing", true)
		utils.TLog.Error("Can not find this client, make sure you connect to server before subscribing")
		return
	}

	query := request["query"].(map[string]interface{})
	// 计算 query 的 hash ，参与计算的字段包括： className 与 where
	subscriptionHash := utils.QueryHash(query)
	className := query["className"].(string)
	if _, ok := l.subscriptions[className]; ok == false {
		l.subscriptions[className] = map[string]*server.Subscription{}
	}
	// 取出当前 className 对应的 订阅对象列表
	classSubscriptions := l.subscriptions[className]
	// 取出当前 queryHash 对应的 订阅对象，如果没有则根据 queryHash 生成，并保存到 订阅对象列表
	var subscription *server.Subscription
	if s, ok := classSubscriptions[subscriptionHash]; ok {
		subscription = s
	} else {
		where := query["where"].(map[string]interface{})
		subscription = server.NewSubscription(className, where, subscriptionHash)
		classSubscriptions[subscriptionHash] = subscription
	}

	// 生成订阅信息对象，用于设置到 client 中
	subscriptionInfo := &server.SubscriptionInfo{
		Subscription: subscription,
	}
	if fields, ok := query["fields"]; ok {
		fieldsArray := []string{}
		// query["fields"] 已经过校验，确定格式为 []string ，无需再次校验
		for _, fld := range fields.([]interface{}) {
			fieldsArray = append(fieldsArray, fld.(string))
		}
		subscriptionInfo.Fields = fieldsArray
	}
	if sessionToken, ok := request["sessionToken"]; ok {
		subscriptionInfo.SessionToken = sessionToken.(string)
	}
	requestID := int(request["requestId"].(float64))
	// 根据 requestID ，把订阅信息对象设置到 client 中
	client.AddSubscriptionInfo(requestID, subscriptionInfo)
	// 更新订阅对象，添加使用该对象的 ClientID 与 requestID
	subscription.AddClientSubscription(ws.ClientID, requestID)
	// 订阅成功
	client.PushSubscribe(requestID, nil)

	utils.TLog.Verbose("Create client", ws.ClientID, "new subscription:", requestID)
	utils.TLog.Verbose("Current client number:", len(l.clients))
}

// handleUpdateSubscription 处理客户端 update 操作
func (l *liveQueryServer) handleUpdateSubscription(ws *server.WebSocket, request t.M) {
	l.handleUnsubscribe(ws, request, false)
	l.handleSubscribe(ws, request)
}

// handleUnsubscribe 处理客户端 Unsubscribe 操作
// notifyClient 默认为 true，即为通知客户端取消订阅行为
func (l *liveQueryServer) handleUnsubscribe(ws *server.WebSocket, request t.M, notifyClient bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
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
	// 取出 requestID 对应的订阅信息对象
	subscriptionInfo := client.GetSubscriptionInfo(requestID)
	if subscriptionInfo == nil {
		server.PushError(ws, 2, "Cannot find subscription with clientId "+strconv.Itoa(ws.ClientID)+" subscriptionId "+strconv.Itoa(requestID)+". Make sure you subscribe to live query server before unsubscribing.", true)
		utils.TLog.Error("Can not find subscription with clientId", ws.ClientID, "subscriptionId", requestID)
		return
	}
	// 从 client 中删除 requestID 对应的 订阅信息对象
	client.DeleteSubscriptionInfo(requestID)
	// 取出 订阅对象
	subscription := subscriptionInfo.Subscription
	className := subscription.ClassName
	// 更新订阅对象，删除使用该对象的 ClientID 与 requestID
	subscription.DeleteClientSubscription(ws.ClientID, requestID)
	// 如果没有任何 client 订阅该对象，则从 订阅对象列表中删除
	classSubscriptions := l.subscriptions[className]
	if subscription.HasSubscribingClient() == false {
		delete(classSubscriptions, subscription.Hash)
	}
	// 如果当前类对应的 订阅对象列表 长度为 0 ，则删除往前类对应的列表
	if len(classSubscriptions) == 0 {
		delete(l.subscriptions, className)
	}
	if notifyClient == false {
		return
	}
	// 退订成功
	client.PushUnsubscribe(requestID, nil)
}

// matchesSubscription 检测对象是否符合订阅条件
func (l *liveQueryServer) matchesSubscription(object t.M, subscription *server.Subscription) bool {
	if object == nil {
		return false
	}

	return utils.MatchesQuery(object, subscription.Query)
}

// matchesACL 检测客户端是否有权限接收消息
func (l *liveQueryServer) matchesACL(acl t.M, client *server.Client, requestID int) bool {
	if acl == nil {
		return true
	}

	if getPublicReadAccess(acl) {
		return true
	}

	subscriptionInfo := client.GetSubscriptionInfo(requestID)
	if subscriptionInfo == nil {
		return false
	}

	subscriptionSessionToken := subscriptionInfo.SessionToken
	userID := l.sessionTokenCache.GetUserID(subscriptionSessionToken)
	if userID == "" {
		return false
	}
	isSubscriptionSessionTokenMatched := getReadAccess(acl, userID)
	if isSubscriptionSessionTokenMatched {
		return true
	}

	// 检测用户的角色是否符合 acl
	aclHasRoles := false
	for key := range acl {
		if strings.HasPrefix(key, "role:") {
			aclHasRoles = true
			break
		}
	}
	if aclHasRoles == false {
		return false
	}

	roles := server.GetUserRoles(userID)
	for _, role := range roles {
		if getReadAccess(acl, role) {
			return true
		}
	}

	return false
}

// validateKeys 校验 connect 请求中是否包含必要的键值对
func (l *liveQueryServer) validateKeys(request t.M, validKeyPairs map[string]string) bool {
	if validKeyPairs == nil || len(validKeyPairs) == 0 {
		return true
	}
	isValid := false
	for key, secret := range validKeyPairs {
		if request[key] == nil {
			continue
		}
		if v, ok := request[key].(string); ok {
			if v != secret {
				continue
			}
		} else {
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
	if acl == nil {
		return true
	}
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
