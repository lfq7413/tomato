package server

import (
	"encoding/json"

	"github.com/lfq7413/tomato/livequery/t"
)

var dafaultFields = []string{"className", "objectId", "updatedAt", "createdAt", "ACL"}

// Client ...
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

// PushError ...
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

// AddSubscriptionInfo ...
func (c *Client) AddSubscriptionInfo(requestID int, subscriptionInfo *SubscriptionInfo) {
	c.SubscriptionInfos[requestID] = subscriptionInfo
}

// GetSubscriptionInfo ...
func (c *Client) GetSubscriptionInfo(requestID int) *SubscriptionInfo {
	return c.SubscriptionInfos[requestID]
}

// DeleteSubscriptionInfo ...
func (c *Client) DeleteSubscriptionInfo(requestID int) {
	delete(c.SubscriptionInfos, requestID)
}

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
func (c *Client) toObjectWithFields(object t.M, fields []string) t.M {
	if fields == nil {
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
