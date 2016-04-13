package controllers

import (
	"encoding/json"

	"github.com/lfq7413/tomato/types"
)

// ResetController ...
type ResetController struct {
	ObjectsController
}

// HandleResetRequest ...
// @router / [post]
func (r *ResetController) HandleResetRequest() {
	var object types.M
	json.Unmarshal(r.Ctx.Input.RequestBody, &object)
	if object == nil && object["email"] == nil {
		// TODO 需要 email
		return
	}
	// TODO 发送邮件
	r.Data["json"] = types.M{}
	r.ServeJSON()
}

// Get ...
// @router / [get]
func (r *ResetController) Get() {
	r.ObjectsController.Get()
}

// Delete ...
// @router / [delete]
func (r *ResetController) Delete() {
	r.ObjectsController.Delete()
}

// Put ...
// @router / [put]
func (r *ResetController) Put() {
	r.ObjectsController.Put()
}
