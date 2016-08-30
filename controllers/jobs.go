package controllers

import (
	"github.com/lfq7413/tomato/cloud"
	"github.com/lfq7413/tomato/errs"
	"github.com/lfq7413/tomato/types"
)

// JobsController 处理 /jobs 接口的请求
type JobsController struct {
	ClassesController
}

// HandleCloudJob 执行后台任务
// @router /:jobName [post]
func (j *JobsController) HandleCloudJob() {
	if j.Auth.IsMaster == false {
		j.Data["json"] = errs.ErrorMessageToMap(errs.OperationForbidden, "need master key")
		j.ServeJSON()
		return
	}
	jobName := j.Ctx.Input.Param(":jobName")
	theJob := cloud.GetJob(jobName)
	if theJob == nil {
		j.Data["json"] = errs.ErrorMessageToMap(errs.ScriptFailed, "Invalid job function.")
		j.ServeJSON()
		return
	}

	if j.JSONBody == nil {
		j.JSONBody = types.M{}
	}

	request := cloud.JobRequest{
		Params: j.JSONBody,
		Master: false,
	}
	if j.Auth != nil {
		request.Master = j.Auth.IsMaster
		request.User = j.Auth.User
	}

	go theJob(request)

	j.Data["json"] = types.M{}
	j.ServeJSON()

}

// Get ...
// @router / [get]
func (j *JobsController) Get() {
	j.ClassesController.Get()
}

// Post ...
// @router / [post]
func (j *JobsController) Post() {
	j.ClassesController.Post()
}

// Delete ...
// @router / [delete]
func (j *JobsController) Delete() {
	j.ClassesController.Delete()
}

// Put ...
// @router / [put]
func (j *JobsController) Put() {
	j.ClassesController.Put()
}
