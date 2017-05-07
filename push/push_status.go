package push

import (
	"encoding/json"
	"time"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

const pushStatusCollection = "_PushStatus"

type pushStatus struct {
	objectID string
	db       *orm.DBController
}

func newPushStatus(objectID string) *pushStatus {
	if objectID == "" {
		objectID = utils.CreateObjectID()
	}
	p := &pushStatus{
		objectID: objectID,
		db:       orm.TomatoDBController,
	}
	return p
}

// setInitial 初始化推送状态
func (p *pushStatus) setInitial(body, where, options types.M) error {
	if body == nil {
		body = types.M{}
	}
	if options == nil {
		options = types.M{"source": "rest"}
	}

	now := time.Now().UTC()
	pushTime := time.Now().UTC()
	status := "pending"

	if t, ok := body["push_time"].(time.Time); ok {
		if config.TConfig.ScheduledPush {
			pushTime = t
			status = "scheduled"
		}
	}

	data := utils.M(body["data"])
	if data == nil {
		data = types.M{}
	}
	payloadString, _ := json.Marshal(data)
	whereString, _ := json.Marshal(where)
	pushHash := "d41d8cd98f00b204e9800998ecf8427e"

	if v, ok := data["alert"].(string); ok {
		pushHash = utils.MD5Hash(v)
	} else if v := utils.M(data["alert"]); v != nil {
		alert, _ := json.Marshal(v)
		pushHash = utils.MD5Hash(string(alert))
	}

	object := types.M{
		"objectId":  p.objectID,
		"createdAt": utils.TimetoString(now),
		"pushTime":  types.M{"__type": "Date", "iso": utils.TimetoString(pushTime)},
		"query":     string(whereString),
		"payload":   string(payloadString),
		"source":    utils.S(options["source"]),
		"title":     utils.S(options["title"]),
		"expiry":    body["expiration_time"],
		"status":    status,
		"numSent":   0,
		"pushHash":  pushHash,
		// lockdown!
		"ACL": types.M{},
	}

	return p.db.Create(pushStatusCollection, object, types.M{})
}

// setRunning 设置正在推送
func (p *pushStatus) setRunning(count int) {
	where := types.M{
		"status":   "pending",
		"objectId": p.objectID,
	}
	update := types.M{
		"status":    "running",
		"updatedAt": utils.TimetoString(time.Now().UTC()),
		"count":     count,
	}
	p.db.Update(pushStatusCollection, where, update, types.M{}, false)
}

// trackSent 推送完成，传入数据格式如下
// {
// 	"device":{
// 		"deviceType":"ios"
// 	},
// 	"transmitted":true
// }
func (p *pushStatus) trackSent(results []types.M) error {
	update := types.M{}
	numSent := 0
	numFailed := 0

	for _, result := range results {
		if result == nil {
			continue
		}
		if result["device"] == nil {
			continue
		}
		device := utils.M(result["device"])
		if device["deviceType"] == nil {
			continue
		}
		deviceType := utils.S(device["deviceType"])
		// 统计发送数据
		if result["transmitted"] != nil && result["transmitted"].(bool) {
			numSent++
			incrementOp(update, `sentPerType.`+deviceType, 1)
		} else {
			numFailed++
			incrementOp(update, `failedPerType.`+deviceType, 1)
		}
	}
	incrementOp(update, "count", -len(results))

	if numSent > 0 {
		update["numSent"] = types.M{
			"__op":   "Increment",
			"amount": numSent,
		}
	}
	if numFailed > 0 {
		update["numFailed"] = types.M{
			"__op":   "Increment",
			"amount": numFailed,
		}
	}
	update["updatedAt"] = utils.TimetoString(time.Now().UTC())

	where := types.M{
		"objectId": p.objectID,
	}

	res, err := p.db.Update(pushStatusCollection, where, update, types.M{}, false)
	if err != nil {
		return err
	}
	if res != nil {
		if c, ok := res["count"].(float64); ok && c == 0 {
			p.complete()
		} else if c, ok := res["count"].(int); ok && c == 0 {
			p.complete()
		}
	}
	return nil
}

// complete 推送完成
func (p *pushStatus) complete() {
	where := types.M{
		"objectId": p.objectID,
	}
	update := types.M{
		"status":    "succeeded",
		"count":     types.M{"__op": "Delete"},
		"updatedAt": utils.TimetoString(time.Now().UTC()),
	}
	p.db.Update(pushStatusCollection, where, update, types.M{}, false)
}

// fail 处理推送失败的情况
func (p *pushStatus) fail(err error) {
	update := types.M{
		"errorMessage": err.Error(),
		"status":       "failed",
		"updatedAt":    utils.TimetoString(time.Now().UTC()),
	}
	where := types.M{
		"objectId": p.objectID,
	}
	p.db.Update(pushStatusCollection, where, update, types.M{}, false)
}

func incrementOp(object types.M, key string, amount int) interface{} {
	value := utils.M(object[key])
	if value == nil {
		value = types.M{
			"__op":   "Increment",
			"amount": amount,
		}
	} else {
		if i, ok := value["amount"].(int); ok {
			value["amount"] = i + amount
		} else if f, ok := value["amount"].(float64); ok {
			value["amount"] = f + float64(amount)
		} else {
			value["amount"] = amount
		}
	}
	object[key] = value
	return value
}
