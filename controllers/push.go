package controllers

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/push"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// PushController 处理 /push 接口的请求
type PushController struct {
	ObjectsController
}

// HandlePost 处理发送推送消息请求
// @router / [post]
func (p *PushController) HandlePost() {
	if p.JSONBody == nil {
		p.Data["json"] = errs.ErrorMessageToMap(errs.InvalidJSON, "request body is empty")
		p.ServeJSON()
		return
	}
	where, err := getQueryCondition(p.JSONBody)
	if err != nil {
		p.Data["json"] = errs.ErrorToMap(err)
		p.ServeJSON()
		return
	}
	push.SendPush(p.JSONBody, where, p.Auth, false)
	p.Data["json"] = types.M{"result": true}
	p.ServeJSON()
}

// getQueryCondition 获取查询条件
func getQueryCondition(body types.M) (types.M, error) {
	hasWhere := (body["where"] != nil)
	hasChannels := (body["channels"] != nil)

	var where types.M
	if hasWhere && hasChannels {
		// 查询与频道不能同时设定
		return nil, errs.E(errs.PushMisconfigured, "Channels and query can not be set at the same time.")
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
		return nil, errs.E(errs.PushMisconfigured, `Sending a push requires either "channels" or a "where" query.`)
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
