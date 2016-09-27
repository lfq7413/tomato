package controllers

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
)

// ResetController 处理 /requestPasswordReset 接口的请求
type ResetController struct {
	ClassesController
}

// HandleResetRequest 处理通过 email 重置密码的请求
// @router / [post]
func (r *ResetController) HandleResetRequest() {
	if r.JSONBody == nil && r.JSONBody["email"] == nil {
		r.HandleError(errs.E(errs.EmailMissing, "you must provide an email"), 0)
		return
	}
	var email string
	if v, ok := r.JSONBody["email"].(string); ok {
		email = v
	} else {
		r.HandleError(errs.E(errs.InvalidEmailAddress, "you must provide a valid email string"), 0)
		return
	}
	err := rest.SendPasswordResetEmail(email)
	if err != nil {
		if errs.GetErrorCode(err) == errs.ObjectNotFound {
			err = errs.E(errs.EmailNotFound, "No user found with email "+email)
		}

		r.HandleError(err, 0)
		return
	}

	r.Data["json"] = types.M{}
	r.ServeJSON()
}

// Get ...
// @router / [get]
func (r *ResetController) Get() {
	r.ClassesController.Get()
}

// Delete ...
// @router / [delete]
func (r *ResetController) Delete() {
	r.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (r *ResetController) Put() {
	r.ClassesController.Put()
}
