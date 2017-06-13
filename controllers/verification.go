package controllers

import (
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/orm"
	"github.com/lfq7413/tomato/rest"
	"github.com/lfq7413/tomato/types"
	"github.com/lfq7413/tomato/utils"
)

// VerificationController 处理 /verificationEmailRequest 接口的请求
type VerificationController struct {
	ClassesController
}

// HandleVerificationEmailRequest 处理 email 验证的请求
// @router / [post]
func (r *VerificationController) HandleVerificationEmailRequest() {
	if r.JSONBody == nil || r.JSONBody["email"] == nil {
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

	results, err := orm.TomatoDBController.Find("_User", types.M{"email": email}, types.M{})
	if err != nil {
		r.HandleError(err, 0)
		return
	}
	if len(results) < 1 {
		err = errs.E(errs.EmailNotFound, "No user found with email "+email)
		r.HandleError(err, 0)
	}

	user := utils.M(results[0])
	if user != nil {
		if emailVerified, ok := user["emailVerified"].(bool); ok && emailVerified {
			err = errs.E(errs.OtherCause, "Email "+email+" is already verified.")
			r.HandleError(err, 0)
		}
	}

	rest.SendVerificationEmail(user)
	r.Data["json"] = types.M{}
	r.ServeJSON()
}

// Get ...
// @router / [get]
func (r *VerificationController) Get() {
	r.ClassesController.Get()
}

// Delete ...
// @router / [delete]
func (r *VerificationController) Delete() {
	r.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (r *VerificationController) Put() {
	r.ClassesController.Put()
}
