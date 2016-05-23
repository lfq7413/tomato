package push

import (
	"strconv"
	"strings"
	"time"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
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

	status := newPushStatus()

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

	var op types.M
	var updateWhere types.M
	data := utils.MapInterface(body["data"])
	if data != nil && data["badge"] != nil {
		badge := data["badge"]
		op = types.M{}
		if strings.ToLower(utils.String(badge)) == "increment" {
			inc := types.M{
				"badge": 1,
			}
			op["$inc"] = inc
		} else if v, ok := badge.(float64); ok {
			set := types.M{
				"badge": v,
			}
			op["$set"] = set
		} else {
			return errs.E(errs.PushMisconfigured, "Invalid value for badge, expected number or 'Increment'")
		}
		updateWhere = utils.CopyMap(where)
	}

	status.setInitial(body, where, nil)

	if op != nil && updateWhere != nil {
		badgeQuery, err := rest.NewQuery(auth, "_Installation", updateWhere, types.M{})
		if err != nil {
			status.fail(err)
			return err
		}
		badgeQuery.BuildRestWhere()
		restWhere := utils.CopyMap(badgeQuery.Where)
		and := utils.SliceInterface(restWhere["$and"])
		if and == nil {
			and = types.S{badgeQuery.Where}
		}
		// badge 只有 iOS 支持，所以只发送 iOS 设备，
		tp := types.M{
			"deviceType": "ios",
		}
		and = append(and, tp)
		restWhere["$and"] = and
		err = orm.AdaptiveCollection("_Installation").UpdateMany(restWhere, op)
		if err != nil {
			status.fail(err)
			return err
		}
	}

	status.setRunning()

	onPushStatusSaved(status.objectID)

	// TODO 处理结果大于100的情况
	response, err := rest.Find(auth, "_Installation", where, types.M{})
	if err != nil {
		status.fail(err)
		return err
	}
	if utils.HasResults(response) == false {
		status.complete([]types.M{})
		return nil
	}
	results := utils.SliceInterface(response["results"])

	res := sendToAdapter(body, results, status)
	status.complete(res)

	return nil
}

// sendToAdapter 发送推送消息
func sendToAdapter(body types.M, installations []interface{}, status *pushStatus) []types.M {
	data := utils.MapInterface(body["data"])
	if data != nil && data["badge"] != nil && strings.ToLower(utils.String(data["badge"])) == "increment" {
		badgeInstallationsMap := types.M{}
		// 按 badge 分组
		for _, v := range installations {
			installation := utils.MapInterface(v)
			var badge string
			if v, ok := installation["badge"].(float64); ok {
				badge = strconv.Itoa(int(v))
			} else {
				continue
			}
			if utils.String(installation["deviceType"]) != "ios" {
				badge = "unsupported"
			}
			installations := types.S{}
			if badgeInstallationsMap[badge] != nil {
				installations = append(installations, utils.SliceInterface(badgeInstallationsMap[badge])...)
			}
			installations = append(installations, installation)
			badgeInstallationsMap[badge] = installations
		}

		var results = []types.M{}

		// 按 badge 分组发送推送
		for k, v := range badgeInstallationsMap {
			payload := utils.CopyMap(body)
			paydata := utils.MapInterface(payload["data"])
			if k == "unsupported" {
				delete(paydata, "badge")
			} else {
				paydata["badge"], _ = strconv.Atoi(k)
			}
			result := adapter.send(payload, utils.SliceInterface(v), status)
			results = append(results, result...)
		}
		return results
	}

	return adapter.send(body, installations, status)
}

// validatePushType 校验查询条件中的推送类型
// where 查询条件
// validPushTypes 当前推送模块支持的类型
func validatePushType(where types.M, validPushTypes []string) error {
	deviceTypeField := where["deviceType"]
	if deviceTypeField == nil {
		deviceTypeField = types.M{}
	}
	deviceTypes := []string{}
	if utils.String(deviceTypeField) != "" {
		deviceTypes = append(deviceTypes, utils.String(deviceTypeField))
	} else if utils.MapInterface(deviceTypeField) != nil {
		m := utils.MapInterface(deviceTypeField)
		if utils.SliceInterface(m["$in"]) != nil {
			s := utils.SliceInterface(m["$in"])
			for _, v := range s {
				deviceTypes = append(deviceTypes, utils.String(v))
			}
		}
	}
	for _, v := range deviceTypes {
		b := false
		for _, t := range validPushTypes {
			if v == t {
				b = true
				break
			}
		}
		if b == false {
			return errs.E(errs.PushMisconfigured, v+" is not supported push type.")
		}
	}

	return nil
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
	send(data types.M, installations types.S, status *pushStatus) []types.M
	getValidPushTypes() []string
}
