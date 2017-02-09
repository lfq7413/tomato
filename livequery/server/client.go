package server

import (
	"encoding/json"

	"reflect"

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
	PushConnect       func(int, t.M, t.M)
	PushSubscribe     func(int, t.M, t.M)
	PushUnsubscribe   func(int, t.M, t.M)
	PushCreate        func(int, t.M, t.M)
	PushEnter         func(int, t.M, t.M)
	PushUpdate        func(int, t.M, t.M)
	PushDelete        func(int, t.M, t.M)
	PushLeave         func(int, t.M, t.M)
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
	go ws.send(msg)
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
func (c *Client) pushEvent(eventType string) func(int, t.M, t.M) {
	return func(subscriptionId int, object, originalObject t.M) {
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
			if eventType != "delete" {
				// 仅在更新与创建时去转换操作符
				transformUpdateOperators(object, originalObject)
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

	return limitedObject
}

// transformUpdateOperators 把更新操作符转换为更新之后的值
func transformUpdateOperators(object, originalObject t.M) {
	if object == nil {
		return
	}
	if originalObject == nil {
		originalObject = t.M{}
	}

	for key, v := range object {
		if value, ok := v.(map[string]interface{}); ok && value != nil {
			var op string
			if s, ok := value["__op"].(string); ok && s != "" {
				op = s
			} else {
				continue
			}
			switch op {
			case "Increment":
				var amount float64
				if a, ok := value["amount"].(float64); ok {
					amount = a
				} else if a, ok := value["amount"].(int); ok {
					amount = float64(a)
				}
				if a, ok := originalObject[key].(float64); ok {
					amount = a + amount
				} else if a, ok := originalObject[key].(int); ok {
					amount = float64(a) + amount
				}
				object[key] = amount
			case "Add":
				objects := []interface{}{}
				if objs, ok := value["objects"].([]interface{}); ok && objs != nil {
					objects = objs
				}
				if objs, ok := originalObject[key].([]interface{}); ok && objs != nil {
					objects = append(objs, objects...)
				}
				object[key] = objects
			case "AddUnique":
				objects := []interface{}{}
				if objs, ok := value["objects"].([]interface{}); ok && objs != nil {
					objects = objs
				}
				if objs, ok := originalObject[key].([]interface{}); ok && objs != nil {
					for _, obj := range objects {
						isUnique := true
						for _, obj2 := range objs {
							if reflect.DeepEqual(obj, obj2) {
								isUnique = false
								break
							}
						}
						if isUnique {
							objs = append(objs, obj)
						}
					}
					objects = objs
				}
				object[key] = objects
			case "Remove":
				objects := []interface{}{}
				if objs, ok := value["objects"].([]interface{}); ok && objs != nil {
					objects = objs
				}
				removeObjs := []interface{}{}
				if objs, ok := originalObject[key].([]interface{}); ok && objs != nil {
					for _, obj := range objs {
						canRemove := false
						for _, obj2 := range objects {
							if reflect.DeepEqual(obj, obj2) {
								canRemove = true
								break
							}
						}
						if canRemove == false {
							removeObjs = append(removeObjs, obj)
						}
					}
				}
				object[key] = removeObjs
			case "Delete":
				delete(object, key)
			case "AddRelation", "RemoveRelation":
				tp := map[string]interface{}{"__type": "Relation"}
				objects := []interface{}{}
				if objs, ok := value["objects"].([]interface{}); ok && objs != nil {
					objects = objs
				}
				for _, obj := range objects {
					if o, ok := obj.(map[string]interface{}); ok && o != nil {
						if className, ok := o["className"].(string); ok && className != "" {
							tp["className"] = className
							break
						}
					}
				}
				object[key] = tp
			}
		}
	}

}
