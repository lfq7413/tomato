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
func (p *pushStatus) setInitial(body, where, options types.M) {
	if options == nil {
		options = types.M{"source": "rest"}
	}

	now := time.Now().UTC()

	data := body["data"]
	payloadString, _ := json.Marshal(data)
	var pushHash = "d41d8cd98f00b204e9800998ecf8427e"
	if d := utils.M(data); d != nil {
		if v, ok := d["alert"].(string); ok {
			pushHash = utils.MD5Hash(v)
		} else if v := utils.M(d["alert"]); v != nil {
			alert, _ := json.Marshal(v)
			pushHash = utils.MD5Hash(string(alert))
		}
	}

	object := types.M{
		"objectId":  p.objectID,
		"pushTime":  utils.TimetoString(now),
		"createdAt": utils.TimetoString(now),
		"query":     where,
		"payload":   string(payloadString),
		"source":    options["source"],
		"title":     options["title"],
		"expiry":    body["expiration_time"],
		"status":    "pending",
		"numSent":   0,
		"pushHash":  pushHash,
		// lockdown!
		"ACL": types.M{},
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
func (p *pushStatus) setRunning(count int) {
	where := types.M{
		"status":   "pending",
		"objectId": p.status["objectId"],
	}
	update := types.M{
		"status":    "running",
		"updatedAt": utils.TimetoString(time.Now().UTC()),
	}
	p.db.Update(pushStatusCollection, where, update, types.M{}, false)
}

func (p *pushStatus) trackSent(results []types.M) error {
	// TODO
	return nil
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
		device := utils.M(result["device"])
		if device["deviceType"] == nil {
			continue
		}
		deviceType := utils.S(device["deviceType"])
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
		"status":   "running",
		"objectId": p.status["objectId"],
	}
	update := types.M{
		"status":        "succeeded",
		"numSent":       numSent,
		"numFailed":     numFailed,
		"sentPerType":   sentPerType,
		"failedPerType": failedPerType,
		"updatedAt":     utils.TimetoString(time.Now().UTC()),
	}
	p.db.Update(pushStatusCollection, where, types.M{"$set": update}, types.M{}, false)
}

// fail 处理推送失败的情况
func (p *pushStatus) fail(err error) {
	update := types.M{
		"errorMessage": err.Error(),
		"status":       "failed",
		"updatedAt":    utils.TimetoString(time.Now().UTC()),
	}
	where := types.M{
		"objectId": p.status["objectId"],
	}
	p.db.Update(pushStatusCollection, where, types.M{"$set": update}, types.M{}, false)
}
