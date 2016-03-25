package push

import "github.com/lfq7413/tomato/rest"
import "github.com/lfq7413/tomato/config"
import "github.com/lfq7413/tomato/utils"

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
		updateWhere = map[string]interface{}{}
		for k, v := range where {
			updateWhere[k] = v
		}
	}

	if op != nil && updateWhere != nil {
		badgeQuery := rest.NewQuery(auth, "_Installation", updateWhere, map[string]interface{}{})
		badgeQuery.BuildRestWhere()
		restWhere := map[string]interface{}{}
		for k, v := range badgeQuery.Where {
			restWhere[k] = v
		}
		and := utils.SliceInterface(restWhere["$and"])
		if and == nil {
			restWhere["$and"] = []interface{}{badgeQuery.Where}
		}
		tp := map[string]interface{}{
			"deviceType": "ios",
		}
		and = append(and, tp)
		restWhere["$and"] = and
		// TODO 更新 badge
	}

	return nil
}

func validatePushType(where map[string]interface{}, validPushTypes []string) error {
	return nil
}

func getExpirationTime(body map[string]interface{}) (string, error) {
	return "", nil
}

type pushAdapter interface {
	send()
	getValidPushTypes() []string
}
