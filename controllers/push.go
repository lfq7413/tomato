package controllers

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/push"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// PushController 处理 /push 接口的请求
type PushController struct {
	ClassesController
}

// HandlePost 处理发送推送消息请求
// @router / [post]
func (p *PushController) HandlePost() {
	if p.EnforceMasterKeyAccess() == false {
		return
	}
	if p.JSONBody == nil {
		p.HandleError(errs.E(errs.InvalidJSON, "request body is empty"), 0)
		return
	}
	where, err := getQueryCondition(p.JSONBody)
	if err != nil {
		p.HandleError(err, 0)
		return
	}
	onPushStatusSaved := func(pushStatusID string) {
		p.Ctx.Output.Header("X-Parse-Push-Status-Id", pushStatusID)
	}
	err = push.SendPush(p.JSONBody, where, p.Auth, onPushStatusSaved)
	if err != nil {
		p.HandleError(err, 0)
		return
	}
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
		where = utils.M(body["where"])
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
	p.ClassesController.Get()
}

// Delete ...
// @router / [delete]
func (p *PushController) Delete() {
	p.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (p *PushController) Put() {
	p.ClassesController.Put()
}
