package controllers

import (
	"encoding/json"

	"github.com/lfq7413/tomato/push"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// PushController ...
type PushController struct {
	ObjectsController
}

// HandlePost ...
// @router / [post]
func (p *PushController) HandlePost() {
	var body types.M
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
	p.Data["json"] = types.M{
		"result": true,
	}
	p.ServeJSON()
}

func getQueryCondition(body types.M) (types.M, error) {
	hasWhere := (body["where"] != nil)
	hasChannels := (body["channels"] != nil)

	var where types.M
	if hasWhere && hasChannels {
		// TODO 不能同时设定
		return nil, nil
	} else if hasWhere {
		where = utils.MapInterface(body["where"])
	} else if hasChannels {
		channels := types.M{
			"$in": body["channels"],
		}
		where = types.M{
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
