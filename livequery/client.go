package livequery

import "encoding/json"

var dafaultFields = []string{"className", "objectId", "updatedAt", "createdAt", "ACL"}

type client struct {
	id                int
	ws                *webSocket
	subscriptionInfos map[int]*subscriptionInfo
	pushConnect       func(int, M)
	pushSubscribe     func(int, M)
	pushUnsubscribe   func(int, M)
	pushCreate        func(int, M)
	pushEnter         func(int, M)
	pushUpdate        func(int, M)
	pushDelete        func(int, M)
	pushLeave         func(int, M)
}

func newClient(id int, ws *webSocket) *client {
	c := &client{
		id: id,
		ws: ws,
	}
	c.subscriptionInfos = map[int]*subscriptionInfo{}
	c.pushConnect = c.pushEvent("connected")
	c.pushSubscribe = c.pushEvent("subscribed")
	c.pushUnsubscribe = c.pushEvent("unsubscribed")
	c.pushCreate = c.pushEvent("create")
	c.pushEnter = c.pushEvent("enter")
	c.pushUpdate = c.pushEvent("update")
	c.pushDelete = c.pushEvent("delete")
	c.pushLeave = c.pushEvent("leave")
	return c
}

func pushResponse(ws *webSocket, msg string) {
	ws.send(msg)
}

func pushError(ws *webSocket, code int, errMsg string, reconnect bool) {
	errResp := M{
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

func (c *client) addSubscriptionInfo(requestID int, subscriptionInfo *subscriptionInfo) {
	c.subscriptionInfos[requestID] = subscriptionInfo
}

func (c *client) getSubscriptionInfo(requestID int) *subscriptionInfo {
	return c.subscriptionInfos[requestID]
}

func (c *client) deleteSubscriptionInfo(requestID int) {
	delete(c.subscriptionInfos, requestID)
}

func (c *client) pushEvent(eventType string) func(int, M) {
	return func(subscriptionId int, object M) {
		response := M{
			"op":        eventType,
			"clientId":  c.id,
			"requestId": subscriptionId,
		}
		if object != nil {
			fields := []string{}
			if info, ok := c.subscriptionInfos[subscriptionId]; ok {
				fields = info.fields
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
func (c *client) toObjectWithFields(object M, fields []string) M {
	if fields == nil {
		return object
	}

	limitedObject := M{}
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
