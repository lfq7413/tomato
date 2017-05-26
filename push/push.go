package push

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

var adapter pushAdapter
var queue *pushQueue
var worker *pushWorker

// init 初始化推送模块
// 当前仅有模拟的推送模块，
// 后续添加 APNS、GCM、以及其他第三方推送模块
func init() {
	a := config.TConfig.PushAdapter
	if a == "tomato" {
		adapter = newTomatoPush()
	} else if a == "FCM" {
		adapter = newFCMPush()
	} else {
		adapter = nil
	}

	worker = newPushWorker(adapter, config.TConfig.PushChannel)
	queue = newPushQueue(config.TConfig.PushChannel, config.TConfig.PushBatchSize)
}

// SendPush 发送推送消息
func SendPush(body types.M, where types.M, auth *rest.Auth, onPushStatusSaved func(string)) error {
	if adapter == nil {
		return errs.E(errs.PushMisconfigured, "Missing push configuration")
	}

	// validatePushType(where, adapter.getValidPushTypes())

	if body["expiration_time"] != nil {
		var err error
		body["expiration_time"], err = getExpirationTime(body)
		if err != nil {
			return err
		}
	}

	if body["push_time"] != nil {
		pushTime, err := getPushTime(body)
		if err != nil {
			return err
		}
		if pushTime != nil {
			body["push_time"] = pushTime
		}
	}

	badgeUpdate := func() error { return nil }

	data := utils.M(body["data"])
	if data != nil && data["badge"] != nil {
		restUpdate := types.M{}
		badge := data["badge"]
		restUpdate = types.M{}
		if strings.ToLower(utils.S(badge)) == "increment" {
			restUpdate["badge"] = types.M{
				"__op":   "Increment",
				"amount": 1,
			}
		} else if v, ok := badge.(float64); ok {
			restUpdate["badge"] = v
		} else if v, ok := badge.(int); ok {
			restUpdate["badge"] = v
		} else {
			return errors.New("Invalid value for badge, expected number or 'Increment'")
		}
		updateWhere := utils.CopyMapM(where)

		badgeUpdate = func() error {
			updateWhere["deviceType"] = "ios"
			restQuery, err := rest.NewQuery(rest.Master(), "_Installation", updateWhere, types.M{}, nil)
			if err != nil {
				return err
			}
			err = restQuery.BuildRestWhere()
			if err != nil {
				return err
			}
			write, err := rest.NewWrite(rest.Master(), "_Installation", restQuery.Where, restUpdate, types.M{}, nil)
			if err != nil {
				return err
			}
			write.RunOptions["many"] = true
			_, err = write.Execute()
			return err
		}
	}

	status := newPushStatus("")

	err := status.setInitial(body, where, nil)
	if err != nil {
		return err
	}

	onPushStatusSaved(status.objectID)
	err = badgeUpdate()
	if err != nil {
		status.fail(err)
		return err
	}

	if _, ok := body["push_time"]; ok && config.TConfig.ScheduledPush {

	} else {
		err = queue.enqueue(body, where, auth, status)
	}

	if err != nil {
		status.fail(err)
	}

	return err
}

// getExpirationTime 把过期时间转换为以毫秒为单位的 Unix 时间
func getExpirationTime(body types.M) (interface{}, error) {
	expirationTimeParam := body["expiration_time"]
	if expirationTimeParam == nil {
		return nil, nil
	}

	var expirationTime time.Time
	var err error
	if v, ok := expirationTimeParam.(float64); ok {
		expirationTime = time.Unix(int64(v), 0)
	} else if v, ok := expirationTimeParam.(int); ok {
		expirationTime = time.Unix(int64(v), 0)
	} else if v, ok := expirationTimeParam.(string); ok {
		expirationTime, err = utils.StringtoTime(v)
		if err != nil {
			return nil, errs.E(errs.PushMisconfigured, fmt.Sprint(expirationTimeParam, "is not valid time."))
		}
	} else {
		// 时间格式错误
		return nil, errs.E(errs.PushMisconfigured, fmt.Sprint(expirationTimeParam, "is not valid time."))
	}

	if expirationTime.Unix() < time.Now().Unix() {
		// 时间非法
		return nil, errs.E(errs.PushMisconfigured, fmt.Sprint(expirationTimeParam, "is not valid time."))
	}

	return expirationTime.Unix() * 1000, nil
}

// getPushTime 获取推送时间
func getPushTime(body types.M) (interface{}, error) {
	pushTimeParam := body["push_time"]
	if pushTimeParam == nil {
		return nil, nil
	}

	var pushTime time.Time
	var err error

	if v, ok := pushTimeParam.(float64); ok {
		pushTime = time.Unix(int64(v), 0)
	} else if v, ok := pushTimeParam.(int); ok {
		pushTime = time.Unix(int64(v), 0)
	} else if v, ok := pushTimeParam.(string); ok {
		pushTime, err = utils.StringtoTime(v)
		if err != nil {
			return nil, errs.E(errs.PushMisconfigured, fmt.Sprint(pushTimeParam, "is not valid time."))
		}
	} else {
		// 时间格式错误
		return nil, errs.E(errs.PushMisconfigured, fmt.Sprint(pushTimeParam, "is not valid time."))
	}

	return pushTime, nil
}

// pushAdapter 推送模块要实现的接口
// send() 中的 status 参数暂时没有使用
type pushAdapter interface {
	// send 发送消息
	/*
		body 数据格式：
		{
			"channels":["aaa","bbb"],
			"where":{
				"key":"v"
			},
			"push_time":time.Time("2015-03-13T22:05:08Z"),
			"expiration_interval": 518400,
			"expiration_time": 14xxxxxxxxx,
			"data":{
				"alert":"hello world."
				"badge":"Increment",
				"sound":"cheering.caf",
				"content-available":1,
				"category":"aaa",
				"uri":"xxxx",
				"title":"hello"
			}
		}
	*/
	send(body types.M, installations types.S, pushStatus string) []types.M
	getValidPushTypes() []string
}
