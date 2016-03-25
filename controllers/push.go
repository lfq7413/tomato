package controllers

import "encoding/json"
import "github.com/lfq7413/tomato/push"

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
	where := getQueryCondition(body)
	push.SendPush(body, where, p.Auth)
	p.Data["json"] = map[string]interface{}{
		"result": true,
	}
	p.ServeJSON()
}

func getQueryCondition(body map[string]interface{}) map[string]interface{} {
	return nil
}

// Get ...
// @router / [get]
func (p *PushController) Get() {
	p.Controller.Get()
}

// Delete ...
// @router / [delete]
func (p *PushController) Delete() {
	p.Controller.Delete()
}

// Put ...
// @router / [put]
func (p *PushController) Put() {
	p.Controller.Put()
}
