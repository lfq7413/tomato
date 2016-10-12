package server

import (
	"encoding/json"

	"github.com/lfq7413/tomato/livequery/t"
)

var dafaultFields = []string{"className", "objectId", "updatedAt", "createdAt", "ACL"}

// Client 客户端信息
// ws 当前对象的 WebSocket 连接
// SubscriptionInfos 当前客户端发起的所有请求对应的订阅信息
type Client struct {
	id                int
	ws                *WebSocket
	SubscriptionInfos map[int]*SubscriptionInfo
	PushConnect       func(int, t.M)
	PushSubscribe     func(int, t.M)
	PushUnsubscribe   func(int, t.M)
	PushCreate        func(int, t.M)
	PushEnter         func(int, t.M)
	PushUpdate        func(int, t.M)
	PushDelete        func(int, t.M)
	PushLeave         func(int, t.M)
}

// NewClient ...
func NewClient(id int, ws *WebSocket) *Client {
	c := &Client{
		id: id,
		ws: ws,
	}
	c.SubscriptionInfos = map[int]*SubscriptionInfo{}
	c.PushConnect = c.pushEvent("connected")
	c.PushSubscribe = c.pushEvent("subscribed")
	c.PushUnsubscribe = c.pushEvent("unsubscribed")
	c.PushCreate = c.pushEvent("create")
	c.PushEnter = c.pushEvent("enter")
	c.PushUpdate = c.pushEvent("update")
	c.PushDelete = c.pushEvent("delete")
	c.PushLeave = c.pushEvent("leave")
	return c
}

func pushResponse(ws *WebSocket, msg string) {
	ws.send(msg)
}

// PushError 发送错误信息
func PushError(ws *WebSocket, code int, errMsg string, reconnect bool) {
	errResp := t.M{
		"op":        "error",
		"error":     errMsg,
		"code":      code,
		"reconnect": reconnect,
	}
	data, err := json.Marshal(errResp)
	if err != nil {
		return
	}
	pushResponse(ws, string(data))
}

// AddSubscriptionInfo 添加 requestID 对应的订阅信息
func (c *Client) AddSubscriptionInfo(requestID int, subscriptionInfo *SubscriptionInfo) {
	c.SubscriptionInfos[requestID] = subscriptionInfo
}

// GetSubscriptionInfo 获取 requestID 对应的订阅信息
func (c *Client) GetSubscriptionInfo(requestID int) *SubscriptionInfo {
	return c.SubscriptionInfos[requestID]
}

// DeleteSubscriptionInfo 删除 requestID 对应的订阅信息
func (c *Client) DeleteSubscriptionInfo(requestID int) {
	delete(c.SubscriptionInfos, requestID)
}

// pushEvent 发送消息
func (c *Client) pushEvent(eventType string) func(int, t.M) {
	return func(subscriptionId int, object t.M) {
		response := t.M{
			"op":       eventType,
			"clientId": c.id,
		}
		if subscriptionId != 0 {
			response["requestId"] = subscriptionId
		}
		if object != nil {
			fields := []string{}
			if info, ok := c.SubscriptionInfos[subscriptionId]; ok {
				fields = info.Fields
			}
			response["object"] = c.toObjectWithFields(object, fields)
		}
		r, err := json.Marshal(response)
		if err != nil {
			return
		}
		pushResponse(c.ws, string(r))
	}
}

// toObjectWithFields 返回指定字段
func (c *Client) toObjectWithFields(object t.M, fields []string) t.M {
	if len(fields) == 0 {
		return object
	}

	limitedObject := t.M{}
	for _, field := range dafaultFields {
		limitedObject[field] = object[field]
	}

	for _, field := range fields {
		if v, ok := object[field]; ok {
			limitedObject[field] = v
		}
	}

	return nil
}
