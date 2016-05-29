package push

import (
	"encoding/json"
	"time"

	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

const pushStatusCollection = "_PushStatus"

type pushStatus struct {
	objectID string
	status   types.M
	db       *orm.DBController
}

func newPushStatus() *pushStatus {
	p := &pushStatus{
		objectID: utils.CreateObjectID(),
		db:       orm.TomatoDBController.WithoutValidation(),
	}
	return p
}

// setInitial 初始化推送状态
func (p *pushStatus) setInitial(body, where, options types.M) {
	if options == nil {
		options = types.M{"source": "rest"}
	}

	now := time.Now().UTC()

	data := body["data"]
	payloadString, _ := json.Marshal(data)

	object := types.M{
		"_id":         p.objectID,
		"pushTime":    utils.TimetoString(now),
		"_created_at": now,
		"query":       where,
		"payload":     string(payloadString),
		"source":      options["source"],
		"title":       options["title"],
		"expiry":      body["expiration_time"],
		"status":      "pending",
		"numSent":     0,
		"pushHash":    utils.MD5Hash(string(payloadString)),
		// lockdown!
		"_wperm": []interface{}{},
		"_rperm": []interface{}{},
	}

	err := p.db.Create(pushStatusCollection, object, types.M{})
	if err != nil {
		p.status = types.M{}
		return
	}

	p.status = types.M{
		"objectId": object["objectId"],
	}
}

// setRunning 设置正在推送
func (p *pushStatus) setRunning() {
	where := types.M{
		"status": "pending",
		"_id":    p.status["objectId"],
	}
	update := types.M{
		"$set": types.M{
			"status": "running",
		},
	}
	p.db.Update(pushStatusCollection, where, update, types.M{})
}

// complete 推送完成，传入数据格式如下
// {
// 	"device":{
// 		"deviceType":"ios"
// 	},
// 	"transmitted":true
// }
func (p *pushStatus) complete(results []types.M) {

	numSent := 0
	numFailed := 0
	sentPerType := map[string]int{}
	failedPerType := map[string]int{}

	for _, result := range results {
		if result == nil {
			continue
		}
		if result["device"] == nil {
			continue
		}
		device := utils.MapInterface(result["device"])
		if device["deviceType"] == nil {
			continue
		}
		deviceType := utils.String(device["deviceType"])
		// 统计发送数据
		if result["transmitted"] != nil && result["transmitted"].(bool) {
			numSent++
			if v, ok := sentPerType[deviceType]; ok {
				sentPerType[deviceType] = v + 1
			} else {
				sentPerType[deviceType] = 1
			}
		} else {
			numFailed++
			if v, ok := failedPerType[deviceType]; ok {
				failedPerType[deviceType] = v + 1
			} else {
				failedPerType[deviceType] = 1
			}
		}
	}

	where := types.M{
		"status": "running",
		"_id":    p.status["objectId"],
	}
	update := types.M{
		"status":        "succeeded",
		"numSent":       numSent,
		"numFailed":     numFailed,
		"sentPerType":   sentPerType,
		"failedPerType": failedPerType,
	}
	p.db.Update(pushStatusCollection, where, types.M{"$set": update}, types.M{})
}

// fail 处理推送失败的情况
func (p *pushStatus) fail(err error) {
	update := types.M{
		"errorMessage": err.Error(),
		"status":       "failed",
	}
	where := types.M{
		"_id": p.status["objectId"],
	}
	p.db.Update(pushStatusCollection, where, types.M{"$set": update}, types.M{})
}
