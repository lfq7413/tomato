package controllers

import "encoding/json"
import "github.com/lfq7413/tomato/push"
import "github.com/lfq7413/tomato/utils"

// PushController ...
type PushController struct {
	ObjectsController
}

// HandlePost ...
// @router / [post]
func (p *PushController) HandlePost() {
	var body map[string]interface{}
	json.Unmarshal(p.Ctx.Input.RequestBody, &body)
	if body == nil {
		// TODO
		return
	}
	where, err := getQueryCondition(body)
	if err != nil {
		// TODO
		return
	}
	push.SendPush(body, where, p.Auth)
	p.Data["json"] = map[string]interface{}{
		"result": true,
	}
	p.ServeJSON()
}

func getQueryCondition(body map[string]interface{}) (map[string]interface{}, error) {
	hasWhere := (body["where"] != nil)
	hasChannels := (body["channels"] != nil)

	var where map[string]interface{}
	if hasWhere && hasChannels {
		// TODO 不能同时设定
		return nil, nil
	} else if hasWhere {
		where = utils.MapInterface(body["where"])
	} else if hasChannels {
		channels := map[string]interface{}{
			"$in": body["channels"],
		}
		where = map[string]interface{}{
			"channels": channels,
		}
	} else {
		// TODO 至少有一个
		return nil, nil
	}

	return where, nil
}

// Get ...
// @router / [get]
func (p *PushController) Get() {
	p.ObjectsController.Get()
}

// Delete ...
// @router / [delete]
func (p *PushController) Delete() {
	p.ObjectsController.Delete()
}

// Put ...
// @router / [put]
func (p *PushController) Put() {
	p.ObjectsController.Put()
}
