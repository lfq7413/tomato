package push

import (
	"strconv"
	"strings"
	"time"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

var adapter pushAdapter

// init 初始化推送模块
// 当前仅有模拟的推送模块，
// 后续添加 APNS、GCM、以及其他第三方推送模块
func init() {
	a := config.TConfig.PushAdapter
	if a == "tomato" {
		adapter = newTomatoPush()
	} else {
		adapter = newTomatoPush()
	}
}

// SendPush 发送推送消息
func SendPush(body types.M, where types.M, auth *rest.Auth, onPushStatusSaved func(string)) error {

	status := newPushStatus("")

	// TODO where 中并不包含 deviceType，此处是否有问题？
	err := validatePushType(where, adapter.getValidPushTypes())
	if err != nil {
		return err
	}

	if body["expiration_time"] != nil {
		body["expiration_time"], err = getExpirationTime(body)
		if err != nil {
			return err
		}
	}

	// TODO 检测通过立即返回，不等待推送发送完成，后续添加

	var restUpdate types.M
	var updateWhere types.M
	data := utils.M(body["data"])
	if data != nil && data["badge"] != nil {
		badge := data["badge"]
		restUpdate = types.M{}
		if strings.ToLower(utils.S(badge)) == "increment" {
			inc := types.M{
				"__op":   "Increment",
				"amount": 1,
			}
			restUpdate["badge"] = inc
		} else if v, ok := badge.(float64); ok {
			restUpdate["badge"] = v
		} else {
			return errs.E(errs.PushMisconfigured, "Invalid value for badge, expected number or 'Increment'")
		}
		updateWhere = utils.CopyMap(where)
	}

	status.setInitial(body, where, nil)

	if restUpdate != nil && updateWhere != nil {
		updateWhere["deviceType"] = "ios"
		restQuery, err := rest.NewQuery(rest.Master(), "_Installation", updateWhere, types.M{}, nil)
		if err != nil {
			status.fail(err)
			return err
		}
		err = restQuery.BuildRestWhere()
		if err != nil {
			status.fail(err)
			return err
		}
		write, err := rest.NewWrite(rest.Master(), "_Installation", restQuery.Where, restUpdate, types.M{}, nil)
		if err != nil {
			status.fail(err)
			return err
		}
		write.RunOptions["many"] = true
		_, err = write.Execute()
		if err != nil {
			status.fail(err)
			return err
		}
	}

	status.setRunning(0)

	onPushStatusSaved(status.objectID)

	// TODO 处理结果大于100的情况
	response, err := rest.Find(auth, "_Installation", where, types.M{}, nil)
	if err != nil {
		status.fail(err)
		return err
	}
	if utils.HasResults(response) == false {
		status.complete([]types.M{})
		return nil
	}
	results := utils.A(response["results"])

	res := sendToAdapter(body, results, status)
	status.complete(res)

	return nil
}

// sendToAdapter 发送推送消息
func sendToAdapter(body types.M, installations []interface{}, status *pushStatus) []types.M {
	data := utils.M(body["data"])
	if data != nil && data["badge"] != nil && strings.ToLower(utils.S(data["badge"])) == "increment" {
		badgeInstallationsMap := types.M{}
		// 按 badge 分组
		for _, v := range installations {
			installation := utils.M(v)
			var badge string
			if v, ok := installation["badge"].(float64); ok {
				badge = strconv.Itoa(int(v))
			} else {
				continue
			}
			if utils.S(installation["deviceType"]) != "ios" {
				badge = "unsupported"
			}
			installations := types.S{}
			if badgeInstallationsMap[badge] != nil {
				installations = append(installations, utils.A(badgeInstallationsMap[badge])...)
			}
			installations = append(installations, installation)
			badgeInstallationsMap[badge] = installations
		}

		var results = []types.M{}

		// 按 badge 分组发送推送
		for k, v := range badgeInstallationsMap {
			payload := utils.CopyMap(body)
			paydata := utils.M(payload["data"])
			if k == "unsupported" {
				delete(paydata, "badge")
			} else {
				paydata["badge"], _ = strconv.Atoi(k)
			}
			result := adapter.send(payload, utils.A(v), status.objectID)
			results = append(results, result...)
		}
		return results
	}

	return adapter.send(body, installations, status.objectID)
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
	} else if v, ok := expirationTimeParam.(string); ok {
		expirationTime, err = utils.StringtoTime(v)
		if err != nil {
			return nil, err
		}
	} else {
		// 时间格式错误
		return nil, errs.E(errs.PushMisconfigured, "expiration_time is not valid time.")
	}

	if expirationTime.Unix() < time.Now().Unix() {
		// 时间非法
		return nil, errs.E(errs.PushMisconfigured, "expiration_time is not valid time.")
	}

	return expirationTime.Unix() * 1000, nil
}

// pushAdapter 推送模块要实现的接口
// send() 中的 status 参数暂时没有使用
type pushAdapter interface {
	send(body types.M, installations types.S, pushStatus string) []types.M
	getValidPushTypes() []string
}
