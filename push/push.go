package push

import (
	"strconv"
	"time"

	"github.com/lfq7413/tomato/config"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/utils"
)

var adapter pushAdapter

func init() {
	a := config.TConfig.PushAdapter
	if a == "tomato" {
		adapter = newTomatoPush()
	} else {
		adapter = newTomatoPush()
	}
}

// SendPush ...
func SendPush(body map[string]interface{}, where map[string]interface{}, auth *rest.Auth) error {
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

	var op map[string]interface{}
	var updateWhere map[string]interface{}
	data := utils.MapInterface(body["data"])
	if data != nil && data["badge"] != nil {
		badge := data["badge"]
		op = map[string]interface{}{}
		if utils.String(badge) == "Increment" {
			inc := map[string]interface{}{
				"badge": 1,
			}
			op["$inc"] = inc
		} else if v, ok := badge.(float64); ok {
			set := map[string]interface{}{
				"badge": v,
			}
			op["$set"] = set
		} else {
			// TODO badge 值不符合要求
			return nil
		}
		updateWhere = utils.CopyMap(where)
	}

	if op != nil && updateWhere != nil {
		badgeQuery := rest.NewQuery(auth, "_Installation", updateWhere, map[string]interface{}{})
		badgeQuery.BuildRestWhere()
		restWhere := utils.CopyMap(badgeQuery.Where)
		and := utils.SliceInterface(restWhere["$and"])
		if and == nil {
			restWhere["$and"] = []interface{}{badgeQuery.Where}
		}
		tp := map[string]interface{}{
			"deviceType": "ios",
		}
		and = append(and, tp)
		restWhere["$and"] = and
		orm.AdaptiveCollection("_Installation").UpdateMany(restWhere, op)
	}

	response := rest.Find(auth, "_Installation", where, map[string]interface{}{})
	if utils.HasResults(response) == false {
		return nil
	}
	results := utils.SliceInterface(response["results"])

	if data != nil && data["badge"] != nil && utils.String(data["badge"]) == "Increment" {
		badgeInstallationsMap := map[string]interface{}{}
		for _, v := range results {
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
			installations := []interface{}{}
			if badgeInstallationsMap[badge] != nil {
				installations = append(installations, utils.SliceInterface(badgeInstallationsMap[badge])...)
			}
			installations = append(installations, installation)
			badgeInstallationsMap[badge] = installations
		}

		for k, v := range badgeInstallationsMap {
			payload := utils.CopyMap(body)
			paydata := utils.MapInterface(payload["data"])
			if k == "unsupported" {
				delete(paydata, "badge")
			} else {
				paydata["badge"], _ = strconv.Atoi(k)
			}
			adapter.send(payload, utils.SliceInterface(v))
		}
		return nil
	}
	adapter.send(body, results)

	return nil
}

func validatePushType(where map[string]interface{}, validPushTypes []string) error {
	deviceTypeField := where["deviceType"]
	if deviceTypeField == nil {
		return nil
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
			// TODO 不支持的类型
			return nil
		}
	}

	return nil
}

func getExpirationTime(body map[string]interface{}) (interface{}, error) {
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
		// TODO 时间格式错误
		return nil, nil
	}

	if expirationTime.Unix() < time.Now().Unix() {
		// TODO 时间非法
		return nil, nil
	}

	return expirationTime.Unix() * 1000, nil
}

type pushAdapter interface {
	send(data map[string]interface{}, installations []interface{})
	getValidPushTypes() []string
}
